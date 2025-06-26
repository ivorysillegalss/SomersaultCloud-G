package middleware

import (
	"SomersaultCloud/app/somersaultcloud-chat/constant/common"
	"SomersaultCloud/app/somersaultcloud-chat/constant/request"
	"SomersaultCloud/app/somersaultcloud-chat/domain"
	"SomersaultCloud/app/somersaultcloud-chat/infrastructure/redis"
	"SomersaultCloud/app/somersaultcloud-common/log"
	"context"
	_ "embed"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"time"
)

//go:embed lua/ratelimit.lua
var rateLimitRateLimit string
var rcl redis.Client
var chatConfig BuckRateLimitConfig

const (
	chatRateLimitPrefix = "chat_rate_limit"
	chatRate            = "5"
	chatCapacity        = "10"

	commonRequested = "1"
)

// BuckRateLimitConfig 无论是什么类型 传到lua里都需要是string
type BuckRateLimitConfig struct {
	prefix    string
	rate      string
	capacity  string
	requested string
}

// TODO 多起来限流条件可以整理为一个map
func init() {
	chatConfig = BuckRateLimitConfig{
		prefix:    chatRateLimitPrefix,
		rate:      chatRate,
		capacity:  chatCapacity,
		requested: commonRequested,
	}
}

// BuckRateLimit 闭包函数 将桶限流作为单个限流的方法独立出来
func BuckRateLimit(r redis.Client, rate string, capacity string, prefix string, requested string) app.HandlerFunc {
	r = rcl
	buckConfig := chatConfig
	if rate != "" {
		buckConfig.rate = rate
	}
	if capacity != "" {
		buckConfig.capacity = capacity
	}
	if prefix != "" {
		buckConfig.prefix = prefix
	}
	if requested != "" {
		buckConfig.requested = requested
	}
	return func(c context.Context, ctx *app.RequestContext) {
		handle(c, ctx, buckConfig)
	}
}

func handle(c context.Context, ctx *app.RequestContext, conf BuckRateLimitConfig) {
	tokenString := ctx.Request.Header.Get("token")
	prefix := conf.prefix + common.Infix + tokenString
	err, v := rcl.ExecuteArgsLuaScript(c, rateLimitRateLimit,
		[]string{prefix}, conf.rate, conf.capacity, time.Now().String(), conf.requested)
	if err != nil {
		log.GetTextLogger().Error("rate limit lua script execute error in token: %s", tokenString)
	}
	if v[0].(int) == common.True {
		ctx.Next(c)
	} else {
		ctx.JSON(http.StatusTooManyRequests, domain.ErrorResponse{Message: "请求次数过多", Code: request.RateLimit})
	}
}
