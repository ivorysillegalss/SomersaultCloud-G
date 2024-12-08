package domain

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"time"
)

type MonitorStatus struct {
	EndpointInfo    *discovery.EndpointInfo
	ServiceRegister *discovery.ServiceRegister
	Time            time.Time
}

type Monitor interface {
	HandleMonit()
	ServiceRegister()
	//TODO Prometheus消息抓取接口
}
