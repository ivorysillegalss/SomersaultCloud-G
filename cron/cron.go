package cron

import (
	"SomersaultCloud/domain"
	"SomersaultCloud/repository"
	"context"
	"fmt"
	"time"
)

//TODO 定时任务模块分包优化

// AsyncPoller 异步轮询chat的generation的函数 仅负责映射到map中 程序启动则执行
// TODO 线程池优化，进一步抽象至仓库中
func (a *asyncService) AsyncPoller(resps <-chan *domain.GenerationResponse, stop <-chan bool) {
	for {
		select {
		case task := <-resps:
			a.generationRepository.CacheLuaPollHistory(context.Background(), *task)
		case <-stop:
			fmt.Println("Stopping async poller")
			return
		default:
			time.Sleep(500 * time.Millisecond) // 控制轮询频率
		}
	}
}

type AsyncService interface {
	// AsyncPoller 只读channel
	AsyncPoller(resps <-chan *domain.GenerationResponse, stop <-chan bool)
}

type asyncService struct {
	generationRepository domain.GenerationRepository
}

func newAsyncService() AsyncService {
	return &asyncService{generationRepository: repository.NewGenerationRepository()}
}

func Setup() {
	service := newAsyncService()
	go service.AsyncPoller(c.RpcRes, c.Stop) //只读
}