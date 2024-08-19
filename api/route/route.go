package route

import (
	"SomersaultCloud/bootstrap"
	"github.com/gin-gonic/gin"
)

func Setup(c *bootstrap.Controllers, e bootstrap.Executor) *gin.Engine {
	r := gin.Default()

	publicRouter := r.Group("")
	// All Public APIs
	RegisterChatRouter(publicRouter, c.ChatController)

	//Cron start
	e.SetupCron()
	return r
}
