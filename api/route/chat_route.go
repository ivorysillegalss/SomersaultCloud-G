package route

import (
	"SomersaultCloud/api/controller"
	"github.com/gin-gonic/gin"
)

func RegisterChatRouter(group *gin.RouterGroup, cc *controller.ChatController) {
	chatGroup := group.Group("/chat")
	chatGroup.POST("/init", cc.InitNewChat)
}
