package source

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"SomersaultCloud/app/somersaultcloud-common/log"
	"context"
	"fmt"
	"math/rand"
	"time"
)

// 模拟服务发现
func testServiceRegister(ctx context.Context, port, node, servicePath string, cli discovery.ServiceDiscovery) {
	go func() {
		ed := discovery.EndpointInfo{
			IP:   "127.0.0.1",
			Port: port,
			MetaData: map[string]interface{}{
				"connect_num":   float64(rand.Int63n(12312321231231131)),
				"message_bytes": float64(rand.Int63n(1231232131556)),
			},
		}
		client := cli.GetClient().Cli
		sr, err := discovery.NewServiceRegister(client, fmt.Sprintf("%s/%s", servicePath, node), ctx, &ed, time.Now().Unix())
		if err != nil {
			log.GetTextLogger().Error("create service register error")
		}

		//go一个协程监听每一个服务的信息
		go sr.ListenLeaseRespChan()

		//模拟服务的信息改变 不断修改服务发现配置中心中的信息
		for {
			ed = discovery.EndpointInfo{
				IP:   "127.0.0.1",
				Port: port,
				MetaData: map[string]interface{}{
					"connect_num":   float64(rand.Int63n(12312321231231131)),
					"message_bytes": float64(rand.Int63n(1231232131556)),
				},
			}
			_ = sr.UpdateValue(&ed)
			time.Sleep(1 * time.Second)
		}
	}()

}
