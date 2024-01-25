package setting

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gopkg.in/ini.v1"
)

// Conf 是一个包含应用程序配置的变量
var Conf = new(AppConfig)

func Init(file string) error {
	return ini.MapTo(Conf, file)
}

// AppConfig 应用程序配置
type AppConfig struct {
	Release       bool `ini:"release"`
	Port          int  `ini:"port"`
	*MySQLConfig  `ini:"mysql"`
	*LogrusConfig `ini:"logrus"`
	*ApiConfig    `ini:"api"`
	*RedisConfig  `ini:"redis"`
}

type ApiConfig struct {
	SecretKey string `ini:"secretKey"`
}

// MySQLConfig mysql配置
type MySQLConfig struct {
	User     string `ini:"user"`
	Password string `ini:"password"`
	DB       string `ini:"db"`
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
}

// RedisConfig redis配置
type RedisConfig struct {
	Client   *redis.Client
	Ctx      context.Context
	Addr     string `ini:"addr"`
	Password string `ini:"password"`
}

// LogrusConfig Logrus日志框架配置
type LogrusConfig struct {
	File string `ini:"file"`
}
