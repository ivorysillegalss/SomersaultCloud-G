package executor

import (
	"SomersaultCloud/app/domain"
	"SomersaultCloud/app/infrastructure/log"
)

type ConsumeExecutor struct {
	storageEvent  domain.StorageEvent
	generateEvent domain.GenerateEvent
}

func (d *ConsumeExecutor) SetupConsume() {
	d.storageEvent.AsyncConsumeDbHistory()
	log.GetTextLogger().Info("AsyncConsumeDbHistory QUEUE start")
	d.storageEvent.AsyncConsumeCacheHistory()
	log.GetTextLogger().Info("AsyncConsumeCacheHistory QUEUE start")
	d.storageEvent.AsyncConsumeDbUpdateTitle()
	log.GetTextLogger().Info("AsyncConsumeDbUpdateTitle QUEUE start")
	d.generateEvent.AsyncStreamStorageDataReady()
	log.GetTextLogger().Info("AsyncStreamStorageDataReady QUEUE start")

	log.GetTextLogger().Info("ALL-----QUEUE----START-----SUCCESSFULLY")
	//TODO
	//在这里全部启动消费者逻辑
}

func NewConsumeExecutor(c domain.StorageEvent, g domain.GenerateEvent) *ConsumeExecutor {
	return &ConsumeExecutor{storageEvent: c, generateEvent: g}
}
