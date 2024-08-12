package repository

import (
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/db"
	"SomersaultCloud/database/mysql"
	"SomersaultCloud/database/redis"
	"SomersaultCloud/domain"
	"context"
	"encoding/json"
	"strconv"
	"time"
)

// 多例是正确的吗？
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
	return res.(int), nil
}

func (c *chatRepository) CacheGetHistory(ctx context.Context, chatId int) (history *[]*domain.Record, isCache bool, isErr error) {
	var h []*domain.Record
	err := c.redis.GetStruct(ctx, strconv.Itoa(chatId), h)
	if c.redis.IsEmpty(err) {
		return nil, true, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &h, false, nil
}

func (c *chatRepository) CacheLuaLruPutHistory(ctx context.Context) {
	//c.redis.Lru(ctx, cache.ContextLruMaxCapacity, cache.ListType).Add(ctx)
	//TODO
}

// DbGetHistory TODO db中数据流式更新flink 更新入Hbase等
//
//	目前的架构一旦db获取历史记录 就是全部获取 初步思路是定时任务 mq以某个时间段为界限（eg7天） 将数据流式更新
//	保证不出现大Key等 主要是为了提高查询效率
func (c *chatRepository) DbGetHistory(ctx context.Context, chatId int) (*[]*domain.Record, error) {
	var h []*domain.Record
	var data string

	//pluck方法在不构建结构体的前提下 获取单个字段的值
	if err := c.mysql.Gorm().Table("chat_re").Where("chat_id = ?", chatId).Pluck("data", data).Error; err != nil {
		return nil, err
	}

	err := json.Unmarshal([]byte(data), &h)
	if err != nil {
		return nil, err
	}

	return &h, nil
}

func (c *chatRepository) DbInsertNewChatId(ctx context.Context, userId int, botId int) {
	chat := &domain.Chat{
		UserId:         userId,
		BotId:          botId,
		Title:          db.DefaultTitle,
		LastUpdateTime: time.Now().Unix(),
		IsDelete:       false,
	}
	if err := c.mysql.Gorm().Table("chat_re").Create(chat).Error; err != nil {
		//TODO 异步 写入日志
	}
	return
}

func NewChatRepository() domain.ChatRepository {
	return &chatRepository{redis: dbs.Redis, mysql: dbs.Mysql}
}
