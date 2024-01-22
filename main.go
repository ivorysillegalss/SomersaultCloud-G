package main

import (
	"fmt"
	"mini-gpt/dao"
	"mini-gpt/routers"
	"mini-gpt/setting"
)

func main() {
	setting.AboutConf()
	logger := setting.GetLogger()
	// 连接数据库
	//err := dao.InitMySQL(Conf.MySQLConfig)
	err := dao.InitMySQL(setting.Conf.MySQLConfig)
	if err != nil {
		//Logger.Error("init mysql failed, err:%v\n", err)
		logger.Error("init mysql failed, err:%v\n", err)
		return
	}
	defer dao.Close() // 程序退出关闭数据库连接

	//上方代码顺序不能改变 日志框架文件路径在配置文件中 数据库初始化中使用了日志框架

	// 模型绑定 service层 此处仅测试
	//dao.DB.AutoMigrate(&models.Todo{})

	// 注册路由
	r := routers.SetupRouter()
	if err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port)); err != nil {
		logger.Error(err)
	}
}
