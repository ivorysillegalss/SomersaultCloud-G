package dto

// InitDTO 初始化映射
type InitDTO struct {
	UserId int `json:"user_id"`
	ChatId int `json:"chat_id"`
}

// BotDTO 调用ai功能映射类
type BotDTO struct {
	BotId int `json:"bot_id"`
	//提示词配置
	Configs []string `json:"configs"`
}
