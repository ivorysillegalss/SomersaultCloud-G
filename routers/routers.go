package routers

import (
	"github.com/gin-gonic/gin"
	"mini-gpt/setting"
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

	//测试接口
	r.GET("/test", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"lang": "Golang",
			"tag":  "<br>",
		})
	})

	//r.GET("/", controller.IndexHandler)

	// v1
	//v1Group := r.Group("v1")
	//{
	//	// crud
	//	v1Group.POST("/todo", controller.CreateTodo)
	//	v1Group.GET("/todo", controller.GetTodoList)
	//	v1Group.PUT("/todo/:id", controller.UpdateATodo)
	//	v1Group.DELETE("/todo/:id", controller.DeleteATodo)
	//}
	return r
}
