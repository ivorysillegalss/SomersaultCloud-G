package grpc

import (
	"SomersaultCloud/app/somersaultcloud-exporter/bootstrap"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

func NewGrpcConn(env *bootstrap.ExporterEnv) *bootstrap.GrpcConn {
	gEnv := env.Grpc
	serverAddr := gEnv.Monitor.Server + ":" + strconv.Itoa(gEnv.Monitor.Port)

	//TODO withSecure
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	return &bootstrap.GrpcConn{Conn: conn}
}
