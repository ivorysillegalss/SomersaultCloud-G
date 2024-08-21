package cron

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/domain"
	"context"
	"fmt"
	"time"
)

//TODO 定时任务模块分包优化

// AsyncPoller 异步轮询chat的generation的函数 仅负责映射到map中 程序启动则执行
// TODO 线程池优化，进一步抽象至仓库中
func (a *asyncService) AsyncPoller() {
	for {
		fmt.Println("-----------------------------------------")
		select {
		case task := <-a.channels.RpcRes:
			fmt.Println("get -----------------------------------")
			//a.generationRepository.CacheLuaPollHistory(context.Background(), *task)
			a.generationRepository.InMemoryPollHistory(context.Background(), task)
		case <-a.channels.Stop:
			fmt.Println("Stopping async poller")
			//TODO 打日志
			return
		default:
			time.Sleep(500 * time.Millisecond) // 控制轮询频率
		}
	}
}

type AsyncService interface {
	// AsyncPoller 只读channel
	AsyncPoller()
}

type asyncService struct {
	generationRepository domain.GenerationRepository
	channels             *bootstrap.Channels
}

func NewAsyncService(generationRepository domain.GenerationRepository, channels *bootstrap.Channels) AsyncService {
	return &asyncService{generationRepository: generationRepository, channels: channels}
}
