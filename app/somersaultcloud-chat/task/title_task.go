package task

import (
	"SomersaultCloud/app/somersaultcloud-chat/api/middleware/taskchain"
	"SomersaultCloud/app/somersaultcloud-chat/bootstrap"
	"SomersaultCloud/app/somersaultcloud-chat/constant/common"
	"SomersaultCloud/app/somersaultcloud-chat/constant/dao"
	"SomersaultCloud/app/somersaultcloud-chat/constant/task"
	"SomersaultCloud/app/somersaultcloud-chat/domain"
	"SomersaultCloud/app/somersaultcloud-chat/handler"
	"SomersaultCloud/app/somersaultcloud-chat/internal/checkutil"
	"context"
	"github.com/thoas/go-funk"
)

type ChatTitleTask struct {
	chatRepository domain.ChatRepository
	env            *bootstrap.Env
	channels       *bootstrap.Channels
}

func NewChatTitleTask(c domain.ChatRepository, env *bootstrap.Env, cn *bootstrap.Channels) TitleTask {
	return &ChatTitleTask{chatRepository: c, env: env, channels: cn}
}

func (c ChatTitleTask) InitContextData(args ...any) *taskchain.TaskContext {
	userId := args[0].(int)
	chatId := args[2].(int)
	message := args[3].(*[]domain.TextMessage)
	return &taskchain.TaskContext{
		BusinessType:    args[4].(string),
		BusinessCode:    args[5].(int),
		TaskContextData: &domain.AskContextData{UserId: userId, ChatId: chatId, HistoryMessage: &domain.Message{TextMessage: message}, ExecutorId: args[6].(int)},
	}
}

func (c ChatTitleTask) PreTitleTask(tc *taskchain.TaskContext) {
	data := tc.TaskContextData.(*domain.AskContextData)
	chatIdCheck := checkutil.IsLegalID(data.ChatId, common.FalseInt, c.chatRepository.CacheGetNewestChatId(context.Background()))
	message := data.HistoryMessage.TextMessage
	msgCheck := funk.NotEmpty(message)
	if !(msgCheck || chatIdCheck) {
		tc.InterruptExecute(task.InvalidDataFormatMessage)
		return
	}
}

func (c ChatTitleTask) AssembleTitleReqTask(tc *taskchain.TaskContext) {
	data := tc.TaskContextData.(*domain.AskContextData)
	executor := handler.NewLanguageModelExecutor(c.env, c.channels, data.ExecutorId)
	data.Executor = executor
	data.Model = dao.DefaultModel

	titleSysPrompt := c.chatRepository.CacheGetTitlePrompt(context.Background())
	message := *data.HistoryMessage.TextMessage
	SysMessage := &domain.TextMessage{Role: common.SystemRole, Content: titleSysPrompt}
	// 将 SysMessage 放到 message 的前面
	message = append([]domain.TextMessage{*SysMessage}, message...)
	data.HistoryMessage.TextMessage = &message

	request := executor.EncodeReq(data)
	if funk.IsEmpty(request) {
		tc.InterruptExecute(task.ReqDataMarshalFailed)
		return
	}
	client := executor.ConfigureProxy(data)
	data.Conn = *domain.NewConnection(client, request)
}

// Convert2AskTask 预备方法,防止不兼容
func (c ChatTitleTask) Convert2AskTask(tc *taskchain.TaskContext) {
}
