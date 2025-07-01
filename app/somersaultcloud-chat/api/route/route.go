package route

import (
	"SomersaultCloud/app/somersaultcloud-chat/api/middleware"
	"SomersaultCloud/app/somersaultcloud-chat/bootstrap"
	"SomersaultCloud/app/somersaultcloud-chat/infrastructure/redis"
	"github.com/cloudwego/hertz/pkg/app/server"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
)

func Setup(env *bootstrap.Env, c *bootstrap.Controllers, e *bootstrap.Executor, rcl redis.Client) *server.Hertz {
	r := server.Default(server.WithHostPorts(env.ServerAddress),
		server.WithTracer(
			prometheus.NewServerTracer(env.Prometheus.ServerAddress, "/trace"),
		))

	publicRouter := r.Group("/domain")
	// All Public APIs
	RegisterChatRouter(publicRouter, c)

	//注册跨域中间件
	middleware.DefaultCorsConfig(r)
	//注册限流中间件 Redis实现
	middleware.RegisterRateLimit(r, rcl)

	//Cron start
	e.CronExecutor.SetupCron()
	//Consume Start
	e.ConsumeExecutor.SetupConsume()
	//InitRedis start
	e.DataExecutor.InitData()

	return r
}
