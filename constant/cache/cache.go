package cache

import "time"

const (
	// NewestChatIdKey 最新chatId
	NewestChatIdKey = "newestChatId"
	// BotConfig 模型相关配置
	BotConfig = "botConfig"
	// ChatHistoryScore 历史记录LRU权重
	ChatHistoryScore = "chatHistoryLruScore"
	// ChatHistory 历史记录缓存前缀
	ChatHistory = "chatHistory"
	// ChatHistoryTitle 历史记标题前缀
	ChatHistoryTitle = "chatHistoryTitle"
	// MaxBotId 最大的BotId 用于判断数据是否合法
	MaxBotId = "maxBotId"

	// OriginTable 原表缓存 由于两种历史记录的格式不一样 旧表用这个进行过渡
	OriginTable = "ori"

	// RedisListType List数据类型标识符
	RedisListType = 1
	// RedisZSetType ZSet数据类型标识符
	RedisZSetType = 2

	// ContextLruMaxCapacity 历史记录LRU窗口默认大小
	ContextLruMaxCapacity = 5
	// HistoryDefaultWeight 单chat历史记录默认存储记录大小
	HistoryDefaultWeight = 5
	// LruPrefix LRU前缀缓存
	LruPrefix = "lru"
	// ChatGeneration chat的生成缓存
	ChatGeneration = "chatGeneration"
	// ChatGenerationExpired chat生成缓存时间 配合lua脚本设定HSet单键DDL
	ChatGenerationExpired = "chatGenerationExpired"
	// ChatGenerationTTL 生成缓存时间
	ChatGenerationTTL = 500
	// HistoryTitlePrompt 获取标题的prompt
	HistoryTitlePrompt = "historyTitlePrompt"

	// StreamStorageReadyData 提前缓存流信息请求时数据
	StreamStorageReadyData = "StreamStorageReadyData"
	// StreamStorageReadyDataExpire 缓存信息保存的DDL
	StreamStorageReadyDataExpire = 5 * time.Second
)
