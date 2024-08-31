package cron

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/dao"
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/log"
	"context"
	"github.com/thoas/go-funk"
	"time"
)

type generationCron struct {
	generationRepository domain.GenerationRepository
	channels             *bootstrap.Channels
}

// AsyncPollerGeneration 异步轮询chat的generation的函数 仅负责映射到map中 程序启动则执行
// TODO EPOLL&线程池优化，进一步抽象至仓库中
func (g generationCron) AsyncPollerGeneration() {
	i := dao.AsyncPollingFrequency
	for {
		if funk.Equal(i, common.ZeroInt) {
			i = dao.AsyncPollingFrequency
			log.GetTextLogger().Info("RPCRES CHANNEL ---------ASYNC POLLING")
		} else {
			i--
		}

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
