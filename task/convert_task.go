package task

import (
	"SomersaultCloud/api/middleware/taskchain"
	"SomersaultCloud/domain"
)

type ChatConvertTask struct {
	chatEvent domain.GenerateEvent
}

func (c ChatConvertTask) StreamPublishTask(tc *taskchain.TaskContext) {
	c.chatEvent.PublishApiCalling(tc.TaskContextData.(*domain.AskContextData))
}

func (c ChatConvertTask) StreamArgsTask(tc *taskchain.TaskContext) {
	data := tc.TaskContextData.(*domain.AskContextData)
	data.Stream = true
}

func NewConvertTask(event domain.StorageEvent) ConvertTask {
	return ChatConvertTask{chatEvent: event}
}
