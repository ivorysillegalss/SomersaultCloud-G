package routers

import (
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

	// CORS 中间件配置
	r.Use(func(c *gin.Context) {
		// 设置允许访问的源，"*" 表示允许所有源，你也可以指定具体的域名
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 设置允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		// 设置允许的头部，注意添加你的自定义头部 "token"
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, token")
		// 设置浏览器是否应该包含凭证
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 如果是OPTIONS请求，直接返回200
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
		} else {
			c.Next()
		}
	})

	////注册
	//r.POST("/register",
	//	controller.Register)
	////登录
	//r.POST("/login", controller.Login)

	//
	//r.OPTIONS("/", func(context *gin.Context) {
	//	context.
	//})

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
	}

	//历史记录相关的接口 目前只有分享
	//收藏重命名 TBD
	historyGroup := r.Group("/history")
	{
		historyShareGroup := historyGroup.Group("/share")
		{
			//密钥分享chat记录
			historyShareGroup.GET("/:chatId", controller.ShareHistoryWithSk)
			//密钥获取chat记录 预览
			historyShareGroup.GET("/get/:sk", controller.GetSharedHistoryWithSk)
			//依据分享所得chat记录 继续聊天
			historyShareGroup.POST("/get/:sk", controller.ContinueSharedChat)
		}
	}

	//主页面模型
	contextGroup := r.Group("/context")
	{
		//初始化上下文chat 返回一个id给客户端
		contextGroup.POST("/init", controller.InitNewChat)
		//真正调用gpt模型进行上下文交流
		contextGroup.POST("/call", controller.CallContextChat)
		//在chat一次之后 根据已有的历史记录获取一个标题
		contextGroup.POST("/title/init", controller.InitialTitle)
		//根据用户的输入更新标题
		contextGroup.POST("/title", controller.UpdateTitle)
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

	//RPC接口
	rpcGroup := r.Group("/rpc")
	{
		//接受一段json数据 发起api调用并返回对应的标题作为字符串
		rpcGroup.POST("/title", controller.Rpc4Title)
	}
	return r
}
