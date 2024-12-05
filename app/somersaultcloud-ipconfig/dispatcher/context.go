package dispatcher

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

// IpConfContext ipconfig中可能包含使用到的所有上下文封装
type IpConfContext struct {
	ctx       context.Context
	AppCtx    *app.RequestContext // Hertz的req上下文
	ClientCtx *ClientContext
}

type ClientContext struct {
	IP string `json:"ip"`
}

func BuildIpConfContext(c context.Context, ctx *app.RequestContext) *IpConfContext {
	return &IpConfContext{
		ctx:       c,
		AppCtx:    ctx,
		ClientCtx: &ClientContext{},
	}
}
