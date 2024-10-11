package task

import (
	"SomersaultCloud/api/middleware/taskchain"
	"SomersaultCloud/domain"
)

type ChatConvertTask struct {
	generateEvent domain.GenerateEvent
}

// TODO remove
func (c ChatConvertTask) StreamPublishTask(tc *taskchain.TaskContext) {
	//由于包含不可用字段 所以此处不能使用消息队列发布任务
	c.generateEvent.PublishApiCalling(tc.TaskContextData.(*domain.AskContextData))
}

func (c ChatConvertTask) StreamArgsTask(tc *taskchain.TaskContext) {
	data := tc.TaskContextData.(*domain.AskContextData)
	data.Stream = true
}

func NewConvertTask(event domain.GenerateEvent) ConvertTask {
	return &ChatConvertTask{generateEvent: event}
}
