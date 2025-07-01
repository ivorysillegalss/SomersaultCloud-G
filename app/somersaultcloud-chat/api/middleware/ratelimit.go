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
	"github.com/cloudwego/hertz/pkg/app/server"
	"net/http"
	"time"
)

//go:embed lua/ratelimit.lua
var rateLimitRateLimit string
var rcl redis.Client
var chatConfig rateLimitConfig

const (
	chatRateLimitPrefix = "chat_rate_limit"
	chatRate            = "5"
	chatCapacity        = "10"

	commonRequested = "1"
)

// 无论是什么类型 传到lua里都需要是string
type rateLimitConfig struct {
	prefix    string
	rate      string
	capacity  string
	requested string
}

// TODO 多起来限流条件可以整理为一个map
func init() {
	chatConfig = rateLimitConfig{
		prefix:    chatRateLimitPrefix,
		rate:      chatRate,
		capacity:  chatCapacity,
		requested: commonRequested,
	}
}

func RegisterRateLimit(h *server.Hertz, r redis.Client) {
	rcl = r
	h.Use(buckRateLimit)
}

func buckRateLimit(c context.Context, ctx *app.RequestContext) {
	handle(c, ctx, chatConfig)
}

func handle(c context.Context, ctx *app.RequestContext, conf rateLimitConfig) {
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
