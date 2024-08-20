package lru

import (
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/infrastructure/redis"
	"SomersaultCloud/internal/ioutil"
	"context"
)

func NewLru(maxCapacity int, dataType int, middleware any) Lru {
	var redisClient *redis.Client
	if client, ok := middleware.(*redis.Client); ok {
		redisClient = client
	}

	if dataType == cache.RedisListType {
		return &redisLuaLruList{maxCapacity: maxCapacity, rcl: *redisClient}
	}
	if dataType == cache.RedisZSetType {
		return &redisLuaLruZSet{maxCapacity: maxCapacity, rcl: *redisClient}
	}

	//TODO 丰富实现的数据类型 eg InMemoryLru
	panic("error")
}

// Lru 接口定义
type Lru interface {
	List(ctx context.Context, k string) ([]string, error)
	Add(ctx context.Context, k, value string) error
	Remove(ctx context.Context, k string, v string) error
	isExist(ctx context.Context, k string, member string) bool
	//Get(ctx context.Context, field string) (string, error)
}

// ZSet实现类型
type redisLuaLruZSet struct {
	rcl         redis.Client
	maxCapacity int
}

func (r *redisLuaLruZSet) Remove(ctx context.Context, k string, v string) error {
	_, err := r.rcl.ZRem(ctx, k, v)
	return err
}

func (r *redisLuaLruZSet) isExist(ctx context.Context, k string, member string) bool {
	isExist, _ := r.rcl.ZScore(ctx, k, member)
	return isExist
}

func (r *redisLuaLruZSet) List(ctx context.Context, k string) ([]string, error) {
	return r.rcl.LRange(ctx, k, 0, -1)
}

func (r *redisLuaLruZSet) Add(ctx context.Context, key, value string) error {
	luaScript, err := ioutil.LoadLuaScript("lua/zsetlru.lua")
	err = r.rcl.ExecuteArgsLuaScript(ctx, luaScript, []string{key, key + common.Infix + cache.LruPrefix}, value, r.maxCapacity)
	return err
}

// redisLuaLruList List实现类型
type redisLuaLruList struct {
	rcl         redis.Client
	maxCapacity int
}

func (r *redisLuaLruList) Remove(ctx context.Context, k string, v string) error {
	//Not Recommend 不推荐使用

	// 使用 LREM 命令移除列表中的指定元素
	_, err := r.rcl.LRem(ctx, k, 0, v)
	if err != nil {
		return err
	}

	// 同时移除 LRU 集合中的该元素
	_, err = r.rcl.ZRem(ctx, k+":lru", v)
	if err != nil {
		return err
	}

	return nil
}

func (r *redisLuaLruList) isExist(ctx context.Context, k string, v string) bool {
	//Not Recommend 不推荐使用
	// 获取列表的所有元素
	list, _ := r.rcl.LRange(ctx, k, 0, -1)

	// 遍历列表，检查是否存在目标值
	for _, value := range list {
		if value == v {
			return true
		}
	}

	return false
}

func (r *redisLuaLruList) List(ctx context.Context, k string) ([]string, error) {
	return r.rcl.LRange(ctx, k, 0, -1)
}

func (r *redisLuaLruList) Add(ctx context.Context, key, value string) error {
	luaScript, err := ioutil.LoadLuaScript("lua/listlru.lua")
	err = r.rcl.ExecuteArgsLuaScript(ctx, luaScript, []string{key, key + common.Infix + cache.LruPrefix}, value, r.maxCapacity)
	return err
}
