package redis

import (
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/internal/ioutil"
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
	LRange(ctx context.Context, k string, start int, end int) ([]string, error)
	SetStruct(ctx context.Context, k string, vStruct any) error
	SetStructExpire(ctx context.Context, k string, vStruct any, ddl time.Duration) error
	GetStruct(ctx context.Context, k string, targetStruct any) error
	ExecuteLuaScript(ctx context.Context, luaScript string, k string) (any, error)
	ExecuteArgsLuaScript(ctx context.Context, luaScript string, keys []string, args ...interface{}) error
	IsEmpty(err error) bool
	Lru(ctx context.Context, maxCapacity int, dataType int) Lru
}

type redisClient struct {
	rcl *redis.Client
}

func (r *redisClient) Lru(ctx context.Context, maxCapacity int, dataType int) Lru {
	if dataType == cache.ListType {
		return &redisLuaLruList{maxCapacity: maxCapacity}
	}
	//TODO 丰富实现的数据类型
	panic("error")
}

// Lru 接口定义
type Lru interface {
	List(ctx context.Context, k string) ([]string, error)
	Add(ctx context.Context, k, value string) error
	Remove(ctx context.Context, k string, v string) error
	isExist(ctx context.Context, k string, v string) (bool, error)
	//Get(ctx context.Context, field string) (string, error)
}

// redisLuaLruList 实现
type redisLuaLruList struct {
	rcl         *redis.Client
	r           Client
	maxCapacity int
}

func (r *redisLuaLruList) Remove(ctx context.Context, k string, v string) error {
	//Not Recommend 不推荐使用

	// 使用 LREM 命令移除列表中的指定元素
	_, err := r.rcl.LRem(ctx, k, 0, v).Result()
	if err != nil {
		return err
	}

	// 同时移除 LRU 集合中的该元素
	_, err = r.rcl.ZRem(ctx, k+":lru", v).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisLuaLruList) isExist(ctx context.Context, k string, v string) (bool, error) {
	//Not Recommend 不推荐使用
	// 获取列表的所有元素
	list, err := r.rcl.LRange(ctx, k, 0, -1).Result()
	if err != nil {
		return false, err
	}

	// 遍历列表，检查是否存在目标值
	for _, value := range list {
		if value == v {
			return true, nil
		}
	}

	return false, nil
}

func (r *redisLuaLruList) List(ctx context.Context, k string) ([]string, error) {
	return r.r.LRange(ctx, k, 0, -1)
}

func (r *redisLuaLruList) Add(ctx context.Context, key, value string) error {
	luaScript, err := ioutil.LoadLuaScript("lua/listlru.lua")
	//	TODO
	err = r.r.ExecuteArgsLuaScript(ctx, luaScript, []string{key, key + ":lru"}, value, r.maxCapacity)
	return err
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

func (r *redisClient) LRange(ctx context.Context, k string, start int, end int) ([]string, error) {
	return r.rcl.LRange(ctx, k, int64(start), int64(end)).Result()
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

func (r *redisClient) ExecuteArgsLuaScript(ctx context.Context, luaScript string, keys []string, args ...interface{}) error {
	_, err := r.rcl.Eval(ctx, luaScript, keys, args).Result()
	return err
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
