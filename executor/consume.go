package executor

import (
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/log"
)

type ConsumeExecutor struct {
	ChatEvent domain.ChatEvent
}

func (d *ConsumeExecutor) SetupConsume() {
	d.ChatEvent.AsyncConsumeDbHistory()
	log.GetTextLogger().Info("AsyncConsumeDbHistory QUEUE start")
	d.ChatEvent.AsyncConsumeCacheHistory()
	log.GetTextLogger().Info("AsyncConsumeCacheHistory QUEUE start")
	//TODO
	//在这里全部启动消费者逻辑
}

func NewConsumeExecutor(c domain.ChatEvent) *ConsumeExecutor {
	return &ConsumeExecutor{ChatEvent: c}
}