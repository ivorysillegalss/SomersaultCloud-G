package domain

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"time"
)

type MonitorStatus struct {
	EndpointInfo    *discovery.EndpointInfo
	ServiceRegister *discovery.ServiceRegister
	Time            time.Time
}

type Monitor interface {
	// ServiceRegister 初始化服务注册
	ServiceRegister()

	// HandleMonit 服务发现推送到etcd
	HandleMonit()

	// ExposeMonitorInterface 暴露接口供Prometheus拉取
	ExposeMonitorInterface(c context.Context, ctx *app.RequestContext)
}
