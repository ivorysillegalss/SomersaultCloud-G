package cron

import (
	"SomersaultCloud/app/somersaultcloud-chat/bootstrap"
	"SomersaultCloud/app/somersaultcloud-chat/constant/common"
	"SomersaultCloud/app/somersaultcloud-chat/constant/dao"
	"SomersaultCloud/app/somersaultcloud-chat/domain"
	"SomersaultCloud/app/somersaultcloud-chat/handler"
	"SomersaultCloud/app/somersaultcloud-chat/handler/stream"
	"SomersaultCloud/app/somersaultcloud-chat/internal/kvutil"
	log2 "SomersaultCloud/app/somersaultcloud-common/log"
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
			log2.GetTextLogger().Info("RPCRES CHANNEL ---------ASYNC POLLING")
		} else {
			i--
		}

		//TODO，不使用单select语句，架构上优化
		select {
		case task := <-g.channels.RpcRes:
			log2.GetTextLogger().Info("RPCRES CHANNEL < -- SUCCESSFULLY RECEIVED RESPONSE")
			g.generationRepository.InMemoryPollHistory(context.Background(), task)
		case <-g.channels.Stop:
			log2.GetJsonLogger().WithFields("res channel stop", nil).Info("Stopping async poller")
			return
		case streamTask := <-g.channels.StreamRpcRes:
			//log.GetTextLogger().Info("STREAM RPCRES CHANNEL < -- SUCCESSFULLY RECEIVED STREAM RESPONSE")
			consumeAndParse(streamTask, g.env, g.channels, g.generateEvent, g.generationRepository)
		default:
			time.Sleep(500 * time.Millisecond) // 控制轮询频率
		}
	}
}

// consumeAndParse TODO MOVE
func consumeAndParse(streamTask *domain.GenerationResponse, env *bootstrap.Env, channels *bootstrap.Channels, event domain.GenerateEvent, repository domain.GenerationRepository) {
	executor := handler.NewLanguageModelExecutor(env, channels, streamTask.ExecutorId)
	parsedResp, _ := executor.ParseResp(&domain.AskContextData{Resp: *streamTask, Stream: true, ExecutorId: streamTask.ExecutorId})
	index := chatcmplCache.IndexIncIfExist(parsedResp.GetChatcmplId())
	parsedResp.SetIndex(index)
	//log.GetTextLogger().Info("start parsing value for chatcmplId: " + parsedResp.GetChatcmplId())

	newSequencer := stream.NewSequencer()
	newSequencer.Setup(parsedResp)
}
