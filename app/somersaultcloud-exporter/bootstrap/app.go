package bootstrap

import (
	"SomersaultCloud/app/somersaultcloud-exporter/domain"
	"github.com/cloudwego/hertz/pkg/app/server"
	"google.golang.org/grpc"
)

type ExportApplication struct {
	Env     *ExporterEnv
	Monitor domain.Monitor
}

type GrpcConn struct {
	Conn *grpc.ClientConn
}

// Setup pull和push入口
func (e *ExportApplication) Setup() {
	//push数据进入etcd
	monit(e)

	//Prometheus pull数据入口
	s := server.Default(server.WithHostPorts(e.Env.ExporterConfig.ServerAddress))
	s.GET("/metrics", e.Monitor.ExposeMonitorInterface)
	s.Spin()
}

// monit 开启不断grpc抓取消息更新到ipconfig中
func monit(e *ExportApplication) {
	e.Monitor.ServiceRegister()
	go func() {
		e.Monitor.HandleMonit()
	}()
}
