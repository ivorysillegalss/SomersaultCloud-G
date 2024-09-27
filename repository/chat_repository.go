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
	__proto "SomersaultCloud/proto/.proto"
	"context"
	"encoding/json"
	"github.com/jinzhu/gorm"
	jsoniter "github.com/json-iterator/go"
	"github.com/thoas/go-funk"
	"strconv"
	"sync"
	"time"
)

// 存在内存中 	维护一个map 记录哪个bot是使用re表的 哪个是使用chat表的
// 饿汉式单例
var botId2TableMap map[int]string = make(map[int]string)

type chatRepository struct {
	redis redis.Client
	mysql mysql.Client
	env   *bootstrap.Env
}

type historySerializer struct {
	gzipRecord *[]*domain.Record
	pbRecord   *[]*__proto.Record
}

func NewChatRepository(dbs *bootstrap.Databases, env *bootstrap.Env) domain.ChatRepository {
	botId2TableMap[dao.DefaultModelBotId] = dao.RefactorTable
	botId2TableMap[dao.MathSolveBotId] = dao.OriginTable
	botId2TableMap[dao.CommentBotId] = dao.OriginTable
	return &chatRepository{redis: dbs.Redis, mysql: dbs.Mysql, env: env}
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

func (c *chatRepository) AsyncSaveHistory(ctx context.Context, chatId int, askText string, generationText string, botId int) {

	switch botId2TableMap[botId] {
	case dao.OriginTable:
		originTableSaveHistory(c.mysql.Gorm(), chatId, askText, generationText)
	case dao.RefactorTable:
		refactorTableSaveHistory(c.mysql.Gorm(), chatId, askText, generationText, c.env)
	default:
		log.GetTextLogger().Error("get wrong botId mapping table type")
	}
}

// 写入数据库的聊天记录映射类
type recordToStruct struct {
	ID     int `gorm:"primaryKey column:record_id" `
	ChatId int `gorm:"column:chat_id"`
}

func originTableSaveHistory(db *gorm.DB, chatId int, askText string, generationText string) {
	r := &recordToStruct{
		ChatId: chatId,
	}

	//主键回显获取自增后ID
	if err := db.Table("record_info").Create(r).Error; err != nil {
		log.GetTextLogger().Error(err.Error())
	}

	chatAsk := &domain.ChatAsk{
		RecordId: r.ID,
		ChatId:   chatId,
		Message:  askText,
		Time:     time.Now().Unix(),
	}

	chatGeneration := &domain.ChatGeneration{
		RecordId: r.ID,
		Message:  generationText,
		Time:     time.Now().Unix(),
	}

	if err := db.Table("chat_ask").Save(chatAsk).Error; err != nil {
		log.GetTextLogger().Fatal(err.Error())
	}
	if err := db.Table("chat_generation").Save(chatGeneration).Error; err != nil {
		log.GetTextLogger().Fatal(err.Error())
	}
	if err := db.Table("chat").Where("chat_id = ?", chatId).Update("last_update_time", time.Now().Unix()).Error; err != nil {
		log.GetTextLogger().Error(err.Error())
	}
}

func refactorTableSaveHistory(db *gorm.DB, chatId int, askText string, generationText string, env *bootstrap.Env) {
	history, _, err := refactorTableGetHistory(db, chatId, env)
	if err != nil {
		log.GetJsonLogger().WithFields("async history", err.Error()).Warn("async history error")
		panic(err)
	}

	var marshal []byte
	switch env.Serializer {
	case sys.GzipCompress:
		r := &domain.Record{
			ChatAsks:        &domain.ChatAsk{Message: askText},
			ChatGenerations: &domain.ChatGeneration{Message: generationText},
		}
		var records []*domain.Record
		records = make([]*domain.Record, 0)
		records = append(records, r)
		if funk.NotEmpty(history) {
			*history.gzipRecord = append(*history.gzipRecord, records...)
		} else {
			history.gzipRecord = &records
		}
		marshal, err = compressutil.NewCompress(sys.GzipCompress).CompressData(*history.gzipRecord)
		if err != nil {
			panic(err)
		}
	case sys.ProtoBufCompress:
		r := &__proto.Record{
			ChatAsks:        &__proto.ChatAsk{Message: askText},
			ChatGenerations: &__proto.ChatGeneration{Message: generationText},
		}
		var records []*__proto.Record
		records = make([]*__proto.Record, 0)
		records = append(records, r)
		if funk.NotEmpty(history) {
			*history.pbRecord = append(*history.pbRecord, records...)
		} else {
			history.pbRecord = &records
		}
		marshal, err = compressutil.NewCompress(sys.ProtoBufCompress).CompressData(*history.pbRecord)
		if err != nil {
			panic(err)
		}
	}

	if err = db.Table("chat_re").Where("chat_id = ?", chatId).Update("data", marshal).Error; err != nil {
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
	//重构的分成新表和旧表
	switch botType {
	case dao.RefactorTable:

		//如果是新表 根据系统配置的序列化方式 选择gzip or protobuf
		switch c.env.Serializer {
		case sys.GzipCompress:
			getHistory, s, err := refactorTableGetHistory(c.mysql.Gorm(), chatId, c.env)
			return getHistory.gzipRecord, s, err
		case sys.ProtoBufCompress:
			getHistory, s, err := refactorTableGetHistory(c.mysql.Gorm(), chatId, c.env)
			if funk.IsEmpty(getHistory) {
				return nil, s, err
			}
			history2Domain := convertPbHistory2Domain(getHistory)
			return &history2Domain, s, err
		default:
			log.GetTextLogger().Fatal("error Compress sign")
			panic("error Compress sign")
		}

	case dao.OriginTable:
		return originTableGetHistory(c.mysql.Gorm(), chatId)
	default:
		log.GetTextLogger().Fatal("bot mapping db table error")
		return nil, common.ZeroString, nil
	}

}

// 使用带缓冲区channel 将历史记录切片依次切换
func convertPbHistory2Domain(h *historySerializer) []*domain.Record {
	pbRecords := *h.pbRecord
	domainRecords := make([]*domain.Record, len(pbRecords))
	var wg sync.WaitGroup

	// 创建一个带缓冲区的 channel，容量等于 pbRecords 长度
	results := make(chan struct {
		index  int
		record *domain.Record
	}, len(pbRecords))

	for i, pbRecord := range pbRecords {
		wg.Add(1)
		pbRecord := pbRecord
		i := i
		go func() {
			defer wg.Done()
			domainRecord := convertToDomainRecord(pbRecord)
			results <- struct {
				index  int
				record *domain.Record
			}{i, domainRecord} // 将结果发送到 channel
		}()
	}

	// 另一个 goroutine 等待所有任务完成后关闭 channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// 从 channel 中读取结果并写入切片
	for result := range results {
		domainRecords[result.index] = result.record
	}

	return domainRecords
}

func convertToDomainRecord(record *__proto.Record) *domain.Record {
	asks := record.ChatAsks
	generations := record.ChatGenerations
	return &domain.Record{
		RecordId: int(record.RecordId),
		ChatAsks: &domain.ChatAsk{
			ChatId:  int(asks.ChatId),
			Message: asks.Message,
			BotId:   int(asks.BotId),
			Time:    asks.Time,
		},
		ChatGenerations: &domain.ChatGeneration{
			Message: generations.Message,
			Time:    generations.Time,
		},
	}
}

type historyData struct {
	Data  []byte `gorm:"column:data"`
	Title string `gorm:"column:title"`
}

// refactorTableGetHistory  重构后的新表获取记录的方法
func refactorTableGetHistory(db *gorm.DB, chatId int, env *bootstrap.Env) (*historySerializer, string, error) {
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

	var err error
	var gzipHistory []*domain.Record
	var pbHistory []*__proto.Record
	var hs historySerializer
	switch env.Serializer {
	case sys.GzipCompress:
		err = compressutil.NewCompress(sys.GzipCompress).DecompressData(h[0].Data, &gzipHistory)
		hs.gzipRecord = &gzipHistory
	case sys.ProtoBufCompress:
		if funk.IsEmpty(h[0].Data) {
			return &hs, h[0].Title, nil
		}
		err = compressutil.NewCompress(sys.ProtoBufCompress).DecompressData(h[0].Data, &pbHistory)
		hs.pbRecord = &pbHistory
	default:
		log.GetTextLogger().Fatal("error Compress sign")
		panic("error Compress sign")
	}

	if err != nil {
		return nil, common.ZeroString, err
	}

	return &hs, h[0].Title, nil
}

// TODO 得测，重写
func originTableGetHistory(db *gorm.DB, chatId int) (*[]*domain.Record, string, error) {
	var records []*domain.Record
	var titles []string
	//TODO 没问题的话切换为异步
	db.Table("chat").Where("chat_id = ?", chatId).Pluck("title", &titles)
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

	return &records, titles[0], nil
}

func (c *chatRepository) DbInsertNewChat(ctx context.Context, userId int, botId int) {
	var data []byte
	switch c.env.Serializer {
	case sys.GzipCompress:
		marshal, _ := jsoniter.Marshal(dao.DefaultData)
		data, _ = compressutil.NewCompress(sys.GzipCompress).CompressData(marshal)
	case sys.ProtoBufCompress:
		data, _ = compressutil.NewCompress(sys.ProtoBufCompress).CompressData(&__proto.Record{})
	}
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
