package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"mini-gpt/controller"
	"mini-gpt/setting"
	"net/http"
)

// SetupRouter 布尔型的配置项，控制应用程序的发布（release）模式。
// 如果setting.Conf.Release为true， 则表示应用程序正在发布模式下运行，否则为开发模式。
// 当应用程序处于发布模式时，通过gin.SetMode(gin.ReleaseMode)将Gin框架的运行模式设置为发布模式。
// 在发布模式下，Gin框架会关闭掉一些调试信息和中间件，以提高性能并减少潜在的安全风险。
func SetupRouter() *gin.Engine {
	//如果 Release为真 则启用发布模式
	if setting.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	//注册路由

	r.Use(cors.Default())
	// config := cors.DefaultConfig()
	// config.AllowAllOrigins = true
	// router.Use(cors.New(config))
	//此处注册跨域cors 中间件   默认配置

	////注册
	//r.POST("/register",
	//	controller.Register)
	////登录
	//r.POST("/login", controller.Login)

	//测试接口
	r.GET("/test", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"lang": "Golang",
			"tag":  "<br>",
		})
	})

	//r.GET("/", controller.IndexHandler)

	//原始请求格式 暂定如下
	r.POST("/chat/new", controller.CreateChat)

	mainPageGroup := r.Group("/init")
	{
		//主页面查询的chat历史记录
		mainPageGroup.GET("/", controller.InitChatHistory)
		//查询特定chat的历史记录
		mainPageGroup.GET("/:chatId", controller.GetChatHistory)
		//密钥分享chat记录
		mainPageGroup.GET("/share/:chatId/:ddl", controller.ShareHistory)
	}

	//主页面模型
	contextGroup := r.Group("/context")
	{
		//初始化上下文chat 返回一个id给客户端
		contextGroup.POST("/init", controller.InitNewChat)
		//真正调用gpt模型进行上下文交流
		contextGroup.POST("/call", controller.CallContextChat)
	}

	//小机器人功能
	botGroup := r.Group("/bot")
	{
		botGroup.POST("/", controller.CallBot)
	}

	//管理员功能
	adminGroup := r.Group("/admin")
	{
		adminBotGroup := adminGroup.Group("/bot")
		{
			//获取机器人信息及其提示词
			adminBotGroup.GET("/:isOfficial/:botId", controller.AdminGetBot)
			//管理员设置新机器人
			adminBotGroup.POST("/", controller.AdminSaveNewBot)
			//管理员更新现有机器人
			adminBotGroup.PUT("/", controller.AdminModifyBot)
		}
	}
	return r
}
