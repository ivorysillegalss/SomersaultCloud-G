package dto

import "mini-gpt/models"

// InitDTO 初始化映射
type InitDTO struct {
	UserId int `json:"user_id"`
	ChatId int `json:"chat_id"`
}

// ExecuteBotDTO 调用ai功能映射
type ExecuteBotDTO struct {
	BotId int `json:"bot_id"`
	//提示词配置
	Configs []string `json:"configs"`
}

// CreateBotDTO 创建机器人映射
type CreateBotDTO struct {
	BotInfo   *models.BotInfo
	BotConfig *models.BotConfig
	BotId     int
}
