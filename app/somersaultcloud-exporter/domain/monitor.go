package domain

import "SomersaultCloud/app/somersaultcloud-common/discovery"

type MonitorStatus struct {
	EndpointInfo    *discovery.EndpointInfo
	ServiceRegister *discovery.ServiceRegister
}

type Monitor interface {
	HandleMonit(serviceName string)
	ServiceRegister()
}
