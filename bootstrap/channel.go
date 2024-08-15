package bootstrap

import (
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/sys"
	"SomersaultCloud/infrastructure/channel"
	"SomersaultCloud/infrastructure/redis"
	"context"
	"fmt"
	"strconv"
	"time"
)

func NewChannel(r redis.Client) *Channels {
	//TODO channel类型的增多可以在channel结构体中增加 并在此处初始化
	//只读
	go asyncPoller(r, make(<-chan *channel.GenerationResponse), make(<-chan bool))
	return &Channels{RpcRes: make(chan *channel.GenerationResponse, sys.GenerationResponseChannelBuffer)}
}

//	TODO 迁移到定时任务板块中
//
// 异步轮询chat的generation的函数
// 仅负责映射到map中
func asyncPoller(r redis.Client, resps <-chan *channel.GenerationResponse, stop <-chan bool) {
	for {
		select {
		case task := <-resps:
			err := r.SetStructExpire(context.Background(), cache.ChatGeneration+common.Infix+strconv.Itoa(task.ChatId), task, cache.ChatGenerationDDL*time.Minute)
			if err != nil {
				//TODO 打日志
			}
		case <-stop:
			fmt.Println("Stopping async poller")
			return
		default:
			time.Sleep(500 * time.Millisecond) // 控制轮询频率
		}
	}
}
