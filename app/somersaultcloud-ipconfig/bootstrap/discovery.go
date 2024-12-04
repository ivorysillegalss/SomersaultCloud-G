package bootstrap

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"SomersaultCloud/app/somersaultcloud-ipconfig/dispatcher"
	"SomersaultCloud/app/somersaultcloud-ipconfig/source"
	"context"
	"time"
)

func NewServiceDiscovery(env *IpConfigEnv) discovery.ServiceDiscovery {
	config := env.DiscoveryConfig
	return discovery.NewServiceDiscovery(context.Background(), config.Endpoints, config.Timeout*time.Second)
}

func NewDispatcher(env *IpConfigEnv) *dispatcher.Dispatcher {
	return &dispatcher.Dispatcher{IpConfigEnv: env}
}

func NewDataHandler(env *IpConfigEnv, dis discovery.ServiceDiscovery) *source.DataHandler {
	return &source.DataHandler{IpConfigEnv: env, ServiceDiscovery: dis}
}
