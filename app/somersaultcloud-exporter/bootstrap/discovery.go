package bootstrap

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"context"
)

func NewServiceDiscovery(env *ExporterEnv) discovery.ServiceDiscovery {
	config := env.DiscoveryConfig
	return discovery.NewServiceDiscovery(context.Background(), config.Endpoints, config.Timeout,
		config.Username, config.Password)
}
