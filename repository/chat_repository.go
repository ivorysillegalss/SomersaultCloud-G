package repository

import (
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/database/mysql"
	"SomersaultCloud/database/redis"
	"SomersaultCloud/domain"
	"context"
	"strconv"
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
	// TODO MQ异步写入sql lua改进
}

func NewChatRepository() domain.ChatRepository {
	return &chatRepository{redis: dbs.Redis, mysql: dbs.Mysql}
}
