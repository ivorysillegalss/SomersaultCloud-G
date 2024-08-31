package usecase

import (
	"SomersaultCloud/api/middleware/taskchain"
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	task2 "SomersaultCloud/constant/task"
	"SomersaultCloud/domain"
	"SomersaultCloud/internal/ioutil"
	"SomersaultCloud/internal/tokenutil"
	"SomersaultCloud/task"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type chatUseCase struct {
	env            *bootstrap.Env
	chatRepository domain.ChatRepository
	botRepository  domain.BotRepository
	chatTask       task.AskTask
	tokenUtil      *tokenutil.TokenUtil
	chatEvent      domain.ChatEvent
}

func NewChatUseCase(e *bootstrap.Env, c domain.ChatRepository, b domain.BotRepository, ct task.AskTask, util *tokenutil.TokenUtil, ce domain.ChatEvent) domain.ChatUseCase {
	chat := &chatUseCase{chatRepository: c, botRepository: b, env: e, chatTask: ct, tokenUtil: util, chatEvent: ce}
	return chat
}

func (c *chatUseCase) InitChat(ctx context.Context, token string, botId int) int {
	//ctx, cancel := context.WithTimeout(ctx, time.Duration(c.env.ContextTimeout))
	//defer cancel()

	script, err := ioutil.LoadLuaScript("usecase/lua/increment.lua")
	if err != nil {
		return common.FalseInt
	}

	chatId, err := c.chatRepository.CacheLuaInsertNewChatId(ctx, script, cache.NewestChatIdKey)
	if err != nil {
		return common.FalseInt
	}

	id, err := c.tokenUtil.DecodeToId(token)
	if err != nil {
		return common.FalseInt
	}

	// 同样提供依赖mq or not
	//go c.chatRepository.DbInsertNewChat(ctx, id, botId)
	c.chatEvent.PublishDbNewChat(&domain.ChatStorageData{BotId: botId, UserId: id})

	return chatId
}

func (c *chatUseCase) ContextChat(ctx context.Context, token string, botId int, chatId int, askMessage string) (isSuccess bool, message domain.ParsedResponse, code int) {
	chatTask := c.chatTask

	userId, err := c.tokenUtil.DecodeToId(token)
	if err != nil {
		return false, &domain.OpenAIParsedResponse{GenerateText: common.ZeroString}, common.FalseInt
	}

	//我他妈太优雅了
	taskContext := chatTask.InitContextData(userId, botId, chatId, askMessage)
	factory := taskchain.NewTaskContextFactory()
	factory.TaskContext = taskContext
	factory.Puts(chatTask.PreCheckDataTask, chatTask.GetHistoryTask, chatTask.GetBotTask,
		chatTask.AssembleReqTask, chatTask.CallApiTask, chatTask.ParseRespTask)
	factory.ExecuteChain()
	//TODO 异步数据缓存

	// TODO 接入消息队列

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

// func (c *chatUseCase) InitMainPage(ctx context.Context, token string) (titles []string, err error) {
// TODO 适应前端接口修改
func (c *chatUseCase) InitMainPage(ctx context.Context, token string) (titles []*domain.TitleData, err error) {
	userId, err := c.tokenUtil.DecodeToId(token)
	if err != nil {
		return nil, err
	}
	titleStr, err := c.chatRepository.CacheGetTitles(ctx, userId)
	return titleStr, nil
}

func (c *chatUseCase) GetChatHistory(ctx *gin.Context, chatId int) (*[]*domain.Record, error) {
	var history *[]*domain.Record
	history, isCache, err := c.chatRepository.CacheGetHistory(ctx, chatId)
	if err != nil {
		return nil, err
	}
	if isCache && funk.IsEmpty(history) {
		history, _, err = c.chatRepository.DbGetHistory(ctx, chatId)
	} else {
		return history, nil
	}
	return history, err
}
