package cache

const (
	// NewestChatIdKey 最新chatId
	NewestChatIdKey = "newestChatId"
	// BotConfig 模型相关配置
	BotConfig = "botConfig"
	// ChatHistory 历史记录
	ChatHistory = "chatHistory"
	// MaxBotId 最大的BotId 用于判断数据是否合法
	MaxBotId = "maxBotId"
	// ListType List数据类型标识符
	ListType = 1
	// ContextLruMaxCapacity 历史记录LRU窗口默认大小
	ContextLruMaxCapacity = 5
)
