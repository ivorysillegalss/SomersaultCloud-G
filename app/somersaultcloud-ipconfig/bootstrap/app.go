package bootstrap

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"SomersaultCloud/app/somersaultcloud-ipconfig/dispatcher"
	"SomersaultCloud/app/somersaultcloud-ipconfig/source"
)

type IpConfigApplication struct {
	Env         *IpConfigEnv
	Api         *Api
	Discovery   discovery.ServiceDiscovery
	DataHandler *source.DataHandler
	Dispatcher  *dispatcher.Dispatcher
}
