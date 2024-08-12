package domain

import (
	"SomersaultCloud/api/middleware/taskchain"
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
}

type Record struct {
	RecordId        int             `json:"record_id"`
	ChatAsks        *ChatAsk        `json:"chat_asks"`
	ChatGenerations *ChatGeneration `json:"chat_generations"`
	//Weights         float64
}

// ChatAsk 一次问题
type ChatAsk struct {
	RecordId int    `json:"record_id"`
	ChatId   int    `json:"chat_id"`
	Message  string `json:"message"`
	BotId    int    `json:"bot_id" gorm:"-"`
	Time     int64  `json:"time"`
}

// ChatGeneration 一次生成
type ChatGeneration struct {
	RecordId int    `json:"record_id"`
	ChatId   int    `json:"chat_id"`
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
	// DbInsertNewChatId 异步使用 存入SQL持久化方法
	DbInsertNewChatId(ctx context.Context, token int, id int)

	// CacheGetHistory 从缓存中取出历史记录 存的时候确保最大条数 取时无需注意
	CacheGetHistory(ctx context.Context, chatId int) (*[]*Record, bool, error)
	// DbGetHistory miss缓存 从DB中获取历史记录
	DbGetHistory(ctx context.Context, chatId int) (*[]*Record, error)

	CacheLuaLruPutHistory(ctx context.Context, k string, v string) error
}

type ChatUseCase interface {
	InitChat(ctx context.Context, token string, botId int) int
}

type ChatTask interface {
	// PreCheckDataTask 数据的前置检查 & 组装TaskContextData对象
	PreCheckDataTask(tc *taskchain.TaskContext)
	// GetHistoryTask 从DB or Cache获取历史记录
	GetHistoryTask(tc taskchain.TaskContext)
	// GetBotTask 获取prompt & model
	GetBotTask(tc taskchain.TaskContext)
	// TODO 微调 TBD
	AdjustmentTask(tc taskchain.TaskContext)
	// AssembleReqTask 组装rpc请求体
	AssembleReqTask(tc *taskchain.TaskContext)
	// CallApiTask 调用api
	CallApiTask(tc *taskchain.TaskContext)
	// ParseRespTask 转换rpc后响应数据
	ParseRespTask(tc *taskchain.TaskContext)
}
