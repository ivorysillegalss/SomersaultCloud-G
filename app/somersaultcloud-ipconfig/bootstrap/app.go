package bootstrap

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"SomersaultCloud/app/somersaultcloud-ipconfig/domain"
)

type IpConfigApplication struct {
	Env         *IpConfigEnv
	Api         *Api
	Discovery   discovery.ServiceDiscovery
	DataHandler domain.DataHandler
	Dispatcher  domain.Dispatcher
}
