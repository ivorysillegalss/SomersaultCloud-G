package executor

import (
	"SomersaultCloud/app/somersaultcloud-chat/domain"
	log2 "SomersaultCloud/app/somersaultcloud-common/log"
)

type ConsumeExecutor struct {
	storageEvent  domain.StorageEvent
	generateEvent domain.GenerateEvent
}

func (d *ConsumeExecutor) SetupConsume() {
	d.storageEvent.AsyncConsumeDbHistory()
	log2.GetTextLogger().Info("AsyncConsumeDbHistory QUEUE start")
	d.storageEvent.AsyncConsumeCacheHistory()
	log2.GetTextLogger().Info("AsyncConsumeCacheHistory QUEUE start")
	d.storageEvent.AsyncConsumeDbUpdateTitle()
	log2.GetTextLogger().Info("AsyncConsumeDbUpdateTitle QUEUE start")
	d.generateEvent.AsyncStreamStorageDataReady()
	log2.GetTextLogger().Info("AsyncStreamStorageDataReady QUEUE start")

	log2.GetTextLogger().Info("ALL-----QUEUE----START-----SUCCESSFULLY")
	//TODO
	//在这里全部启动消费者逻辑
}

func NewConsumeExecutor(c domain.StorageEvent, g domain.GenerateEvent) *ConsumeExecutor {
	return &ConsumeExecutor{storageEvent: c, generateEvent: g}
}
