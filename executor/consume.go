package executor

import (
	"SomersaultCloud/domain"
)

type ConsumeExecutor struct {
	ChatEvent domain.ChatEvent
}

func (d *ConsumeExecutor) SetupConsume() {
	d.ChatEvent.AsyncConsumeDbHistory()
	d.ChatEvent.AsyncConsumeCacheHistory()
	//TODO
	//在这里全部启动消费者逻辑
}

func NewConsumeExecutor(c domain.ChatEvent) *ConsumeExecutor {
	return &ConsumeExecutor{ChatEvent: c}
}
