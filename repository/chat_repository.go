package repository

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/dao"
	"SomersaultCloud/constant/sys"
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/log"
	"SomersaultCloud/infrastructure/lru"
	"SomersaultCloud/infrastructure/mysql"
	"SomersaultCloud/infrastructure/redis"
	"SomersaultCloud/internal/compressutil"
	"context"
	"encoding/json"
	"github.com/jinzhu/gorm"
	jsoniter "github.com/json-iterator/go"
	"github.com/thoas/go-funk"
	"strconv"
	"time"
)

// 存在内存中 	维护一个map 记录哪个bot是使用re表的 哪个是使用chat表的
// 饿汉式单例
var botId2TableMap map[int]string = make(map[int]string)

type chatRepository struct {
	redis redis.Client
	mysql mysql.Client
}

func getInfixType(botId int) (originInfix string) {
	originInfix = common.ZeroString
	switch botId2TableMap[botId] {
	case dao.RefactorTable:
		originInfix = common.ZeroString
	case dao.OriginTable:
		originInfix = cache.OriginTable + common.Infix

	default:
		log.GetTextLogger().Error("get wrong botId mapping table type")
	}
	return originInfix
}

func (c *chatRepository) CacheGetNewestChatId(ctx context.Context) int {
	newestId, err := c.redis.Get(ctx, cache.NewestChatIdKey)
	if err != nil {
		return common.FalseInt
	}
	newId, err := strconv.Atoi(newestId)
	if err != nil {
		return common.FalseInt
	}
	return newId
}

func (c *chatRepository) CacheInsertNewChat(ctx context.Context, id int) {
	_ = c.redis.Set(ctx, cache.NewestChatIdKey, id+1)
	// MQ异步写入sql lua改进 暂时无用
}

func (c *chatRepository) CacheLuaInsertNewChatId(ctx context.Context, luaScript string, k string) (int, error) {
	res, err := c.redis.ExecuteLuaScript(ctx, luaScript, k)
	if err != nil {
		return common.FalseInt, err
	}
	return int(res.(int64)), nil
}

func (c *chatRepository) CacheGetHistory(ctx context.Context, chatId int, botId int) (history *[]*domain.Record, isCache bool, isErr error) {
	//TODO 获取旧表历史记录可能会有bug,待测试
	var h []*domain.Record
	var v string
	var err error
	v, err = c.redis.Get(ctx, cache.ChatHistory+common.Infix+getInfixType(botId)+strconv.Itoa(chatId))
	_ = jsoniter.Unmarshal([]byte(v), &h)
	if c.redis.IsEmpty(err) {
		return nil, true, nil
	}
	if err != nil {
		return nil, true, err
	}
	return &h, false, nil
}

func (c *chatRepository) AsyncSaveHistory(ctx context.Context, chatId int, askText string, generationText string) {

	r := &domain.Record{
		ChatAsks:        &domain.ChatAsk{Message: askText},
		ChatGenerations: &domain.ChatGeneration{Message: generationText},
	}
	var records []*domain.Record
	records = make([]*domain.Record, 0)
	records = append(records, r)

	//history, _, err := c.DbGetHistory(ctx, chatId)
	history, _, err := refactorTableGetHistory(c.mysql.Gorm(), chatId)
	if err != nil {
		log.GetJsonLogger().WithFields("async history", err.Error()).Warn("async history error")
		panic(err)
	}

	if funk.NotEmpty(history) {
		*history = append(*history, records...)
	} else {
		history = &records
	}

	marshal, err := compressutil.NewCompress(sys.GzipCompress).CompressData(*history)
	if err != nil {
		panic(err)
	}
	if err = c.mysql.Gorm().Table("chat_re").Where("chat_id = ?", chatId).Update("data", marshal).Error; err != nil {
		panic(err)
	}
}

func (c *chatRepository) CacheGetGeneration(ctx context.Context, chatId int) (*domain.GenerationResponse, error) {
	hGet, err := c.redis.HGet(ctx, cache.ChatGenerationExpired, strconv.Itoa(chatId))
	if c.redis.IsEmpty(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	s := hGet.(string)
	keyExpiredTTL, _ := strconv.Atoi(s)
	currentTime := int(time.Now().Unix())

	var resAny any
	if currentTime > keyExpiredTTL {
		return nil, nil
	} else {
		resAny, _ = c.redis.HGet(ctx, cache.ChatGeneration, strconv.Itoa(chatId))
	}

	resStr := resAny.(string)
	var res domain.GenerationResponse
	err = json.Unmarshal([]byte(resStr), &res)

	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *chatRepository) MemoryGetGeneration(ctx context.Context, chatId int) *domain.GenerationResponse {
	return chatGenerationMap[chatId]
}

func (c *chatRepository) CacheDelGeneration(ctx context.Context, chatId int) error {
	return c.redis.Del(ctx, cache.ChatGeneration+common.Infix+strconv.Itoa(chatId))
}

func (c *chatRepository) MemoryDelGeneration(ctx context.Context, chatId int) {
	delete(chatGenerationMap, chatId)
}

func (c *chatRepository) CacheLuaLruResetHistory(ctx context.Context, cacheKey string, history *[]*domain.Record, chatId int, title string, botId int) error {
	originInfix := getInfixType(botId)

	marshalToString, err2 := jsoniter.MarshalToString(*history)
	if err2 != nil {
		log.GetJsonLogger().WithFields("marshal", err2.Error()).Warn("reset history failed")
	}

	err2 = c.redis.Set(ctx, cache.ChatHistory+common.Infix+originInfix+strconv.Itoa(chatId), marshalToString)
	if err2 != nil {
		log.GetJsonLogger().WithFields("redis_set", err2.Error()).Warn("reset history failed")
	}

	err2 = c.redis.Set(ctx, cache.ChatHistoryTitle+common.Infix+originInfix+strconv.Itoa(chatId), title)
	if err2 != nil {
		log.GetJsonLogger().WithFields("redis_set", err2.Error()).Warn("reset history failed")
	}

	newLru := lru.NewLru(cache.ContextLruMaxCapacity, cache.RedisZSetType, c.redis)
	//返回最老的一个元素
	err, oldest := newLru.Add(ctx, cacheKey, strconv.Itoa(chatId))
	if funk.NotEqual(oldest, common.FalseInt) || funk.NotEqual(oldest, common.ZeroInt) {
		//证明有元素被移除了
		_ = c.redis.Del(ctx, cache.ChatHistory+common.Infix+originInfix+strconv.Itoa(oldest))
		_ = c.redis.Del(ctx, cache.ChatHistoryTitle+common.Infix+originInfix+strconv.Itoa(oldest))
	}
	return err
}

func (c *chatRepository) CacheLuaLruPutHistory(ctx context.Context, cacheKey string, history *[]*domain.Record, askText string, generationText string, chatId int, botId int, title string) error {

	r := &domain.Record{
		ChatAsks:        &domain.ChatAsk{Message: askText},
		ChatGenerations: &domain.ChatGeneration{Message: generationText},
	}
	if funk.IsEmpty(history) {
		history = new([]*domain.Record)
	}
	a := append(*history, r)

	//控制单chat内最大缓存数量
	if len(a) > cache.HistoryDefaultWeight {
		a = a[len(a)-cache.HistoryDefaultWeight:]
	}

	marshalToString, err2 := jsoniter.MarshalToString(a)
	if err2 != nil {
		log.GetJsonLogger().WithFields("marshal_res", err2.Error()).Warn("lru put history failed")
	}

	originInfix := getInfixType(botId)
	_ = c.redis.Set(ctx, cache.ChatHistory+common.Infix+originInfix+strconv.Itoa(chatId), marshalToString)
	_ = c.redis.Set(ctx, cache.ChatHistoryTitle+common.Infix+originInfix+strconv.Itoa(chatId), title)

	newLru := lru.NewLru(cache.ContextLruMaxCapacity, cache.RedisZSetType, c.redis)
	//返回最老的一个元素
	err, oldest := newLru.Add(ctx, cacheKey, strconv.Itoa(chatId))
	if !(funk.Equal(oldest, common.FalseInt) || funk.Equal(oldest, common.ZeroInt)) {
		//TODO 这里应该是有问题的,没有正确删除LRU外数据
		//证明有元素被移除了
		_ = c.redis.Del(ctx, cache.ChatHistory+common.Infix+originInfix+strconv.Itoa(oldest))
		_ = c.redis.Del(ctx, cache.ChatHistoryTitle+common.Infix+originInfix+strconv.Itoa(oldest))
	}
	if c.redis.IsEmpty(err) {
		return nil
	}
	return err
}

// DbGetHistory TODO db中数据流式更新flink 更新入Hbase等
//
//	目前的架构一旦db获取历史记录 就是全部获取 初步思路是定时任务 mq以某个时间段为界限（eg7天） 将数据流式更新
//	保证不出现大Key等 主要是为了提高查询效率
func (c *chatRepository) DbGetHistory(ctx context.Context, chatId int, botId int) (history *[]*domain.Record, title string, err error) {
	botType := botId2TableMap[botId]
	switch botType {
	case dao.RefactorTable:
		return refactorTableGetHistory(c.mysql.Gorm(), chatId)
	case dao.OriginTable:
		return originTableGetHistory(c.mysql.Gorm(), chatId)
	default:
		log.GetTextLogger().Fatal("bot mapping db table error")
		return nil, common.ZeroString, nil
	}

}

// TODO 映射数据库历史纪录结构体，修改
type historyData struct {
	Data  []byte `gorm:"column:data"`
	Title string `gorm:"column:title"`
}

// refactorTableGetHistory  重构后的新表获取记录的方法
func refactorTableGetHistory(db *gorm.DB, chatId int) (*[]*domain.Record, string, error) {
	var h []*historyData

	if err := db.Table("chat_re").Where("chat_id = ?", chatId).
		Select("data,title").
		Scan(&h).
		Error; err != nil {
		return nil, common.ZeroString, err
	}

	if funk.IsEmpty(h) {
		return nil, common.ZeroString, nil
	}

	var history []*domain.Record
	err := compressutil.NewCompress(sys.GzipCompress).DecompressData(h[0].Data, &history)
	if err != nil {
		return nil, common.ZeroString, err
	}

	return &history, h[0].Title, nil
}

// TODO 重写
func originTableGetHistory(db *gorm.DB, chatId int) (*[]*domain.Record, string, error) {
	var records []*domain.Record
	var title string
	//TODO 没问题的话切换为异步
	db.Table("chat").Where("chat_id = ?", chatId).Pluck("title", &title)
	err := db.Table("record_info").Where("chat_id = ?", chatId).Find(&records).Error
	if err != nil {
		return nil, common.ZeroString, err
	}
	for index, record := range records {

		//如果获取了足够的历史记录 直接跳出 不再获取
		if funk.Equal(index, cache.HistoryDefaultWeight) {
			break
		}

		// 确保 ChatAsks 和 ChatGenerations 是指向结构体的指针
		if records[index].ChatAsks == nil {
			records[index].ChatAsks = &domain.ChatAsk{}
		}
		if records[index].ChatGenerations == nil {
			records[index].ChatGenerations = &domain.ChatGeneration{}
		}

		err := db.Table("chat_ask").Where("record_id = ?", record.RecordId).First(records[index].ChatAsks).Error
		//如果同一段chat在数据库中没找到记录 有可能是这个机器人这一次不需要问题
		if err != nil && err.Error() != dao.RecordNotFoundError {
			return nil, common.ZeroString, nil
		}
		err = db.Table("chat_generation").Where("record_id = ?", record.RecordId).First(records[index].ChatGenerations).Error
		if err != nil && err.Error() != dao.RecordNotFoundError {
			return nil, common.ZeroString, nil
		}
	}

	return &records, title, nil
}

func (c *chatRepository) DbInsertNewChat(ctx context.Context, userId int, botId int) {
	marshal, _ := jsoniter.Marshal(dao.DefaultData)
	data, _ := compressutil.NewCompress(sys.GzipCompress).CompressData(marshal)
	chat := &domain.Chat{
		UserId:         userId,
		BotId:          botId,
		Title:          dao.DefaultTitle,
		LastUpdateTime: time.Now().Unix(),
		IsDelete:       false,
		Data:           data,
	}

	if err := c.mysql.Gorm().Table("chat_re").Create(chat).Error; err != nil {
		//TODO 异步 写入日志
	}
	return
}

// CacheGetTitles 通过LRU找出权重相关 根据权重找出具体值
func (c *chatRepository) CacheGetTitles(ctx context.Context, userId int, botId int) ([]*domain.TitleData, error) {
	newLru := lru.NewLru(cache.ContextLruMaxCapacity, cache.RedisZSetType, c.redis)
	list, err := newLru.List(ctx, cache.ChatHistoryScore+common.Infix+strconv.Itoa(userId)+common.Infix+strconv.Itoa(botId))
	//获取到指定bot下的热数据 最近交互的5个会话的chatId
	if err != nil {
		return nil, err
	}
	var res []*domain.TitleData
	for _, chatId := range list {
		t, err := c.redis.Get(ctx, cache.ChatHistoryTitle+common.Infix+getInfixType(botId)+chatId)
		if err != nil {
			return nil, err
		}
		chatId, _ := strconv.Atoi(chatId)
		v1 := &domain.TitleData{Title: t, ChatId: chatId}
		res = append(res, v1)
	}
	return res, nil
}

func (c *chatRepository) CacheGetTitlePrompt(ctx context.Context) string {
	v, _ := c.redis.Get(ctx, cache.HistoryTitlePrompt)
	return v
}

func (c *chatRepository) DbUpdateTitle(ctx context.Context, chatId int, newTitle string) {
	if err := c.mysql.Gorm().Table("chat_re").Where("chat_id = ?", chatId).Update("title", newTitle).Error; err != nil {
		log.GetJsonLogger().WithFields("chat_id", chatId).Error("update title error")
	}
}

// CacheUpdateTitle 默认是有这个缓存的 因为这个方法的执行时机在第一次聊天后 并且已经获得了聊天记录。所以一般情况下必定有此次聊天记录的LRU缓存
func (c *chatRepository) CacheUpdateTitle(ctx context.Context, chatId int, newTitle string, botId int) {
	err := c.redis.Set(ctx, cache.ChatHistoryTitle+common.Infix+getInfixType(botId)+strconv.Itoa(chatId), newTitle)
	if err != nil {
		log.GetTextLogger().Error("update title failed")
	}
}

func NewChatRepository(dbs *bootstrap.Databases) domain.ChatRepository {
	botId2TableMap[dao.DefaultModelBotId] = dao.RefactorTable
	botId2TableMap[dao.MathSolveBotId] = dao.OriginTable
	botId2TableMap[dao.CommentBotId] = dao.OriginTable
	return &chatRepository{redis: dbs.Redis, mysql: dbs.Mysql}
}
