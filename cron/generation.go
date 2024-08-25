package cron

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/log"
	"context"
	"time"
)

type generationCron struct {
	generationRepository domain.GenerationRepository
	channels             *bootstrap.Channels
}

// AsyncPollerGeneration 异步轮询chat的generation的函数 仅负责映射到map中 程序启动则执行
// TODO EPOLL&线程池优化，进一步抽象至仓库中
func (g generationCron) AsyncPollerGeneration() {
	for {
		log.GetTextLogger().Info("RPCRES CHANNEL ---------ASYNC POLLING")
		select {
		case task := <-g.channels.RpcRes:
			log.GetTextLogger().Info("RPCRES CHANNEL < -- SUCCESSFULLY RECEIVED RESPONSE")
			g.generationRepository.InMemoryPollHistory(context.Background(), task)
		case <-g.channels.Stop:
			log.GetJsonLogger().WithFields("res channel stop", nil).Info("Stopping async poller")
			return
		default:
			time.Sleep(500 * time.Millisecond) // 控制轮询频率
		}
	}
}

func NewGenerationCron(g domain.GenerationRepository, c *bootstrap.Channels) domain.GenerationCron {
	return generationCron{generationRepository: g, channels: c}
}
