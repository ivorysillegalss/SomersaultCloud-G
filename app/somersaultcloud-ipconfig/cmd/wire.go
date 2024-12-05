//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"SomersaultCloud/app/somersaultcloud-ipconfig/bootstrap"
	"SomersaultCloud/app/somersaultcloud-ipconfig/dispatcher"
	"SomersaultCloud/app/somersaultcloud-ipconfig/source"
	"github.com/google/wire"
)

var appSet = wire.NewSet(
	bootstrap.NewEnv,
	bootstrap.NewServiceDiscovery,
	bootstrap.NewApi,

	dispatcher.NewDispatcher,
	source.NewDataHandler,

	wire.Struct(new(bootstrap.IpConfigApplication), "*"),
)

// InitializeApp init application.
func InitializeApp() (*bootstrap.IpConfigApplication, error) {
	wire.Build(appSet)
	return &bootstrap.IpConfigApplication{}, nil
}
