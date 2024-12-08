package monitor

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"SomersaultCloud/app/somersaultcloud-common/log"
	pb "SomersaultCloud/app/somersaultcloud-common/proto/.monitor"
	"SomersaultCloud/app/somersaultcloud-exporter/bootstrap"
	"SomersaultCloud/app/somersaultcloud-exporter/domain"
	"context"
	"fmt"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc"
	"math/rand"
	"time"
)

type monitor struct {
	conn      *bootstrap.GrpcConn
	env       *bootstrap.ExporterEnv
	discovery discovery.ServiceDiscovery
}

func NewMonitor(conn *bootstrap.GrpcConn, env *bootstrap.ExporterEnv, dis discovery.ServiceDiscovery) domain.Monitor {
	return &monitor{conn: conn, env: env, discovery: dis}
}

var statusMap map[string]domain.MonitorStatus

const (
	//多服务时服务的数量 每一个服务对应一个etcd的客户端 独赢一个chan的value
	applicationNums = 2
)

func init() {
	statusMap = make(map[string]domain.MonitorStatus, applicationNums)
}

func (m *monitor) ServiceRegister() {
	// 从配置文件读取所有服务的初始化信息
	intiAddress := m.env.BusinessConfig.Address
	for _, address := range intiAddress {
		randRaw := &pb.StatusResponse{
			Status:       "healthy",
			AvailableMem: uint64(rand.Int63n(12312321231231131)),
			CpuIdleTime:  float64(rand.Int63n(1231232131556)),
		}
		ed := discovery.EndpointInfo{
			IP:       address.IP,
			Port:     string(address.Port),
			MetaData: convertMetaData(randRaw),
		}

		client := m.discovery.GetClient().Cli
		sr, err := discovery.NewServiceRegister(client, fmt.Sprintf("%s/%s",
			m.env.DiscoveryConfig.ServicePath, address.Name), context.Background(), &ed, time.Now().Unix())

		go sr.ListenLeaseRespChan()
		statusMap[address.Name] = domain.MonitorStatus{ServiceRegister: sr, EndpointInfo: &ed, Time: time.Now()}

		if err != nil {
			log.GetTextLogger().Error("create service register error for node %s", address.Name)
		}
	}

}

func (m *monitor) HandleMonit() {
	//k is the name of the node(service)
	for k, _ := range statusMap {
		go func(k string) {
			for {
				handle(k, m.conn.Conn)
				time.Sleep(time.Second)
				log.GetTextLogger().Info("Success updating value with node: %s", k)
			}
		}(k)
	}
}

func handle(serviceName string, grpcConn *grpc.ClientConn) {
	client := pb.NewMonitoringServiceClient(grpcConn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := client.GetStatus(ctx, &pb.EmptyRequest{Name: serviceName})
	if err != nil {
		log.GetTextLogger().Error("Error calling GetStatus: %v", err)
		return
	}
	log.GetTextLogger().Info("Success calling GetStatus %v", err)

	info := discovery.EndpointInfo{
		IP:       resp.GetIp(),
		Port:     string(resp.GetPort()),
		MetaData: convertMetaData(resp),
	}

	status := statusMap[serviceName]
	if funk.IsEmpty(status) {
		log.GetTextLogger().Error("cannot get service status for service %s", serviceName)
	}

	status.EndpointInfo = &info
	status.Time = time.Now()
	statusMap[serviceName] = status

	_ = status.ServiceRegister.UpdateValue(&info)
}

func convertMetaData(raw *pb.StatusResponse) map[string]any {
	m := make(map[string]any)
	m["status"] = raw.GetStatus()
	m["available_mem"] = raw.GetAvailableMem()
	m["cpu_idle_time"] = raw.GetCpuIdleTime()
	return m
}
