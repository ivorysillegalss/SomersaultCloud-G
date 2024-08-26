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
	jsoniter "github.com/json-iterator/go"
	"github.com/thoas/go-funk"
	"strconv"
	"time"
)

type chatRepository struct {
	redis redis.Client
	mysql mysql.Client
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

func (c *chatRepository) CacheGetHistory(ctx context.Context, chatId int) (history *[]*domain.Record, isCache bool, isErr error) {
	var h []*domain.Record
	v, err := c.redis.Get(ctx, cache.ChatHistory+common.Infix+strconv.Itoa(chatId))
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

	history, err := c.DbGetHistory(ctx, chatId)
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

func (c *chatRepository) CacheLuaLruResetHistory(ctx context.Context, cacheKey string, history *[]*domain.Record, chatId int) error {
	marshalToString, err2 := jsoniter.MarshalToString(*history)
	if err2 != nil {
		log.GetJsonLogger().WithFields("marshal", err2.Error()).Warn("reset history failed")
	}

	err2 = c.redis.Set(ctx, cache.ChatHistory+common.Infix+strconv.Itoa(chatId), marshalToString)
	if err2 != nil {
		log.GetJsonLogger().WithFields("redis_set", err2.Error()).Warn("reset history failed")
	}

	newLru := lru.NewLru(cache.ContextLruMaxCapacity, cache.RedisZSetType, c.redis)
	//返回最老的一个元素
	err, oldest := newLru.Add(ctx, cacheKey, strconv.Itoa(chatId))
	if funk.NotEqual(oldest, common.FalseInt) || funk.NotEqual(oldest, common.ZeroInt) {
		//证明有元素被移除了
		_ = c.redis.Del(ctx, cache.ChatHistory+common.Infix+strconv.Itoa(oldest))
	}
	return err
}

func (c *chatRepository) CacheLuaLruPutHistory(ctx context.Context, cacheKey string, history *[]*domain.Record, askText string, generationText string, chatId int) error {

	r := &domain.Record{
		ChatAsks:        &domain.ChatAsk{Message: askText},
		ChatGenerations: &domain.ChatGeneration{Message: generationText},
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

	_ = c.redis.Set(ctx, cache.ChatHistory+common.Infix+strconv.Itoa(chatId), marshalToString)

	newLru := lru.NewLru(cache.ContextLruMaxCapacity, cache.RedisZSetType, c.redis)
	//返回最老的一个元素
	err, oldest := newLru.Add(ctx, cacheKey, strconv.Itoa(chatId))
	if !(funk.Equal(oldest, common.FalseInt) || funk.Equal(oldest, common.ZeroInt)) {
		//证明有元素被移除了
		_ = c.redis.Del(ctx, cache.ChatHistory+common.Infix+strconv.Itoa(oldest))
	}
	return err
}

// DbGetHistory TODO db中数据流式更新flink 更新入Hbase等
//
//	目前的架构一旦db获取历史记录 就是全部获取 初步思路是定时任务 mq以某个时间段为界限（eg7天） 将数据流式更新
//	保证不出现大Key等 主要是为了提高查询效率
func (c *chatRepository) DbGetHistory(ctx context.Context, chatId int) (*[]*domain.Record, error) {
	var h []*domain.Record
	var data [][]byte

	//pluck方法在不构建结构体的前提下 获取单个字段的值
	if err := c.mysql.Gorm().Table("chat_re").
		Where("chat_id = ?", chatId).Pluck("data", &data).Error; err != nil {
		return nil, err
	}

	if funk.IsEmpty(data[0]) {
		return nil, nil
	}

	err := compressutil.NewCompress(sys.GzipCompress).DecompressData(data[0], &h)
	if err != nil {
		return nil, err
	}

	return &h, nil
}

func (c *chatRepository) DbInsertNewChat(ctx context.Context, userId int, botId int) {
	marshal, _ := jsoniter.Marshal(dao.DefaultData)
	chat := &domain.Chat{
		UserId:         userId,
		BotId:          botId,
		Title:          dao.DefaultTitle,
		LastUpdateTime: time.Now().Unix(),
		IsDelete:       false,
		Data:           marshal,
	}
	if err := c.mysql.Gorm().Table("chat_re").Create(chat).Error; err != nil {
		//TODO 异步 写入日志
	}
	return
}

func NewChatRepository(dbs *bootstrap.Databases) domain.ChatRepository {
	return &chatRepository{redis: dbs.Redis, mysql: dbs.Mysql}
}
