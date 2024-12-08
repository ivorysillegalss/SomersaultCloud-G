package domain

import (
	"context"
	"github.com/hertz-contrib/sse"
)

type ChatHistoryTitle struct {
	Title []string `json:"title"`
}

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
	RecordId        int             `json:"record_id"`
	ChatAsks        *ChatAsk        `json:"chat_asks"`
	ChatGenerations *ChatGeneration `json:"chat_generations"`
	//Weights         float64
}

//	TODO 加一个imageurllist，应付提问中有多张图片
//
// ChatAsk 一次问题
type ChatAsk struct {
	//TODO 兼容旧表
	RecordId int `json:"record_id"`
	//ChatId   int    `json:"chat_id,omitempty" gorm:"-"`
	ChatId  int    `json:"chat_id,omitempty"`
	Message string `json:"message"`
	BotId   int    `json:"bot_id,omitempty" gorm:"-"`
	Time    int64  `json:"time"`
}

// ChatGeneration 一次生成
type ChatGeneration struct {
	//TODO 兼容旧表
	ChatId   int    `json:"chat_id,omitempty"`
	RecordId int    `json:"record_id"`
	Message  string `json:"message"`
	Time     int64  `json:"time"`
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
	CacheGetHistory(ctx context.Context, chatId int, botId int) (*[]*Record, bool, error)
	// DbGetHistory miss缓存 从DB中获取历史记录
	DbGetHistory(ctx context.Context, chatId int, botId int) (history *[]*Record, title string, err error)

	// AsyncSaveHistory 异步保存历史记录
	AsyncSaveHistory(ctx context.Context, chatId int, askText string, generationText string, botId int)
	// CacheLuaLruResetHistory 这个是在生成前 把从DB拿到的数据回写缓存 维护热点数据          feat:生成前取消回写 在获取冷历史记录时回写
	CacheLuaLruResetHistory(ctx context.Context, cacheKey string, history *[]*Record, chatId int, title string, botId int) error
	// CacheLuaLruPutHistory 这个是在生成完毕后 回写完整历史记录
	CacheLuaLruPutHistory(ctx context.Context, cacheKey string, history *[]*Record, askText string, generationText string, chatId int, botId int, title string) error

	//由于http.response对象不可序列化 转为inmemory存储
	MemoryGetGeneration(ctx context.Context, chatId int) *GenerationResponse
	CacheGetGeneration(ctx context.Context, chatId int) (*GenerationResponse, error)
	MemoryDelGeneration(ctx context.Context, chatId int)
	CacheDelGeneration(ctx context.Context, chatId int) error

	// CacheGetTitlePrompt 获取根据历史记录获取标题的prompt
	CacheGetTitlePrompt(ctx context.Context) string
	CacheGetTitles(ctx context.Context, userId int, botId int) ([]*TitleData, error)

	DbUpdateTitle(ctx context.Context, chatId int, newTitle string)
	CacheUpdateTitle(ctx context.Context, chatId int, newTitle string, botId int)
}

type ChatUseCase interface {
	InitChat(ctx context.Context, token string, botId int) int
	ContextChat(ctx context.Context, token string, botId int, chatId int, askMessage string, adjustment bool) (isSuccess bool, message ParsedResponse, code int)

	// StreamContextChatSetup 流式输出启动
	StreamContextChatSetup(ctx context.Context, token string, botId int, chatId int, askMessage string, adjustment bool) (isSuccess bool, message ParsedResponse, code int)
	// StreamContextChatWorker 流式输出 信息下发
	StreamContextChatWorker(ctx context.Context, token string, s *sse.Stream)
	// StreamContextStorage 流式输出缓存
	StreamContextStorage(ctx context.Context, token string) bool

	DisposableVisionChat(ctx context.Context, token string, chatId int, botId int, askMessage string, picUrl string) (isSuccess bool, message ParsedResponse, code int)

	//TODO 同上适应前端接口
	//InitMainPage(ctx context.Context, token string) (titles []string, err error)
	InitMainPage(ctx context.Context, token string, botId int) (titles []*TitleData, err error)
	// GetChatHistory 获取历史记录 若从DB取的记录 则回写缓存
	GetChatHistory(ctx context.Context, chatId int, botId int, tokenString string) (*[]*Record, error)

	GenerateUpdateTitle(ctx context.Context, message *[]TextMessage, token string, chatId int) (string, error)
	InputUpdateTitle(ctx context.Context, title string, token string, chatId int, botId int) bool
}

type StorageEvent interface {
	//对repository中方法进行二次封装
	DbPutHistory(b []byte) error
	CachePutHistory(b []byte) error
	DbNewChat(b []byte) error
	DbUpdateTitle(b []byte) error

	PublishSaveDbHistory(data *AskContextData)
	PublishSaveCacheHistory(data *AskContextData)
	PublishDbNewChat(data *ChatStorageData)
	PublishDbSaveTitle(data *AskContextData)

	AsyncConsumeDbHistory()
	AsyncConsumeCacheHistory()
	AsyncConsumeDbNewChat()
	AsyncConsumeDbUpdateTitle()
}

type GenerateEvent interface {
	// StreamDataReady 提前缓存流请求需保存信息 成功则换成不成功则删除
	StreamDataReady(b []byte) error
	PublishStreamReadyStorageData(data *StreamGenerationReadyStorageData)
	AsyncStreamStorageDataReady()
}

type ChatStorageData struct {
	UserId int
	ChatId int
	BotId  int
}

// TODO 适应前端接口
type TitleData struct {
	Title  string `json:"title"`
	ChatId int    `json:"chat_id"`
}
