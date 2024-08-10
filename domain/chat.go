package domain

import "context"

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
	CacheGetNewestChatId(ctx context.Context) int
	CacheInsertNewChat(ctx context.Context, id int)
	CacheLuaInsertNewChatId(ctx context.Context, luaScript string, k string) (int, error)
	DbInsertNewChatId(ctx context.Context, token int, id int)
}

type ChatUseCase interface {
	InitChat(ctx context.Context, token string, botId int) int
}
