package dto

import "mini-gpt/models"

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

// CreateBotDTO 创建机器人映射
type CreateBotDTO struct {
	BotInfo   *models.BotInfo   `json:"bot_info"`
	BotConfig *models.BotConfig `json:"bot_config"`
	BotId     int               `json:"bot_id"`
}

// UpdateBotDTO 这个映射类和models里的bot没区别
// 单独写多一个是因为不想controller和models层有耦合
type UpdateBotDTO struct {
	Bot *models.Bot `json:"bot"`
}
