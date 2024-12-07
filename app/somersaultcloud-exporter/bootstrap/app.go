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
	monit(e)

	s := server.Default(server.WithHostPorts(":7890"))
	//TODO 补充，Prometheus，消息抓取接口
	s.Spin()
}

// monit 开启不断grpc抓取消息更新到ipconfig中
func monit(e *ExportApplication) {
	e.Monitor.ServiceRegister()
	go func() {
		e.Monitor.HandleMonit()
	}()
}
