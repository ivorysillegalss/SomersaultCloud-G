package usecase

import (
	"SomersaultCloud/api/dto"
	"SomersaultCloud/api/middleware/taskchain"
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	task2 "SomersaultCloud/constant/task"
	"SomersaultCloud/domain"
	"SomersaultCloud/internal/ioutil"
	"SomersaultCloud/internal/tokenutil"
	"SomersaultCloud/repository"
	"SomersaultCloud/task"
	"context"
	"time"
)

type chatUseCase struct {
	chatRepository domain.ChatRepository
	botRepository  domain.BotRepository
	chatTask       task.AskTask
}

func NewChatUseCase() domain.ChatUseCase {
	chat := &chatUseCase{chatRepository: repository.NewChatRepository(), botRepository: repository.NewBotRepository()}
	chat.chatTask = task.NewAskChatTask(chat.botRepository, chat.chatRepository)
	return chat
}

func (c *chatUseCase) InitChat(ctx context.Context, token string, botId int) int {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(env.ContextTimeout))
	defer cancel()

	script, err := ioutil.LoadLuaScript("lua/increment.lua")
	if err != nil {
		return common.FalseInt
	}

	chatId, err := c.chatRepository.CacheLuaInsertNewChatId(ctx, script, cache.NewestChatIdKey)
	if err != nil {
		return common.FalseInt
	}

	id, err := tokenutil.DecodeToId(token)
	if err != nil {
		return common.FalseInt
	}
	go c.chatRepository.DbInsertNewChatId(ctx, id, botId)
	// TODO mq异步写入MYSQL

	return chatId
}

func (c *chatUseCase) ContextChat(ctx context.Context, token string, ask *dto.AskDTO) (isSuccess bool, message domain.ParsedResponse, code int) {
	chatTask := c.chatTask

	userId, err := tokenutil.DecodeToId(token)
	if err != nil {
		return false, &domain.OpenAIParsedResponse{GenerateText: common.ZeroString}, common.FalseInt
	}
	ask.UserId = userId

	//我他妈太优雅了
	taskContext := chatTask.InitContextData()
	factory := taskchain.NewTaskContextFactory()
	factory.TaskContext = taskContext
	taskContext.TData = ask
	factory.Puts(chatTask.PreCheckDataTask, chatTask.GetHistoryTask, chatTask.GetBotTask,
		chatTask.AssembleReqTask, chatTask.CallApiTask, chatTask.ParseRespTask)

	//TODO 异步数据缓存

	// TODO 接入消息队列

	//按理来说 上面的taskContext == factory.TaskContext 但是下面再赋值一下比较稳妥一点
	taskContext = factory.TaskContext
	if taskContext.Exception {
		return false, &domain.OpenAIParsedResponse{GenerateText: taskContext.TaskContextResponse.Message}, taskContext.TaskContextResponse.Code
	}
	parsedResponse := taskContext.TaskContextData.ParsedResponse
	response := parsedResponse.(*domain.OpenAIParsedResponse)
	return true, response, task2.SuccessCode
}
