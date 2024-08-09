package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Client interface {
	Ping(ctx context.Context) error
}

type redisClient struct {
	rcl *redis.Client
}

func (r *redisClient) Ping(ctx context.Context) error {
	_, err := r.rcl.Ping(ctx).Result()
	return err
}

type InitRedisApplication struct {
	UserAddr string
	Password string
}

func NewRedisApplication(addr string, password string) *InitRedisApplication {
	return &InitRedisApplication{
		UserAddr: addr,
		Password: password,
	}
}

func NewRedisClient(r *InitRedisApplication) (Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     r.UserAddr,
		Password: r.Password,
	})

	return &redisClient{rcl: client}, nil
}
