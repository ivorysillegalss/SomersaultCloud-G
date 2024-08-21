package repository

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/mysql"
	"SomersaultCloud/infrastructure/redis"
	"context"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

type botRepository struct {
	redis redis.Client
	mysql mysql.Client
}

func (b *botRepository) CacheGetBotHistory(ctx context.Context, chatId int) *[]*domain.Record {
	var a []*domain.Record
	err := b.redis.GetStruct(ctx, cache.ChatHistory+strconv.Itoa(chatId), a)
	if err != nil {
		return nil
	}
	return &a
}

func (b *botRepository) CacheGetBotConfig(ctx context.Context, botId int) *domain.BotConfig {
	var a domain.BotConfig
	//err := b.redis.GetStruct(ctx, cache.BotConfig+strconv.Itoa(botId), a)
	get, err := b.redis.Get(ctx, cache.BotConfig+strconv.Itoa(botId))
	if err != nil {
		return nil
	}
	err = jsoniter.Unmarshal([]byte(get), &a)
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

func NewBotRepository(dbs *bootstrap.Databases) domain.BotRepository {
	return &botRepository{redis: dbs.Redis, mysql: dbs.Mysql}
}
