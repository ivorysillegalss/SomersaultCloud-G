package task

import (
	"SomersaultCloud/app/api/middleware/taskchain"
	task2 "SomersaultCloud/app/constant/task"
	"SomersaultCloud/app/domain"
)

type ChatConvertTask struct {
	generateEvent        domain.GenerateEvent
	generationRepository domain.GenerationRepository
}

func (c ChatConvertTask) InitStreamStorageTask(args ...any) *taskchain.TaskContext {
	userId := args[0].(int)
	return &taskchain.TaskContext{BusinessType: task2.StorageStreamType, BusinessCode: task2.StorageStreamCode, TaskContextData: &domain.AskContextData{UserId: userId}}
}

func (c ChatConvertTask) StreamArgsTask(tc *taskchain.TaskContext) {
	data := tc.TaskContextData.(*domain.AskContextData)
	data.Stream = true
}

func (c ChatConvertTask) StreamStorageTask(tc *taskchain.TaskContext) {
	data := tc.TaskContextData.(*domain.AskContextData)
	dataReady := &domain.StreamGenerationReadyStorageData{ChatId: data.ChatId, UserContent: data.Message, UserId: data.UserId, BotId: data.BotId, Records: data.History}
	c.generateEvent.PublishStreamReadyStorageData(dataReady)
}

func NewConvertTask(event domain.GenerateEvent, g domain.GenerationRepository) ConvertTask {
	return &ChatConvertTask{generateEvent: event, generationRepository: g}
}
