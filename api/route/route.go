package route

import (
	"SomersaultCloud/bootstrap"
	"github.com/gin-gonic/gin"
)

func Setup(c *bootstrap.Controllers, e *bootstrap.Executor) *gin.Engine {
	r := gin.Default()
	defaultCorsConfig(r)

	publicRouter := r.Group("")
	// All Public APIs
	RegisterChatRouter(publicRouter, c)

	//Cron start
	e.CronExecutor.SetupCron()
	//Consume Start
	e.ConsumeExecutor.SetupConsume()

	return r
}

// CORS 中间件配置
func defaultCorsConfig(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		// 设置允许访问的源，"*" 表示允许所有源，你也可以指定具体的域名
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 设置允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		// 设置允许的头部，注意添加你的自定义头部 "token"
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, token")
		// 设置浏览器是否应该包含凭证
		//c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 如果是OPTIONS请求，直接返回200
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
		} else {
			c.Next()
		}
	})
}
