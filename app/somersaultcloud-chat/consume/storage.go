package consume

import (
	"SomersaultCloud/app/somersaultcloud-chat/constant/cache"
	"SomersaultCloud/app/somersaultcloud-chat/constant/common"
	"SomersaultCloud/app/somersaultcloud-chat/constant/dao"
	"SomersaultCloud/app/somersaultcloud-chat/constant/mq"
	"SomersaultCloud/app/somersaultcloud-chat/domain"
	"SomersaultCloud/app/somersaultcloud-chat/infrastructure/log"
	log2 "SomersaultCloud/app/somersaultcloud-common/log"
	"context"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

type storageEvent struct {
	*baseMessageHandler
	chatRepository domain.ChatRepository
}

type storageHistory struct {
	Records           *[]*domain.Record
	UserContent       string
	GenerationContent string
	ChatId            int
	UserId            int
	BotId             int
	Title             string
}

func storageDataReady(data *domain.AskContextData) *storageHistory {
	return &storageHistory{
		Records:           data.History,
		UserContent:       data.Message,
		GenerationContent: data.ParsedResponse.GetGenerateText(),
		ChatId:            data.ChatId,
		UserId:            data.UserId,
		BotId:             data.BotId,
		Title:             data.ParsedResponse.GetGenerateText(),
	}
}

func (c storageEvent) DbPutHistory(b []byte) error {
	var data storageHistory
	_ = jsoniter.Unmarshal(b, &data)
	c.chatRepository.AsyncSaveHistory(context.Background(),
		data.ChatId,
		data.UserContent,
		data.GenerationContent,
		data.BotId,
	)
	return nil
}

func (c storageEvent) PublishSaveDbHistory(data *domain.AskContextData) {
	dataReady := storageDataReady(data)
	marshal, _ := jsoniter.Marshal(dataReady)

	c.PublishMessage(mq.HistoryDbSaveQueue, marshal)
}
func (c storageEvent) AsyncConsumeDbHistory() {
	c.ConsumeMessage(mq.HistoryDbSaveQueue, c.DbPutHistory)
}
func (c storageEvent) CachePutHistory(b []byte) error {
	var data storageHistory
	_ = jsoniter.Unmarshal(b, &data)
	err := c.chatRepository.CacheLuaLruPutHistory(context.Background(),
		cache.ChatHistoryScore+common.Infix+strconv.Itoa(data.UserId)+common.Infix+strconv.Itoa(data.BotId),
		data.Records,
		data.UserContent,
		data.GenerationContent,
		data.ChatId,
		data.BotId,
		dao.DefaultTitle)
	if err != nil {
		log2.GetTextLogger().Error("mq cache put history error:" + err.Error())
	}
	return err
}

func (c storageEvent) PublishSaveCacheHistory(data *domain.AskContextData) {
	dataReady := storageDataReady(data)
	marshal, _ := jsoniter.Marshal(dataReady)
	c.PublishMessage(mq.HistoryCacheSaveQueue, marshal)
}
func (c storageEvent) AsyncConsumeCacheHistory() {
	c.ConsumeMessage(mq.HistoryCacheSaveQueue, c.CachePutHistory)
}

func (c storageEvent) DbNewChat(b []byte) error {
	var data domain.ChatStorageData
	err := jsoniter.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	c.chatRepository.DbInsertNewChat(context.Background(), data.UserId, data.BotId)
	return nil
}
func (c storageEvent) PublishDbNewChat(data *domain.ChatStorageData) {
	marshal, _ := jsoniter.Marshal(data)
	log2.GetJsonLogger().WithFields("user", data.UserId, "activity", "createNewChat").Info("create new chat")
	c.PublishMessage(mq.InsertNewChatQueue, marshal)
}
func (c storageEvent) AsyncConsumeDbNewChat() {
	c.ConsumeMessage(mq.InsertNewChatQueue, c.DbNewChat)
}

func (c storageEvent) PublishDbSaveTitle(data *domain.AskContextData) {
	dataReady := storageDataReady(data)
	marshal, _ := jsoniter.Marshal(dataReady)
	log2.GetJsonLogger().WithFields("user", data.UserId, "activity", "update chat title", "chatId", data.ChatId)
	c.PublishMessage(mq.UpdateChatTitleQueue, marshal)
}
func (c storageEvent) DbUpdateTitle(b []byte) error {
	var data storageHistory
	_ = jsoniter.Unmarshal(b, &data)
	c.chatRepository.DbUpdateTitle(context.Background(), data.ChatId, data.Title)
	return nil
}

func (c storageEvent) AsyncConsumeDbUpdateTitle() {
	c.ConsumeMessage(mq.UpdateChatTitleQueue, c.DbUpdateTitle)
}

func NewStorageEvent(c domain.ChatRepository, h MessageHandler) domain.StorageEvent {
	//TODO 丑
	handler := h.(*baseMessageHandler)
	dbSave := &MessageQueueArgs{
		ExchangeName: mq.HistoryDbSaveExchange,
		QueueName:    mq.HistoryDbSaveQueue,
		KeyName:      mq.HistoryDbSaveKey,
	}
	cacheSave := &MessageQueueArgs{
		ExchangeName: mq.HistoryCacheSaveExchange,
		QueueName:    mq.HistoryCacheSaveQueue,
		KeyName:      mq.HistoryCacheSaveKey,
	}
	dbNewChat := &MessageQueueArgs{
		ExchangeName: mq.InsertNewChatExchange,
		QueueName:    mq.InsertNewChatQueue,
		KeyName:      mq.InsertNewChatKey,
	}
	dbNewTitle := &MessageQueueArgs{
		ExchangeName: mq.UpdateChatTitleExchange,
		QueueName:    mq.UpdateChatTitleQueue,
		KeyName:      mq.UpdateChatTitleKey,
	}
	handler.InitMessageQueue(dbSave, cacheSave, dbNewChat, dbNewTitle)
	return &storageEvent{
		baseMessageHandler: handler,
		chatRepository:     c,
	}
}
