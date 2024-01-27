package models

import (
	"mini-gpt/dao"
)

// Chat 一个Chat
type Chat struct {
	ChatId         int       `json:"chat_id"  gorm:"primaryKey"`
	UserId         int       `json:"user_id"`
	Title          int       `json:"title"`
	LastUpdateTime int64     `json:"last_update_time"`
	IsDelete       bool      `json:"is_delete"`
	Records        *[]Record `json:"records"`
}

// Record 一次问答
type Record struct {
	RecordId    int `json:"record_id"`
	asks        *Ask
	generations *Generation
}

// Ask 一次问题
type Ask struct {
	RecordId int    `json:"record_id"`
	ChatId   int    `json:"chat_id"`
	Message  string `json:"message"`
}

// Generation 一次生成
type Generation struct {
	RecordId int `json:"record_id"`
	ChatId   int `json:"chat_id"`
	Message  int `json:"message"`
}

// ShowChatTitle 主页面展示已有chat的标题
func ShowChatTitle(userId int) ([]*Chat, error) {
	var chats []*Chat
	dao.DB.Table("chat").Where("is_delete = ?", 0).Where("user_id = ?", userId).Find(&chats)
	return chats, nil
}
