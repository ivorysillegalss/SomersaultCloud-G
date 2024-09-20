package dto

import "SomersaultCloud/domain"

// InitDTO 初始化映射
type InitDTO struct {
	UserId int `json:"user_id"`
	ChatId int `json:"chat_id"`
}

// ExecuteBotDTO 调用ai功能映射
type ExecuteBotDTO struct {
	UserId string `json:"user_id"`
	BotId  int    `json:"bot_id"`
	//提示词配置
	Configs []string `json:"configs"`
}

// AskDTO 这个映射类是问题的映射类
type AskDTO struct {
	ChatId     int             `json:"chat_id"`
	BotId      int             `json:"bot_id"`
	Ask        *domain.ChatAsk `json:"ask"`
	UserId     int             `json:"user_id"`
	Adjustment bool            `json:"adjustment"`
	////下方的 Reference是指引用
	//ReferenceToken  string         `json:"reference_token"`
	//ReferenceRecord *domain.Record `json:"reference_record"`
}

// VisionDTO 输入图片时的结构
type VisionDTO struct {
	ChatId  int    `json:"chat_id"`
	BotId   int    `json:"bot_id"`
	Message string `json:"message"`
	PicUrl  string `json:"pic_url"`
}

// CreateChatDTO 创建新chat时候的初始化机器人配置
type CreateChatDTO struct {
	UserId int `json:"user_id"`
	BotId  int `json:"bot_id"`
}

type ChatDTO struct {
	InputPrompt string `json:"prompt"`
	Model       string `json:"model"`
	MaxToken    int    `json:"max_token"`
}

type ShareDTO struct {
	CloneChatId int `json:"clone_chat_id"`
}

type TitleDTO struct {
	Messages []domain.TextMessage `json:"messages"`
	Title    string               `json:"title"`
	ChatId   int                  `json:"chat_id"`
}
