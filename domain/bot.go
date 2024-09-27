package domain

import "context"

// Bot struct
type Bot struct {
	*BotInfo   `json:"bot_info"`
	*BotConfig `json:"bot_config"`
	BotId      int `gorm:"primaryKey" json:"bot_id"`
	//是否已经删除
	IsDelete bool `json:"is_delete"`
	//是否官方bot
	IsOfficial bool `json:"is_official"`
}

type BotInfo struct {
	BotId       int    `json:"bot_id"`
	Name        string `json:"bot_name"`
	Avatar      string `json:"bot_avatar"`
	Description string `json:"bot_description"`
}

type BotConfig struct {
	BotId            int    `json:"bot_id"`
	InitPrompt       string `json:"init_prompt"`
	Model            string `json:"model"`
	AdjustmentPrompt string `json:"adjustment_prompt"`
}

type BotRepository interface {
	CacheGetBotConfig(ctx context.Context, botId int) *BotConfig
	CacheGetMaxBotId(ctx context.Context) int
}

type BotUseCase interface {
}
