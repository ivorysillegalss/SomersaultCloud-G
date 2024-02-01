package main

import (
	"fmt"
	"mini-gpt/dao"
	"mini-gpt/routers"
	"mini-gpt/setting"
)

func main() {
	//加载配置文件
	setting.AboutConf()

	//加载日志框架
	logger := setting.GetLogger()

	// 连接数据库
	//mysql
	mysqlErr := dao.InitMySQL(setting.Conf.MySQLConfig)
	//redisUtils
	redisErr := dao.InitRedis(setting.Conf.RedisConfig)
	if redisErr != nil || mysqlErr != nil {
		logger.Error("init database failed, err:%v\n")
		return
	}
	// 注册  程序退出关闭数据库连接
	defer func() {
		err := dao.Close()
		if err != nil {
			panic(err)
		}
	}()

	//上方代码顺序不能改变 日志框架文件路径在配置文件中 数据库初始化中使用了日志框架

	// 模型绑定 service层 此处仅测试
	//dao.DB.AutoMigrate(&models.Todo{})

	// 注册路由
	r := routers.SetupRouter()
	if err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port)); err != nil {
		logger.Error(err)
	}
}
