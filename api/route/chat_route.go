package route

import (
	"SomersaultCloud/bootstrap"
	"github.com/gin-gonic/gin"
)

func RegisterChatRouter(group *gin.RouterGroup, controllers *bootstrap.Controllers) {
	cc := controllers.ChatController
	hmc := controllers.HistoryMessageController
	chatGroup := group.Group("/context")
	{
		//开启新chat
		chatGroup.POST("/init", cc.InitNewChat)
		//启动上下文chat
		chatGroup.POST("/call", cc.ContextTextChat)
		//在chat一次之后 根据已有的历史记录获取一个标题
		chatGroup.POST("/title/init", hmc.UpdateInitTitle)
		//根据用户的输入更新标题
		chatGroup.POST("/title", hmc.InputTitle)
	}

	streamGroup := chatGroup.Group("/stream")
	{
		streamGroup.POST("/setup", cc.StreamContextTextChatSetup)
		streamGroup.POST("/work", cc.StreamContextTextChatWorker)
	}

	mc := controllers.HistoryMessageController
	mainPageGroup := group.Group("/init")
	{
		//主页面查询的chat历史记录
		mainPageGroup.GET("/title/:botId", mc.HistoryTitle)
		//查询特定chat的历史记录
		mainPageGroup.GET("/:chatId/:botId", mc.GetChatHistory)
	}
}
