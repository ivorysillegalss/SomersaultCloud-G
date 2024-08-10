package repository

import (
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/db"
	"SomersaultCloud/database/mysql"
	"SomersaultCloud/database/redis"
	"SomersaultCloud/domain"
	"context"
	"strconv"
	"time"
)

// 多例是正确的吗？
type chatRepository struct {
	redis redis.Client
	mysql mysql.Client
}

// CacheGetNewestChatId 获取最新chatId 不能保证原子性 弃用
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

// CacheInsertNewChat 增加新Id 不能保证原子性 弃用
func (c *chatRepository) CacheInsertNewChat(ctx context.Context, id int) {
	_ = c.redis.Set(ctx, cache.NewestChatIdKey, id+1)
	// MQ异步写入sql lua改进 暂时无用
}

// CacheLuaInsertNewChatId lua脚本保证高并发时获取chatId的一致性
func (c *chatRepository) CacheLuaInsertNewChatId(ctx context.Context, luaScript string, k string) (int, error) {
	res, err := c.redis.ExecuteLuaScript(ctx, luaScript, k)
	if err != nil {
		return common.FalseInt, err
	}
	return res.(int), nil
}

// DbInsertNewChatId 异步使用 存入SQL持久化方法
func (c *chatRepository) DbInsertNewChatId(ctx context.Context, userId int, botId int) {
	chat := &domain.Chat{
		UserId:         userId,
		BotId:          botId,
		Title:          db.DefaultTitle,
		LastUpdateTime: time.Now().Unix(),
		IsDelete:       false,
	}
	if err := c.mysql.Gorm().Table("chat").Create(chat).Error; err != nil {
		//TODO 异步 写入日志
	}
	return
}

func NewChatRepository() domain.ChatRepository {
	return &chatRepository{redis: dbs.Redis, mysql: dbs.Mysql}
}
