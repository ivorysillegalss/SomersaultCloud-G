package domain

import (
	"context"
)

type Chat struct {
	ID             int        `json:"chat_id"  gorm:"primaryKey"`
	UserId         int        `json:"user_id"`
	BotId          int        `json:"bot_id"`
	Title          string     `json:"title"`
	LastUpdateTime int64      `json:"last_update_time"`
	IsDelete       bool       `json:"is_delete"`
	Records        *[]*Record `json:"records" gorm:"-"`
	Data           []byte     `json:"data"`
}

type Record struct {
	//ChatId int `json:"chat_id"`
	//RecordId        int             `json:"record_id"`
	ChatAsks        *ChatAsk        `json:"chat_asks"`
	ChatGenerations *ChatGeneration `json:"chat_generations"`
	//Weights         float64
}

// ChatAsk 一次问题
type ChatAsk struct {
	//RecordId int    `json:"record_id"`
	ChatId  int    `json:"chat_id,omitempty" gorm:"-"`
	Message string `json:"message"`
	BotId   int    `json:"bot_id,omitempty" gorm:"-"`
	Time    int64  `json:"time"`
}

// ChatGeneration 一次生成
type ChatGeneration struct {
	Message string `json:"message"`
	Time    int64  `json:"time"`
}

type ChatRepository interface {
	// CacheGetNewestChatId 获取最新chatId 不能保证原子性 弃用
	CacheGetNewestChatId(ctx context.Context) int
	// CacheInsertNewChat 增加新Id 不能保证原子性 弃用
	CacheInsertNewChat(ctx context.Context, id int)

	// CacheLuaInsertNewChatId lua脚本保证高并发时获取chatId的一致性
	CacheLuaInsertNewChatId(ctx context.Context, luaScript string, k string) (int, error)
	// DbInsertNewChat 异步使用 存入SQL持久化方法
	DbInsertNewChat(ctx context.Context, userId int, botId int)

	// CacheGetHistory 从缓存中取出历史记录 存的时候确保最大条数 取时无需注意
	CacheGetHistory(ctx context.Context, chatId int) (*[]*Record, bool, error)
	// DbGetHistory miss缓存 从DB中获取历史记录
	DbGetHistory(ctx context.Context, chatId int) (*[]*Record, error)

	// AsyncSaveHistory 异步保存历史记录
	AsyncSaveHistory(ctx context.Context, chatId int, askText string, generationText string)
	// CacheLuaLruResetHistory 这个是在生成前 把从DB拿到的数据回写缓存 维护热点数据
	CacheLuaLruResetHistory(ctx context.Context, cacheKey string, history *[]*Record, chatId int) error
	// CacheLuaLruPutHistory 这个是在生成完毕后 回写完整
	CacheLuaLruPutHistory(ctx context.Context, cacheKey string, history *[]*Record, askText string, generationText string, chatId int) error

	//由于http.response对象不可序列化 转为inmemory存储
	MemoryGetGeneration(ctx context.Context, chatId int) *GenerationResponse
	CacheGetGeneration(ctx context.Context, chatId int) (*GenerationResponse, error)
	MemoryDelGeneration(ctx context.Context, chatId int)
	CacheDelGeneration(ctx context.Context, chatId int) error
}

type ChatUseCase interface {
	InitChat(ctx context.Context, token string, botId int) int
	ContextChat(ctx context.Context, token string, botId int, chatId int, askMessage string) (isSuccess bool, message ParsedResponse, code int)
}

type ChatEvent interface {
	//对repository中方法进行二次封装
	DbPutHistory(b []byte) error
	CachePutHistory(b []byte) error
	DbNewChat(b []byte) error

	PublishSaveDbHistory(data *AskContextData)
	PublishSaveCacheHistory(data *AskContextData)
	PublishDbNewChat(data *ChatStorageData)

	AsyncConsumeDbHistory()
	AsyncConsumeCacheHistory()
	AsyncConsumeDbNewChat()
}

type ChatStorageData struct {
	UserId int
	ChatId int
	BotId  int
}
