package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Client interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, k string, v any) error
	SetExpire(ctx context.Context, k string, v any, ddl time.Duration) error
	Get(ctx context.Context, k string) (string, error)
	SetStruct(ctx context.Context, k string, vStruct any) error
	SetStructExpire(ctx context.Context, k string, vStruct any, ddl time.Duration) error
	GetStruct(ctx context.Context, k string, targetStruct any) error
	ExecuteLuaScript(ctx context.Context, luaScript string, k string) (any, error)
	IsEmpty(err error) bool
}

type redisClient struct {
	rcl *redis.Client
}

func (r *redisClient) Ping(ctx context.Context) error {
	_, err := r.rcl.Ping(ctx).Result()
	return err
}

func (r *redisClient) Set(ctx context.Context, k string, v any) error {
	return r.rcl.Set(ctx, k, v, 0).Err()
}

func (r *redisClient) SetExpire(ctx context.Context, k string, v any, ddl time.Duration) error {
	return r.rcl.Set(ctx, k, v, ddl).Err()
}

func (r *redisClient) Get(ctx context.Context, k string) (string, error) {
	return r.rcl.Get(ctx, k).Result()
}

func (r *redisClient) SetStruct(ctx context.Context, k string, vStruct any) error {
	vJsonData, _ := json.Marshal(vStruct)
	return r.Set(ctx, k, vJsonData)
}

func (r *redisClient) SetStructExpire(ctx context.Context, k string, vStruct any, ddl time.Duration) error {
	vJsonData, _ := json.Marshal(vStruct)
	return r.SetExpire(ctx, k, vJsonData, ddl)
}

// GetStruct 获取自定义结构体
func (r *redisClient) GetStruct(ctx context.Context, k string, targetStruct any) error {
	vJsonData, err := r.rcl.Get(ctx, k).Result() // 获取存储的 JSON 字符串
	if err != nil {
		return err
	}
	// 将 JSON 字符串反序列化为结构体
	err = json.Unmarshal([]byte(vJsonData), targetStruct)
	if err != nil {
		return err
	}
	return nil
}

// ExecuteLuaScript 执行lua脚本 保证操作原子性
func (r *redisClient) ExecuteLuaScript(ctx context.Context, luaScript string, k string) (any, error) {
	result, err := r.rcl.Eval(ctx, luaScript, []string{k}).Result()
	if err != nil {
		return nil, err
	}
	return result, err
}

func (r *redisClient) IsEmpty(err error) bool {
	return errors.Is(err, redis.Nil)
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
