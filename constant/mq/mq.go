package mq

const (
	// RabbitMqReconnectDelay reconnect after delay seconds
	RabbitMqReconnectDelay = 3

	// MqPublishErr Mq相关的处理错误
	MqPublishErr = "PUBLISH ERR"
	MqConsumeErr = "CONSUME ERR"

	InsertNewChatExchange = "user.new.history.db.direct"
	InsertNewChatQueue    = "user.new.history.db.queue"
	InsertNewChatKey      = "user.new.history.db.event"

	HistoryCacheSaveExchange = "user.save.history.cache.direct"
	HistoryCacheSaveQueue    = "user.save.history.cache.queue"
	HistoryCacheSaveKey      = "user.save.history.cache.event"

	HistoryDbSaveExchange = "user.save.history.db.direct"
	HistoryDbSaveQueue    = "user.save.history.db.queue"
	HistoryDbSaveKey      = "user.save.history.db.event"

	UpdateChatTitleExchange = "user.update.title.direct"
	UpdateChatTitleQueue    = "user.update.title.queue"
	UpdateChatTitleKey      = "user.update.title.event"

	UserChatReadyCallingExchange = "user.call.chat.ready.direct"
	UserChatReadyCallingQueue    = "user.call.chat.ready.queue"
	UserChatReadyCallingKey      = "user.call.chat.ready.event"
)
