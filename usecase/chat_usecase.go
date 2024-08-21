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
)

type chatUseCase struct {
	env            *bootstrap.Env
	chatRepository domain.ChatRepository
	botRepository  domain.BotRepository
	chatTask       task.AskTask
	tokenUtil      *tokenutil.TokenUtil
}

func NewChatUseCase(e *bootstrap.Env, c domain.ChatRepository, b domain.BotRepository, ct task.AskTask, util *tokenutil.TokenUtil) domain.ChatUseCase {
	chat := &chatUseCase{chatRepository: c, botRepository: b, env: e, chatTask: ct, tokenUtil: util}
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
	go c.chatRepository.DbInsertNewChatId(ctx, id, botId)
	// TODO mq异步写入MYSQL

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
