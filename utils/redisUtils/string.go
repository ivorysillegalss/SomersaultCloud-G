package redisUtils

import (
	"encoding/json"
	"mini-gpt/dao"
	"reflect"
	"time"
)

//此go文件为redis常用函数二次封装

// 面向redis 自动化字段映射 将结构体利用反射转为map
func StructToMap(obj interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		data[typeOfS.Field(i).Name] = field.Interface()
	}

	return data
}

// String set 不设定时
func Set(k string, v any) error {
	//v只建议存入 byte[] 或 string
	return dao.Client.Set(dao.Ctx, k, v, 0).Err()
}

// String set 设置定时时间
func SetWithExpire(k string, v any, ddl time.Duration) error {
	return dao.Client.Set(dao.Ctx, k, v, ddl).Err()
}

// String set 结合json 放入结构体 不设定时
func SetStruct(k string, vStruct any) error {
	vJsonData, _ := json.Marshal(vStruct)
	return Set(k, vJsonData)
}

// String set 结合json 放入结构体 设置定时时间
func SetStructWithExpire(k string, vStruct any, ddl time.Duration) error {
	vJsonData, _ := json.Marshal(vStruct)
	return SetWithExpire(k, vJsonData, ddl)
}

// String get 获取值 不返回剩余时间
func Get(k string) (string, error) {
	v, err := dao.Client.Get(dao.Ctx, k).Result()
	//如果get不到需要的值 会返回redis.errNil
	return v, err
}

// string get 结合json泛型  获取自定义结构体
func GetStruct[T any](k string) (T, error) {
	var vStruct T
	v, err := Get(k)
	if err != nil {
		return vStruct, err
	}
	err2 := json.Unmarshal([]byte(v), &vStruct)
	if err2 != nil {
		return vStruct, err
	}
	return vStruct, nil
}

// 获取特定键值对剩余时间
func DDL(k string) (time.Duration, error) {
	return dao.Client.TTL(dao.Ctx, k).Result()
}

// 获取值的基础上返回其存活时间
func GetWithExpire(k string) (string, time.Duration, error) {
	v, err := Get(k)
	if err != nil {
		return "", 0, err
	}
	ddl, err := DDL(k)
	if err != nil {
		return "", 0, err
	}
	return v, ddl, nil
}

func Del(k ...string) error {
	return dao.Client.Del(dao.Ctx, k...).Err()
}

func TTL(k string) error {
	return dao.Client.TTL(dao.Ctx, k).Err()
}
