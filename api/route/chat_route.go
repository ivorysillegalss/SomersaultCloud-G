package route

import (
	"SomersaultCloud/api/controller"
	"github.com/gin-gonic/gin"
)

func NewChatRouter(group *gin.RouterGroup) {
	chatController := controller.NewChatController()
	chatGroup := group.Group("/chat")
	chatGroup.POST("/init", chatController.InitNewChat)
}
