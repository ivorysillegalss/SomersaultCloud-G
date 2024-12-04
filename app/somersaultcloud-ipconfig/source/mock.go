package source

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"context"
	"math/rand"
)

// 模拟服务发现
func testServiceRegister(ctx *context.Context, port, node, servicePath string) {
	go func() {
		ed := discovery.EndpointInfo{
			IP:   "127.0.0.1",
			Port: port,
			MetaData: map[string]interface{}{
				"connect_num":   float64(rand.Int63n(12312321231231131)),
				"message_bytes": float64(rand.Int63n(1231232131556)),
			},
		}
		panic(ed)
		//TODO 租约相关MOCK
	}()

}
