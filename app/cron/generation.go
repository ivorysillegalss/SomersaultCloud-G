package cron

import (
	"SomersaultCloud/app/bootstrap"
	"SomersaultCloud/app/constant/common"
	"SomersaultCloud/app/constant/dao"
	"SomersaultCloud/app/domain"
	"SomersaultCloud/app/handler"
	"SomersaultCloud/app/handler/stream"
	"SomersaultCloud/app/infrastructure/log"
	"SomersaultCloud/app/internal/kvutil"
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
			//log.GetTextLogger().Info("STREAM RPCRES CHANNEL < -- SUCCESSFULLY RECEIVED STREAM RESPONSE")
			//TODO planA本身是go这个转码的函数，然后再下发的时候，通过排序器对失序做处理再下发的，但是太多bug了，先用同步上线
			// 在同步的情况下，由于openai是经过http1.1or2的流式传输，所以是可以保证消息的有序的。唯一的可失序的地方就是mq 消息丢失
			// 待优化
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
