package route

import (
	"SomersaultCloud/app/somersaultcloud-chat/bootstrap"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
)

func Setup(env *bootstrap.Env, c *bootstrap.Controllers, e *bootstrap.Executor) *server.Hertz {
	r := server.Default(server.WithHostPorts(env.ServerAddress),
		server.WithTracer(
			prometheus.NewServerTracer(env.Prometheus.ServerAddress, "/trace"),
		))
	defaultCorsConfig(r)

	publicRouter := r.Group("/domain")
	// All Public APIs
	RegisterChatRouter(publicRouter, c)

	//Cron start
	e.CronExecutor.SetupCron()
	//Consume Start
	e.ConsumeExecutor.SetupConsume()
	//InitRedis start
	e.DataExecutor.InitData()

	return r
}

// CORS middleware configuration for Hertz
func defaultCorsConfig(h *server.Hertz) {
	// Correct middleware signature to match app.HandlerFunc
	h.Use(
		func(ctx context.Context, c *app.RequestContext) {
			// Set allowed origins, "*" means all domains are allowed
			c.Response.Header.Set("Access-Control-Allow-Origin", "*")
			// Set allowed methods
			c.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			// Set allowed headers, include your custom headers like "token"
			c.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, token")
			// Set whether credentials should be included (optional)
			// c.Response.Header.Set("Access-Control-Allow-Credentials", "true")

			// Handle OPTIONS request (CORS preflight request)

			if consts.MethodOptions == string(c.Request.Method()[0]) {
				c.AbortWithStatus(204)
			} else {
				c.Next(ctx)
			}
		})
}
