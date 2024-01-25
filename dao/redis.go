package dao

import (
	"context"
	"github.com/redis/go-redis/v9"
	exception "mini-gpt/error"
	"mini-gpt/setting"
	"time"
)

var (
	// UserAddr 用户连接信息
	UserAddr string
	// Client 客户端
	Client *redis.Client
	// Ctx 连接上下文
	Ctx context.Context
)

func InitRedis(rfg *setting.RedisConfig) error {
	UserAddr = rfg.Addr
	rdb := redis.NewClient(&redis.Options{
		Addr:     rfg.Addr,
		Password: rfg.Password,
	})
	Client = rdb
	Ctx = context.Background()
	//赋值到全局变量
	ctx2 := context.Background()
	redisTemplate := &setting.RedisConfig{
		Client: rdb,
		Ctx:    ctx2,
	}
	//验证连接是否正常
	isConn := Ping(redisTemplate)
	if isConn != nil {
		return isConn
	}
	return nil
}

func Ping(redisTemplate *setting.RedisConfig) error {
	ctx := redisTemplate.Ctx
	_, err := redisTemplate.Client.Ping(ctx).Result()
	if err != nil {
		connError := exception.ConnError{
			ConnTime: time.Now(),
			Addr:     UserAddr,
		}
		return &connError
	}
	return nil
}
