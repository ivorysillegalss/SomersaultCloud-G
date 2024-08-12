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

	// RedisListType List数据类型标识符
	RedisListType = 1
	// RedisZSetType ZSet数据类型标识符
	RedisZSetType = 2

	// ContextLruMaxCapacity 历史记录LRU窗口默认大小
	ContextLruMaxCapacity = 5
	// LruPrefix LRU前缀缓存
	LruPrefix = "lru"
)
