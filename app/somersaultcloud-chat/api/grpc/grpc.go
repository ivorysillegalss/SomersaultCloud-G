package grpc

import (
	pb "SomersaultCloud/app/somersaultcloud-chat/proto/.monitor"
	"SomersaultCloud/app/somersaultcloud-common/monitor"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

const (
	healthy   = "healthy"
	unhealthy = "unhealthy"
)

func Setup(env struct {
	Monitor struct {
		Port int `mapstructure:"port" yaml:"port"`
	}
}) {
	monitorPort := env.Monitor.Port
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(monitorPort))
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}
	s := grpc.NewServer()
	pb.RegisterMonitoringServiceServer(s, &MonitoringServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}
}

type MonitoringServer struct {
	pb.UnimplementedMonitoringServiceServer
}

func (s *MonitoringServer) GetStatus(ctx context.Context, req *pb.EmptyRequest) (*pb.StatusResponse, error) {
	availableMem, cpuIdleTime := monitor.GetSystemMetrics()
	if availableMem == 0 || cpuIdleTime == 0.0 {
		return &pb.StatusResponse{Status: unhealthy}, new(monitorError)
	}

	return &pb.StatusResponse{
		Status:       healthy,
		AvailableMem: availableMem,
		CpuIdleTime:  cpuIdleTime,
	}, nil
}

type monitorError struct {
}

func (m *monitorError) Error() string {
	return "empty serialization"
}
