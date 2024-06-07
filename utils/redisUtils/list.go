package redisUtils

import (
	"mini-gpt/dao"
	"time"
)

func Rpush(k string, v ...any) error {
	return dao.Client.RPush(dao.Ctx, k, v).Err()
}

func GetList(k string) ([]string, error) {
	return dao.Client.LRange(dao.Ctx, k, 0, -1).Result()
}

func Expire(k string, t time.Duration) error {
	return dao.Client.Expire(dao.Ctx, k, t).Err()
}

// 删除k中值为v的第一个元素
func LRemFirst(k string, v any) error {
	return dao.Client.LRem(dao.Ctx, k, 1, v).Err()
}
