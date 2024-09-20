package task

import (
	"SomersaultCloud/api/middleware/taskchain"
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/sys"
	"SomersaultCloud/constant/task"
	"SomersaultCloud/domain"
	"SomersaultCloud/handler"
	"SomersaultCloud/internal/checkutil"
	"context"
	"github.com/thoas/go-funk"
	"strconv"
	"sync"
	"time"
)

// ChatAskTask rpc调用责任链任务实现
type ChatAskTask struct {
	chatRepository domain.ChatRepository
	botRepository  domain.BotRepository
	chatEvent      domain.ChatEvent
	env            *bootstrap.Env
	channels       *bootstrap.Channels
	poolFactory    *bootstrap.PoolsFactory
}

func NewAskChatTask(b domain.BotRepository, c domain.ChatRepository, e *bootstrap.Env, ch *bootstrap.Channels, p *bootstrap.PoolsFactory, ce domain.ChatEvent) AskTask {
	return &ChatAskTask{chatRepository: c, botRepository: b, env: e, channels: ch, poolFactory: p, chatEvent: ce}
}

func (c *ChatAskTask) InitContextData(args ...any) *taskchain.TaskContext {
	userId := args[0].(int)
	botId := args[1].(int)
	chatId := args[2].(int)
	message := args[3].(string)
	return &taskchain.TaskContext{
		BusinessType:    args[4].(string),
		BusinessCode:    args[5].(int),
		TaskContextData: &domain.AskContextData{UserId: userId, BotId: botId, ChatId: chatId, Message: message, ExecutorId: args[6].(int)},
	}
}

func (c *ChatAskTask) PreCheckDataTask(tc *taskchain.TaskContext) {

	data := tc.TaskContextData.(*domain.AskContextData)
	//TODO 运行前redis加缓存
	chatIdCheck := checkutil.IsLegalID(data.ChatId, common.FalseInt, c.chatRepository.CacheGetNewestChatId(context.Background()))
	botIdCheck := checkutil.IsLegalID(data.BotId, common.FalseInt, c.botRepository.CacheGetMaxBotId(context.Background()))
	message := data.Message
	msgCheck := funk.NotEmpty(message)
	if !(msgCheck && chatIdCheck && botIdCheck) {
		tc.InterruptExecute(task.InvalidDataFormatMessage)
		return
	}
}

// GetHistoryTask 2情况 判断是否存在缓存 hit拿缓存 miss则db
func (c *ChatAskTask) GetHistoryTask(tc *taskchain.TaskContext) {

	data := tc.TaskContextData.(*domain.AskContextData)

	var history *[]*domain.Record
	// 1. 缓存找
	history, NotHaveCache, err := c.chatRepository.CacheGetHistory(context.Background(), data.ChatId)
	if err != nil {
		tc.InterruptExecute(task.HistoryRetrievalFailed)
		return
	}

	// 2. 缓存miss db找
	//TODO 目前查DB后需要截取历史记录，实现数据流式更新后可取消
	if NotHaveCache {
		var title string
		history, title, err = c.chatRepository.DbGetHistory(context.Background(), data.ChatId)
		if err != nil {
			tc.InterruptExecute(task.HistoryRetrievalFailed)
			return
		}

		if funk.IsEmpty(history) {
			history = new([]*domain.Record)

		} else {
			// 截取数据
			if len(*history) >= cache.HistoryDefaultWeight {
				*history = (*history)[:cache.HistoryDefaultWeight]
			}

			// 2.1 回写缓存 (把从DB拿到的回写缓存 维护热点数据)
			//TODO 目前架构下，chat一次请求回写两次缓存，可优化，取消此次回写
			go c.chatRepository.CacheLuaLruResetHistory(context.Background(),
				cache.ChatHistoryScore+common.Infix+strconv.Itoa(data.UserId), data.History, data.ChatId, title)
		}
	}

	data.History = history
}

func (c *ChatAskTask) GetBotTask(tc *taskchain.TaskContext) {

	data := tc.TaskContextData.(*domain.AskContextData)

	botConfig := c.botRepository.CacheGetBotConfig(context.Background(), data.BotId)
	if funk.IsEmpty(botConfig) {
		tc.InterruptExecute(task.BotRetrievalFailed)
		return
	}

	data.SysPrompt = botConfig.InitPrompt
	data.Model = botConfig.Model
}

func (c *ChatAskTask) AdjustmentTask(tc *taskchain.TaskContext) {
	//TODO implement me
	panic("implement me")
}

func (c *ChatAskTask) AssembleReqTask(tc *taskchain.TaskContext) {

	data := tc.TaskContextData.(*domain.AskContextData)

	//TODO id在此处没什么作用 主要为了之后多实现 策略化 先随便传一个
	executor := handler.NewLanguageModelExecutor(c.env, c.channels, data.ExecutorId)

	data.Executor = executor
	data.HistoryMessage = executor.AssemblePrompt(data)
	//无需判空 因为第一次聊情况下就是没有历史记录的

	request := executor.EncodeReq(data)
	if funk.IsEmpty(request) {
		tc.InterruptExecute(task.ReqDataMarshalFailed)
		return
	}
	client := executor.ConfigureProxy(data)
	data.Conn = *domain.NewConnection(client, request)
}

func (c *ChatAskTask) CallApiTask(tc *taskchain.TaskContext) {

	data := tc.TaskContextData.(*domain.AskContextData)

	var wg sync.WaitGroup
	//包装提交的任务
	t := func() {
		defer wg.Done()
		data.Executor.Execute(data)
	}
	config := c.poolFactory.Pools[sys.ExecuteRpcGoRoutinePool]
	//使用Invoke方法 所返回的是线程池本身在操作中遇到的err
	wg.Add(1)
	err := config.Submit(t)

	//err := config.Invoke(t)
	if err != nil {
		tc.InterruptExecute(task.ReqUploadError)
		return
	}

	//TODO 消息队列
}

func (c *ChatAskTask) ParseRespTask(tc *taskchain.TaskContext) {

	data := tc.TaskContextData.(*domain.AskContextData)

	var generation *domain.GenerationResponse
	//没查到的话有可能是没处理完 等个300ms再查
	//循环查询最多10次 超过则宣布失败
	for i := 0; i < sys.GenerateQueryRetryLimit; i++ {
		if funk.IsEmpty(generation) {
			//轮询等待
			g := c.chatRepository.MemoryGetGeneration(context.Background(), data.ChatId)
			generation = g
			time.Sleep(600 * time.Millisecond)
		} else {
			break
		}
	}

	//TODO 这里其实是有个bug的 如果超过10次收不到 大部分情况下是rpc失败的 但是也有小部分情况调用成功
	//	但是未存储 这种造成了一种情况 因为下方删除是确定缓存存在了才删的 而超时的情况则默认了缓存不存在
	//	假如缓存在超时之后 来到了    在下一次请求的时候 就会读到同一个chat上次rpc时因为超时 而未渲染出来的generation
	//	暂时还没想对策
	if funk.IsEmpty(generation) {
		tc.InterruptExecute(task.ReqCatchError)
		return
	} else {
		//旁路缓存思想 如果缓存获取成功删掉他 防止短时间内 生成内容未覆盖就读到上次的generation
		c.chatRepository.MemoryDelGeneration(context.Background(), data.ChatId)
	}

	//直到此处成功获取到resp对象 此处关流
	data.Resp = *generation
	defer generation.Resp.Body.Close()

	resp, _ := data.Executor.ParseResp(data)

	if funk.IsEmpty(resp) {
		tc.InterruptExecute(task.RespParedError)
	}

	data.ParsedResponse = resp

	if funk.Equal(tc.BusinessCode, task.ExecuteChatAskCode) {
		//回写缓存&DB
		c.chatEvent.PublishSaveCacheHistory(data)
		c.chatEvent.PublishSaveDbHistory(data)
	} else if funk.Equal(tc.BusinessCode, task.ExecuteTitleAskCode) {
		c.chatEvent.PublishDbSaveTitle(data)
		//感觉没必要都上mq 写缓存就直接go好了
		go c.chatRepository.CacheUpdateTitle(context.Background(), data.ChatId, data.ParsedResponse.GetGenerateText())
	}
}
