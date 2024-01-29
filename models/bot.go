package models

import "mini-gpt/dao"

// TalkBot 语言模型接口
type TalkBot interface {
	// 更新名字 模型 描述
	//UpdateName()
	//UpdateModel()
	//UpdateDescription()
	TextToGenerate()
}

// 自定义 bot的结构体
type Bot struct {
	*BotInfo
	*BotConfig
	BotId int `gorm:"primaryKey"`
	//是否已经删除
	IsDelete bool
	//是否官方bot
	IsOfficial bool
}

type BotInfo struct {
	BotId       int    `json:"bot_id"`
	Name        string `json:"bot_name"`
	Avatar      string `json:"bot_avatar"`
	Description string `json:"bot_description"`
}

type BotConfig struct {
	BotId      int `json:"bot_id"`
	InitPrompt string
	Model      string
}

func CreateBotInfo(botInfo BotInfo) error {
	if err := dao.DB.Create(botInfo).Error; err != nil {
		return err
	}
	return nil
}

func CreateBotConfig(config BotConfig) error {
	if err := dao.DB.Create(config).Error; err != nil {
		return err
	}
	return nil
}

// 获取特定bot的信息
func GetBotConfig(botId int) (BotConfig, error) {
	var botConfig BotConfig
	err := dao.DB.Table("bot_config").Where("bot_id = ?", botId).Find(botConfig).Error
	return botConfig, err
}
