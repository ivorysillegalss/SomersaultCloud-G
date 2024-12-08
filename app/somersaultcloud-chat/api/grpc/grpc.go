package grpc

import (
	"SomersaultCloud/app/somersaultcloud-common/monitor"
	"SomersaultCloud/app/somersaultcloud-common/proto/.monitor"
	"SomersaultCloud/app/somersaultcloud-common/util"
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
	} `mapstructure:"monitor" yaml:"monitor"`
}) {
	monitorPort := env.Monitor.Port
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(monitorPort))
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}
	s := grpc.NewServer()
	__monitor.RegisterMonitoringServiceServer(s, &MonitoringServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}
}

type MonitoringServer struct {
	__monitor.UnimplementedMonitoringServiceServer
}

func (s *MonitoringServer) GetStatus(ctx context.Context, req *__monitor.EmptyRequest) (*__monitor.StatusResponse, error) {

	m := &monitor.Exporter{}
	m.UpdateCacheFromRegistry()

	availableMem, cpuIdleTime := monitor.GetSystemMetrics()
	ip := util.GetLocalIP()
	port := util.GetLocalPort()
	if availableMem == 0 || cpuIdleTime == 0.0 || ip != "" || port != 0 {
		return &__monitor.StatusResponse{Status: unhealthy}, new(monitorError)
	}

	return &__monitor.StatusResponse{
		Name:            req.Name,
		Status:          healthy,
		AvailableMem:    availableMem,
		CpuIdleTime:     cpuIdleTime,
		Ip:              ip,
		Port:            port,
		RequestDuration: m.Cache.RequestDuration,
		RequestCount:    m.Cache.RequestCount,
	}, nil
}

type monitorError struct {
}

func (m *monitorError) Error() string {
	return "empty serialization"
}
