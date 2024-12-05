package source

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"SomersaultCloud/app/somersaultcloud-common/log"
	"SomersaultCloud/app/somersaultcloud-ipconfig/bootstrap"
	"context"
)

type DataHandler struct {
	IpConfigEnv      *bootstrap.IpConfigEnv
	ServiceDiscovery discovery.ServiceDiscovery
}

// Handle 服务发现处理 （当各服务有新的改动来到时，以etcd监听机制实现热更新）
func (d *DataHandler) Handle() {
	eventChan = make(chan *Event)
	go handle(d.ServiceDiscovery, d.IpConfigEnv.ServicePath)
	//测试环境下mock出对应的测试诗句进行测试
	if d.IpConfigEnv.AppEnv == "debug" {
		ctx := context.Background()
		testServiceRegister(ctx, "7896", "node1", d.IpConfigEnv.ServicePath, d.ServiceDiscovery)
		testServiceRegister(ctx, "7897", "node2", d.IpConfigEnv.ServicePath, d.ServiceDiscovery)
		testServiceRegister(ctx, "7898", "node3", d.IpConfigEnv.ServicePath, d.ServiceDiscovery)
	}
}

func handle(dis discovery.ServiceDiscovery, servicePath string) {
	setFunc := func(key, value string) {
		if ed, err := discovery.UnMarshal([]byte(value)); err == nil {
			if event := NewEvent(ed); ed != nil {
				event.Type = AddNodeEvent
				eventChan <- event
			}
		}
	}

	delFunc := func(key, value string) {
		if info, err := discovery.UnMarshal([]byte(value)); err == nil {
			if event := NewEvent(info); info != nil {
				event.Type = DelNodeEvent
				eventChan <- event
			}
		}
	}

	err := dis.WatchService(servicePath, setFunc, delFunc)
	if err != nil {
		log.GetTextLogger().Fatal(err.Error())
	}
}
