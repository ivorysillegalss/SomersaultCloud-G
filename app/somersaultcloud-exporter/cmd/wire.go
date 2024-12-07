//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"SomersaultCloud/app/somersaultcloud-exporter/api/grpc"
	"SomersaultCloud/app/somersaultcloud-exporter/bootstrap"
	"SomersaultCloud/app/somersaultcloud-exporter/monitor"
	"github.com/google/wire"
)

var appSet = wire.NewSet(
	bootstrap.NewEnv,
	bootstrap.NewServiceDiscovery,
	grpc.NewGrpcConn,
	monitor.NewMonitor,

	wire.Struct(new(bootstrap.ExportApplication), "*"),
)

// InitializeApp init application.
func InitializeApp() (*bootstrap.ExportApplication, error) {
	wire.Build(appSet)
	return &bootstrap.ExportApplication{}, nil
}
