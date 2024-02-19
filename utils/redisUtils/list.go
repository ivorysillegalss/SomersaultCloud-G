package redisUtils

import "mini-gpt/dao"

func Rpush(k string, v ...any) error {
	return dao.Client.RPush(dao.Ctx, k, v).Err()
}

func GetList(k string) ([]string, error) {
	return dao.Client.LRange(dao.Ctx, k, 0, -1).Result()
}
