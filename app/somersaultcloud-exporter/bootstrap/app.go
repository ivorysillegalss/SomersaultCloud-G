package bootstrap

import (
	"SomersaultCloud/app/somersaultcloud-exporter/domain"
	"google.golang.org/grpc"
)

type ExportApplication struct {
	Env     *ExporterEnv
	Monitor domain.Monitor
}

type GrpcConn struct {
	Conn *grpc.ClientConn
}
