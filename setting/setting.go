package setting

import "gopkg.in/ini.v1"

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
}

type ApiConfig struct {
	SecretKey string `ini:"secretKey"`
}

// MySQLConfig 数据库配置
type MySQLConfig struct {
	User     string `ini:"user"`
	Password string `ini:"password"`
	DB       string `ini:"db"`
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
}

// LogrusConfig Logrus日志框架配置
type LogrusConfig struct {
	File string `ini:"file"`
}
