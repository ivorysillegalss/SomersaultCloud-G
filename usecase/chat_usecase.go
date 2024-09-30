package usecase

import (
	"SomersaultCloud/api/middleware/taskchain"
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	task2 "SomersaultCloud/constant/task"
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/log"
	"SomersaultCloud/internal/tokenutil"
	"SomersaultCloud/task"
	"context"
	_ "embed"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"strconv"
)

//go:embed lua/increment.lua
var incrementLuaScript string

type chatUseCase struct {
	env            *bootstrap.Env
	chatRepository domain.ChatRepository
	botRepository  domain.BotRepository
	chatTask       task.AskTask
	tokenUtil      *tokenutil.TokenUtil
	chatEvent      domain.ChatEvent
	titleTask      task.TitleTask
}

func NewChatUseCase(e *bootstrap.Env, c domain.ChatRepository, b domain.BotRepository, ct task.AskTask, util *tokenutil.TokenUtil, ce domain.ChatEvent, tt task.TitleTask) domain.ChatUseCase {
	chat := &chatUseCase{chatRepository: c, botRepository: b, env: e, chatTask: ct, tokenUtil: util, chatEvent: ce, titleTask: tt}
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
		chatTask.AssembleReqTask, chatTask.CallApiTask, chatTask.ParseRespTask)
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

func (c *chatUseCase) DisposableVisionChat(ctx context.Context, token string, chatId int, botId int, askMessage string, picUrl string) (isSuccess bool, message domain.ParsedResponse, code int) {
	chatTask := c.chatTask

	userId, err := c.tokenUtil.DecodeToId(token)
	//TODO
	if err != nil {
		return false, &domain.OpenAIParsedResponse{GenerateText: common.ZeroString}, common.FalseInt
	}

	//我他妈太优雅了
	taskContext := chatTask.InitContextData(userId, botId, chatId, picUrl, task2.ExecuteChatVisionAskType, task2.ExecuteChatVisionAskCode, task2.ChatVisionAskExecutorId)
	factory := taskchain.NewTaskContextFactory()
	factory.TaskContext = taskContext
	factory.Puts(chatTask.PreCheckDataTask, chatTask.GetBotTask,
		chatTask.AssembleReqTask, chatTask.CallApiTask, chatTask.ParseRespTask)
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
		chatTask.CallApiTask, chatTask.ParseRespTask)
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
