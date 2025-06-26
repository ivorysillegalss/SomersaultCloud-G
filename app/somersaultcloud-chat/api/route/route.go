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
	registerCorsConfig(r, env)
	//注册限流中间件 Redis实现
	registerRateLimit(r, rcl, env.Net.RateLimit)

	//Cron start
	e.CronExecutor.SetupCron()
	//Consume Start
	e.ConsumeExecutor.SetupConsume()
	//InitRedis start
	e.DataExecutor.InitData()

	return r
}

// 注册多种限流
func registerRateLimit(h *server.Hertz, r redis.Client, rl bootstrap.RateLimit) {
	h.Use(middleware.BuckRateLimit(r, rl.Buck.Rate, rl.Buck.Capacity, rl.Buck.Prefix, rl.Buck.Requested))
	h.Use()
}

// 注册跨域配置
func registerCorsConfig(h *server.Hertz, env *bootstrap.Env) {
	cors := env.Net.Cors
	if cors.Default {
		h.Use(middleware.DefaultCorsConfig())
	} else {
		h.Use(middleware.CustomCorsConfig(cors.AllowAllOrigin, cors.AllowAllCredentials, cors.Headers, cors.Methods))
	}
}
