package setting

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

const defaultConfFile = "./conf/config.ini"

var Logger logrus.Logger

// AboutConf 配置文件相关
func AboutConf() {
	confFile := defaultConfFile

	//存在多个配置文件 挑一个
	//if len(os.Args) > 2 {
	//	fmt.Println("use specified conf file: ", os.Args[1])
	//	confFile = os.Args[1]
	//} else {
	//	fmt.Println("no configuration file was specified, use ../conf/config.ini")
	//	fmt.Println("\n")
	//}

	// 加载配置文件
	if err := Init(confFile); err != nil {
		fmt.Printf("load config from file failed, err:%v\n", err)
		return
	}
}

// initLog 日志框架相关
func initLog() {
	Logger = *logrus.New()
	//输出为json格式
	logrus.SetFormatter(new(logrus.JSONFormatter))
	//日志级别
	Logger.SetLevel(logrus.InfoLevel)
	//输出路径
	fileName := Conf.LogrusConfig.File
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.SetOutput(file)

	Logger.Warn("a")

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			Logger.Error(err)
		}
	}(file)

}

// GetLogger 单例模式获取日志框架对象
func GetLogger() *logrus.Logger {
	initLog()
	return &Logger
}
