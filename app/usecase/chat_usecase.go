package usecase

import (
	"SomersaultCloud/app/api/middleware/taskchain"
	"SomersaultCloud/app/bootstrap"
	"SomersaultCloud/app/constant/cache"
	"SomersaultCloud/app/constant/common"
	task2 "SomersaultCloud/app/constant/task"
	"SomersaultCloud/app/domain"
	"SomersaultCloud/app/handler/stream"
	"SomersaultCloud/app/infrastructure/log"
	"SomersaultCloud/app/internal/tokenutil"
	"SomersaultCloud/app/task"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/thoas/go-funk"
	"net/http"
	"strconv"
	"time"
)

//go:embed lua/increment.lua
var incrementLuaScript string

// StreamStorageMap 以用户ID为细粒度 存储成功下发的流信息
var streamStorageMap = make(map[int]chan string)

type chatUseCase struct {
	env                  *bootstrap.Env
	chatRepository       domain.ChatRepository
	botRepository        domain.BotRepository
	generationRepository domain.GenerationRepository
	chatTask             task.AskTask
	tokenUtil            *tokenutil.TokenUtil
	chatEvent            domain.StorageEvent
	titleTask            task.TitleTask
	convertTask          task.ConvertTask
}

func NewChatUseCase(e *bootstrap.Env, c domain.ChatRepository, b domain.BotRepository, ct task.AskTask, util *tokenutil.TokenUtil, ce domain.StorageEvent, tt task.TitleTask, cvt task.ConvertTask, gr domain.GenerationRepository) domain.ChatUseCase {
	chat := &chatUseCase{chatRepository: c, botRepository: b, env: e, chatTask: ct, tokenUtil: util, chatEvent: ce, titleTask: tt, convertTask: cvt, generationRepository: gr}
	return chat
}

func (c *chatUseCase) InitChat(ctx context.Context, token string, botId int) int {
	//ctx, cancel := context.WithTimeout(ctx, time.Duration(c.env.ContextTimeout))
	//defer cancel()

	script := incrementLuaScript
	log.GetJsonLogger().Info("load new chat lua script")

	chatId, err := c.chatRepository.CacheLuaInsertNewChatId(ctx, script, cache.NewestChatIdKey)
	if err != nil {
		log.GetTextLogger().Error(err.Error())
		return common.FalseInt
	}

	id, err := c.tokenUtil.DecodeToId(token)
	if err != nil {
		log.GetTextLogger().Error(err.Error())
		return common.FalseInt
	}

	// 同样提供依赖mq or not
	go c.chatRepository.DbInsertNewChat(ctx, id, botId)
	//c.chatEvent.PublishDbNewChat(&domain.ChatStorageData{BotId: botId, UserId: id})

	return chatId
}

func (c *chatUseCase) ContextChat(ctx context.Context, token string, botId int, chatId int, askMessage string, adjustment bool) (isSuccess bool, message domain.ParsedResponse, code int) {
	chatTask := c.chatTask

	userId, err := c.tokenUtil.DecodeToId(token)
	if err != nil {
		return false, &domain.OpenAIParsedResponse{GenerateText: common.ZeroString}, common.FalseInt
	}

	//我他妈太优雅了
	taskContext := chatTask.InitContextData(userId, botId, chatId, askMessage, task2.ExecuteChatAskType, task2.ExecuteChatAskCode, task2.ChatAskExecutorId, adjustment)
	factory := taskchain.NewTaskContextFactory()
	factory.TaskContext = taskContext
	factory.Puts(chatTask.PreCheckDataTask, chatTask.GetHistoryTask, chatTask.GetBotTask,
		chatTask.AssembleReqTask, chatTask.CallApiTask, chatTask.ParseRespTask, chatTask.StorageTask)
	factory.ExecuteChain()

	//按理来说 上面的taskContext == factory.TaskContext 但是下面再赋值一下比较稳妥一点
	taskContext = factory.TaskContext
	if taskContext.Exception {
		return false, &domain.OpenAIParsedResponse{GenerateText: taskContext.TaskContextResponse.Message}, taskContext.TaskContextResponse.Code
	}
	data := taskContext.TaskContextData.(*domain.AskContextData)
	parsedResponse := data.ParsedResponse

	response := parsedResponse.(*domain.OpenAIParsedResponse)
	return true, response, task2.SuccessCode
}

func (c *chatUseCase) StreamContextChatSetup(ctx context.Context, token string, botId int, chatId int, askMessage string, adjustment bool) (isSuccess bool, message domain.ParsedResponse, code int) {
	chatTask := c.chatTask
	convertTask := c.convertTask

	userId, err := c.tokenUtil.DecodeToId(token)
	if err != nil {
		return false, &domain.OpenAIParsedResponse{GenerateText: common.ZeroString}, common.FalseInt
	}

	taskContext := chatTask.InitContextData(userId, botId, chatId, askMessage, task2.ExecuteChatAskType, task2.ExecuteChatAskCode, task2.ChatAskExecutorId, adjustment)
	factory := taskchain.NewTaskContextFactory()
	factory.TaskContext = taskContext

	//StreamTask开启流式输出
	//至此组装好请求 向mq发布任务 mq消费 向指定客户端send generation
	//TODO remove
	factory.Puts(chatTask.PreCheckDataTask, chatTask.GetHistoryTask, chatTask.GetBotTask,
		convertTask.StreamArgsTask, chatTask.AssembleReqTask, convertTask.StreamStorageTask, chatTask.CallApiTask)
	//问题 我现在可以获取到用户的请求体 用以存储 但是如果想要存储的话 还需要生成 问题是 如何将生成的内容匹配请求体
	factory.ExecuteChain()

	taskContext = factory.TaskContext
	if taskContext.Exception {
		return false, &domain.OpenAIParsedResponse{GenerateText: taskContext.TaskContextResponse.Message}, taskContext.TaskContextResponse.Code
	}
	return true, nil, task2.SuccessCode
}

func (c *chatUseCase) StreamContextChatWorker(ctx context.Context, token string, gc *gin.Context, flusher http.Flusher) {
	userId, err := c.tokenUtil.DecodeToId(token)
	if err != nil {
		log.GetTextLogger().Error(err.Error())
		return
	}

	newSequencer := stream.NewSequencer()
	streamDataChan, _ := newSequencer.GetData(userId)
	log.GetTextLogger().Info("successfully getting the channel for: userId:" + strconv.Itoa(userId))

	var generateText string

	for {
		select {
		case v := <-streamDataChan:
			// 模拟数据推送
			marshal, _ := jsoniter.Marshal(v)
			// 发送符合SSE格式的数据到前端
			_, err = fmt.Fprintf(gc.Writer, "data: %s\n\n", marshal)
			if err != nil {
				log.GetTextLogger().Error(err.Error())
			}
			flusher.Flush() // 刷新输出到客户端

			generateText += v.GetGenerateText()

			if funk.NotEmpty(v.GetFinishReason()) {
				log.GetTextLogger().Info(fmt.Sprintf("Finish once push with finish reason " + v.GetFinishReason() + "  ,with chatcmplId:" + v.GetChatcmplId()))
				streamStorage := streamStorageMap[userId]
				if streamStorage == nil {
					streamStorage = make(chan string, 1)
				}
				streamStorage <- generateText
				return
			}
		case <-ctx.Done():
			// 上下文取消信号，优雅退出
			log.GetTextLogger().Info("Context canceled, stopping worker")

			streamStorage := streamStorageMap[userId]
			if streamStorage == nil {
				streamStorage = make(chan string, 1)
			}
			streamStorage <- generateText

			return
		default:
			//带超时的select语句 就算是在等待的时候 如果发生了事件
			//也可以立马响应 本质上是一个事件驱动架构的一个体现
			select {
			case <-time.After(time.Second):
			}
		}
	}
}

func (c *chatUseCase) StreamContextStorage(ctx context.Context, token string) bool {
	chatTask := c.chatTask
	convertTask := c.convertTask

	userId, err := c.tokenUtil.DecodeToId(token)
	if err != nil {
		log.GetTextLogger().Error(err.Error())
		return false
	}

	var generateText string

	streamResChan := streamStorageMap[userId]
	if streamResChan == nil {
		log.GetTextLogger().Error("cannot get target channel resource")
		return false
	} else {
		generateText = <-streamResChan
	}

	if funk.IsEmpty(generateText) {
		log.GetTextLogger().Error("empty generateText for userId: " + strconv.Itoa(userId))
		return false
	}

	taskContext := convertTask.InitStreamStorageTask(userId)

	factory := taskchain.NewTaskContextFactory()
	factory.TaskContext = taskContext

	storage := c.generationRepository.GetStreamDataStorage(context.Background(), userId)
	storage.ParsedResponse.SetGenerateText(generateText)
	factory.TaskContext.TaskContextData = storage

	factory.Puts(chatTask.StorageTask)

	log.GetTextLogger().Info("successfully storage stream message for userId: " + strconv.Itoa(userId))
	return true
}

func (c *chatUseCase) DisposableVisionChat(ctx context.Context, token string, chatId int, botId int, askMessage string, picUrl string) (isSuccess bool, message domain.ParsedResponse, code int) {
	chatTask := c.chatTask

	userId, err := c.tokenUtil.DecodeToId(token)
	if err != nil {
		return false, &domain.OpenAIParsedResponse{GenerateText: common.ZeroString}, common.FalseInt
	}

	//我他妈太优雅了
	taskContext := chatTask.InitContextData(userId, botId, chatId, picUrl, task2.ExecuteChatVisionAskType, task2.ExecuteChatVisionAskCode, task2.ChatVisionAskExecutorId)
	factory := taskchain.NewTaskContextFactory()
	factory.TaskContext = taskContext
	factory.Puts(chatTask.PreCheckDataTask, chatTask.GetBotTask,
		chatTask.AssembleReqTask, chatTask.CallApiTask, chatTask.ParseRespTask, chatTask.StorageTask)
	factory.ExecuteChain()

	taskContext = factory.TaskContext
	if taskContext.Exception {
		//return false, &domain.OpenAIParsedResponse{GenerateText: taskContext.TaskContextResponse.TextMessage}, taskContext.TaskContextResponse.Code
	}
	data := taskContext.TaskContextData.(*domain.AskContextData)
	parsedResponse := data.ParsedResponse

	response := parsedResponse.(*domain.OpenAIParsedResponse)
	return true, response, task2.SuccessCode
}

// func (c *chatUseCase) InitMainPage(ctx context.Context, token string) (titles []string, err error) {
// TODO 适应前端接口修改
func (c *chatUseCase) InitMainPage(ctx context.Context, token string, botId int) (titles []*domain.TitleData, err error) {
	userId, err := c.tokenUtil.DecodeToId(token)
	if err != nil {
		return nil, err
	}
	titleStr, err := c.chatRepository.CacheGetTitles(ctx, userId, botId)
	return titleStr, nil
}

// TODO 支持旧表
func (c *chatUseCase) GetChatHistory(ctx *gin.Context, chatId int, botId int, tokenString string) (*[]*domain.Record, error) {
	var history *[]*domain.Record
	history, isCache, err := c.chatRepository.CacheGetHistory(ctx, chatId, botId)
	if err != nil {
		return nil, err
	}
	userId, err := c.tokenUtil.DecodeToId(tokenString)
	if err != nil {
		log.GetTextLogger().Error(err.Error())
		return nil, err
	}
	var title string
	if isCache && funk.IsEmpty(history) {
		history, title, err = c.chatRepository.DbGetHistory(ctx, chatId, botId)
		if err != nil {
			log.GetTextLogger().Error(err.Error())
			return nil, err
		}

		// 回写缓存 (把从DB拿到的回写缓存 维护热点数据)
		c.chatRepository.CacheLuaLruResetHistory(context.Background(),
			cache.ChatHistoryScore+common.Infix+strconv.Itoa(userId)+common.Infix+strconv.Itoa(botId), history, chatId, title, botId)

	} else {
		return history, nil
	}
	return history, err
}

func (c *chatUseCase) GenerateUpdateTitle(ctx context.Context, message *[]domain.TextMessage, token string, chatId int) (string, error) {
	userId, err := c.tokenUtil.DecodeToId(token)
	if err != nil {
		log.GetTextLogger().Error(err.Error())
		return common.ZeroString, nil
	}

	taskContext := c.titleTask.InitContextData(userId, common.ZeroInt, chatId, message, task2.ExecuteTitleAskType, task2.ExecuteTitleAskCode, task2.ChatTitleAskExecutorId)
	factory := taskchain.NewTaskContextFactory()

	factory.TaskContext = taskContext
	titleTask := c.titleTask
	chatTask := c.chatTask
	factory.Puts(titleTask.PreTitleTask, titleTask.AssembleTitleReqTask,
		chatTask.CallApiTask, chatTask.ParseRespTask, chatTask.StorageTask)
	factory.ExecuteChain()

	//TODO 包装链子上出现的任务,继续提取其中共同点
	taskContext = factory.TaskContext
	if taskContext.Exception {
		e := errors.New(taskContext.TaskContextResponse.Message)
		return common.ZeroString, e
	}
	data := taskContext.TaskContextData.(*domain.AskContextData)
	parsedResponse := data.ParsedResponse

	response := parsedResponse.(*domain.OpenAIParsedResponse)
	return response.GenerateText, nil
}

func (c *chatUseCase) InputUpdateTitle(ctx context.Context, title string, token string, chatId int, botId int) bool {
	userId, err := c.tokenUtil.DecodeToId(token)
	if err != nil {
		log.GetTextLogger().Error("user signed error:")
		return false
	}
	data := &domain.AskContextData{ChatId: chatId, ParsedResponse: &domain.OpenAIParsedResponse{GenerateText: title}, UserId: userId}
	c.chatEvent.PublishDbSaveTitle(data)
	go c.chatRepository.CacheUpdateTitle(context.Background(), chatId, title, botId)
	return true
}
