package task

import (
	"SomersaultCloud/api/dto"
	"SomersaultCloud/api/middleware/taskchain"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/task"
	"SomersaultCloud/domain"
	"SomersaultCloud/internal/checkutil"
	"SomersaultCloud/repository"
	"context"
	"github.com/thoas/go-funk"
)

// 责任链任务实现
type chatTask struct {
	chatRepository domain.ChatRepository
	botRepository  domain.BotRepository
}

type AskContextData struct {
	chatId  int
	userId  int
	message string
	botId   int
	history *[]*domain.Record
}

func (c *chatTask) PreCheckDataTask(tc *taskchain.TaskContext) {
	askDTO := tc.TData.(dto.AskDTO)
	ask := askDTO.Ask
	//TODO 运行前redis加缓存
	chatIdCheck := checkutil.IsLegalID(ask.ChatId, common.FalseInt, c.chatRepository.CacheGetNewestChatId(context.Background()))
	botIdCheck := checkutil.IsLegalID(ask.BotId, common.FalseInt, c.botRepository.CacheGetMaxBotId(context.Background()))
	message := ask.Message
	msgCheck := funk.NotEmpty(message)

	if !(msgCheck && chatIdCheck && botIdCheck) {
		tc.InterruptExecute(task.InvalidDataFormatMessage)
		return
	}
	tc.TaskContextData.botId = ask.BotId
	tc.TaskContextData.chatId = ask.ChatId
	tc.TaskContextData.userId = askDTO.UserId
	tc.TaskContextData.message = ask.Message
}

// GetHistoryTask 2情况 判断是否存在缓存 hit拿缓存 miss则db
func (c *chatTask) GetHistoryTask(tc taskchain.TaskContext) {
	var history *[]*domain.Record
	// 1. 缓存找
	history, isCache, err := c.chatRepository.CacheGetHistory(context.Background(), tc.TaskContextData.chatId)
	if err != nil {
		tc.InterruptExecute(task.HistoryRetrievalFailed)
		return
	}

	// 2. 缓存miss db找
	//TODO 目前查DB后需要截取历史记录 实现数据流式更新后可取消
	if isCache {
		history, err = c.chatRepository.DbGetHistory(context.Background(), tc.TaskContextData.chatId)

		// 截取数据
		if len(*history) >= common.HistoryDefaultWeight {
			*history = (*history)[:common.HistoryDefaultWeight]
		}

	}

	// 2.1 回写缓存
	//TODO
	if err != nil {
		tc.InterruptExecute(task.HistoryRetrievalFailed)
		return
	}

	tc.TaskContextData.history = history
}

func (c *chatTask) AdjustmentTask(tc taskchain.TaskContext) {
	//TODO implement me
	panic("implement me")
}

func (c *chatTask) AssembleReqTask(tc *taskchain.TaskContext) {

}

func (c *chatTask) CallApiTask(tc *taskchain.TaskContext) {

}
func (c *chatTask) ParseRespTask(tc *taskchain.TaskContext) {

}
func NewChatTask() domain.ChatTask {
	return &chatTask{chatRepository: repository.NewChatRepository(), botRepository: repository.NewBotRepository()}
}
