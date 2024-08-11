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

type botRepository struct {
	redis redis.Client
	mysql mysql.Client
}

func (b *botRepository) CacheGetBotHistory(ctx context.Context, chatId int) *[]*domain.Record {
	var a []*domain.Record
	err := b.redis.GetStruct(ctx, cache.ChatHistory+strconv.Itoa(chatId), a)
	//TODO 写历史记录的时候 记得滑动窗口维护最大条数
	if err != nil {
		return nil
	}
	return &a
}

func (b *botRepository) CacheGetBotConfig(ctx context.Context, botId int) *domain.BotConfig {
	var a domain.BotConfig
	err := b.redis.GetStruct(ctx, cache.BotConfig+strconv.Itoa(botId), a)
	if err != nil {
		return nil
	}
	return &a
}

// CacheGetMaxBotId 获取缓存合法bot最大值 用于判断数据是否合法
func (b *botRepository) CacheGetMaxBotId(ctx context.Context) int {
	maxId, err := b.redis.Get(ctx, cache.MaxBotId)
	if err != nil {
		return common.FalseInt
	}
	m, _ := strconv.Atoi(maxId)
	return m
}

func NewBotRepository() domain.BotRepository {
	return &botRepository{redis: dbs.Redis, mysql: dbs.Mysql}
}
