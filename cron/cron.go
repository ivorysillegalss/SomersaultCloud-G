package cron

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/log"
	"context"
	"time"
)

// AsyncPoller 异步轮询chat的generation的函数 仅负责映射到map中 程序启动则执行
// TODO EPOLL&线程池优化，进一步抽象至仓库中
func (a *asyncService) AsyncPoller() {
	for {
		log.GetTextLogger().Info("RPCRES CHANNEL ---------ASYNC POLLING")
		select {
		case task := <-a.channels.RpcRes:
			log.GetTextLogger().Info("RPCRES CHANNEL < -- SUCCESSFULLY RECEIVED RESPONSE")
			a.generationRepository.InMemoryPollHistory(context.Background(), task)
		case <-a.channels.Stop:
			log.GetJsonLogger().WithFields("res channel stop", nil).Info("Stopping async poller")
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
