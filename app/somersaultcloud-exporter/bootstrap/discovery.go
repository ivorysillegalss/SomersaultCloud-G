package bootstrap

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"context"
	"time"
)

func NewServiceDiscovery(env *ExporterEnv) discovery.ServiceDiscovery {
	config := env.DiscoveryConfig
	return discovery.NewServiceDiscovery(context.Background(), config.Endpoints, config.Timeout*time.Second,
		config.Username, config.Password)
}
