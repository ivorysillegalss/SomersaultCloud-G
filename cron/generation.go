package cron

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/dao"
	"SomersaultCloud/domain"
	"SomersaultCloud/handler"
	"SomersaultCloud/infrastructure/log"
	"SomersaultCloud/internal/kvutil"
	"context"
	"github.com/thoas/go-funk"
	"time"
)

var chatcmplCache kvutil.KVStore

func init() {
	chatcmplCache = *kvutil.NewKVStore(time.Second, 20)
}

type generationCron struct {
	generationRepository domain.GenerationRepository
	channels             *bootstrap.Channels
	env                  *bootstrap.Env
	generateEvent        domain.GenerateEvent
}

// TODO 依赖注入
func NewGenerationCron(g domain.GenerationRepository, c *bootstrap.Channels, e *bootstrap.Env, ge domain.GenerateEvent) domain.GenerationCron {
	return generationCron{generationRepository: g, channels: c, env: e, generateEvent: ge}
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

		//TODO，不使用单select语句，架构上优化
		select {
		case task := <-g.channels.RpcRes:
			log.GetTextLogger().Info("RPCRES CHANNEL < -- SUCCESSFULLY RECEIVED RESPONSE")
			g.generationRepository.InMemoryPollHistory(context.Background(), task)
		case <-g.channels.Stop:
			log.GetJsonLogger().WithFields("res channel stop", nil).Info("Stopping async poller")
			return
		case streamTask := <-g.channels.StreamRpcRes:
			log.GetTextLogger().Info("STREAM RPCRES CHANNEL < -- SUCCESSFULLY RECEIVED STREAM RESPONSE")
			//TODO 测试，目前go一个线程parse是我想到的最好方法
			go consumeAndParse(streamTask, g.env, g.channels, g.generateEvent)
		default:
			time.Sleep(500 * time.Millisecond) // 控制轮询频率
		}
	}
}

// consumeAndParse TODO MOVE
func consumeAndParse(streamTask *domain.GenerationResponse, env *bootstrap.Env, channels *bootstrap.Channels, event domain.GenerateEvent) {
	executor := handler.NewLanguageModelExecutor(env, channels, streamTask.ExecutorId)
	parsedResp, _ := executor.ParseResp(&domain.AskContextData{Resp: *streamTask, Stream: true, ExecutorId: streamTask.ExecutorId})
	index := chatcmplCache.IndexIncIfExist(parsedResp.GetChatcmplId())
	parsedResp.SetIndex(index)
	event.PublishGeneration(parsedResp)
}
