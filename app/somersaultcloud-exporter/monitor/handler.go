package monitor

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"SomersaultCloud/app/somersaultcloud-common/log"
	pb "SomersaultCloud/app/somersaultcloud-common/proto/.monitor"
	"SomersaultCloud/app/somersaultcloud-exporter/bootstrap"
	"SomersaultCloud/app/somersaultcloud-exporter/domain"
	"bytes"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc"
	"math/rand"
	"net/http"
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

var (
	statusMap map[string]domain.MonitorStatus

	gaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "chat_module_load_gauge",
			Help: "A custom metric for chat_module getting instances load",
		},
		[]string{"load"}, // 标签，可以根据需要添加更多
	)
)

const (
	//多服务时服务的数量 每一个服务对应一个etcd的客户端 独赢一个chan的value
	applicationNums = 2

	healthy   = 1
	unhealthy = 0
)

func init() {
	statusMap = make(map[string]domain.MonitorStatus, applicationNums)

	prometheus.MustRegister(gaugeVec)
}

func (m *monitor) ServiceRegister() {
	// 从配置文件读取所有服务的初始化信息
	//初始化的时候 随机状态信息 因为此时所有的服务器状态基本上都是一样的
	initAddress := m.env.BusinessConfig.Address
	for _, address := range initAddress {
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
				//time.Sleep(time.Second)
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

	//TODO 这里可能可以空间复用优化
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

	//更新普罗米修斯 prometheus 中的信息
	updateLabel(resp)

	//更新etcd中的配置信息
	_ = status.ServiceRegister.UpdateValue(&info)
}

// convertMetaData 此方法将pb返回值转换更新 拓展式数据设计map
func convertMetaData(raw *pb.StatusResponse) map[string]any {
	m := make(map[string]any)
	m["status"] = raw.GetStatus()
	m["available_mem"] = raw.GetAvailableMem()
	m["cpu_idle_time"] = raw.GetCpuIdleTime()
	m["request_count"] = raw.GetRequestCount()
	m["request_duration"] = raw.GetRequestDuration()
	return m
}

func updateLabel(raw *pb.StatusResponse) {
	gaugeVec.WithLabelValues(raw.GetName(), "available_mem").Set(float64(raw.GetAvailableMem()))
	gaugeVec.WithLabelValues(raw.GetName(), "request_duration").Set(raw.GetRequestDuration())
	gaugeVec.WithLabelValues(raw.GetName(), "request_count").Set(raw.GetRequestCount())
	gaugeVec.WithLabelValues(raw.GetName(), "cpu_idle_time").Set(raw.GetCpuIdleTime())
	gaugeVec.WithLabelValues(raw.GetName(), "status").Set(healthy)
}

func (m *monitor) ExposeMonitorInterface(ctx context.Context, c *app.RequestContext) {
	log.GetTextLogger().Info("Prometheus pulling.....")
	c.Response.Header.Set("Content-Type", "text/plain; version=0.0.4")
	metrics, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		log.GetTextLogger().Error("Could not gather metrics: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// 使用 expfmt 将指标编码为文本格式
	var buf bytes.Buffer
	encoder := expfmt.NewEncoder(&buf, expfmt.FmtText)
	for _, mf := range metrics {
		if err := encoder.Encode(mf); err != nil {
			hlog.Errorf("Could not encode metric family: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	// 写入响应
	c.Write(buf.Bytes())
}
