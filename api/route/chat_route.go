package route

import (
	"SomersaultCloud/bootstrap"
	"github.com/gin-gonic/gin"
)

func RegisterChatRouter(group *gin.RouterGroup, controllers *bootstrap.Controllers) {
	cc := controllers.ChatController
	chatGroup := group.Group("/context")
	{
		//开启新chat
		chatGroup.POST("/init", cc.InitNewChat)
		//启动上下文chat
		chatGroup.POST("/call", cc.ContextChat)
	}

	mc := controllers.HistoryMessageController
	mainPageGroup := group.Group("/init")
	{
		//主页面查询的chat历史记录
		mainPageGroup.GET("/", mc.HistoryTitle)
		//查询特定chat的历史记录
		mainPageGroup.GET("/:chatId", mc.GetChatHistory)
	}
}
