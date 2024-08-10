package route

import (
	"time"

	"SomersaultCloud/api/middleware"
	"SomersaultCloud/bootstrap"
	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db bootstrap.Databases, gin *gin.Engine) {
	publicRouter := gin.Group("")
	// All Public APIs
	NewChatRouter(publicRouter)

	protectedRouter := gin.Group("")

	// Middleware to verify AccessToken
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))

	// All Private APIs
	//NewProfileRouter(env, timeout, db, protectedRouter)
	//NewTaskRouter(env, timeout, db, protectedRouter)
}
