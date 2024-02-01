package models

import (
	"mini-gpt/constant"
	"mini-gpt/dao"
	"mini-gpt/utils/redisUtils"
)

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

// 获取特定bot的配置信息
func GetBotConfig(botId int) (BotConfig, error) {
	var botConfig BotConfig
	err := dao.DB.Table("bot_config").Where("bot_id = ?", botId).Find(&botConfig).Error
	return botConfig, err
}

// 获取特定bot一般信息
func GetBotInfo(botId int) (BotInfo, error) {
	var botInfo BotInfo
	err := dao.DB.Table("bot_info").Where("bot_id = ?", botId).Find(&botInfo).Error
	return botInfo, err
}

// 获取bot启动信息
func GetUnofficialBot(botId int) (Bot, error) {
	var bot Bot
	var err error

	err = dao.DB.Table("bot").Where("bot_id = ?", botId).Find(&bot).Error

	botInfo, err1 := GetBotInfo(botId)
	botConfig, err2 := GetBotConfig(botId)
	//有点想要创建一个错误管道工具处理批量错误。。。
	if err != nil || err2 != nil || err1 != nil {
		//记录日志
		return ErrorBot(), err
	}
	bot.BotInfo = &botInfo
	bot.BotConfig = &botConfig
	return bot, nil
}

func ErrorBot() Bot {
	return Bot{
		BotInfo:    nil,
		BotConfig:  nil,
		BotId:      0,
		IsDelete:   false,
		IsOfficial: false,
	}
}

// 将官方机器人存到redis当中 如果调用的是官方的 直接从redis中取出
func GetOfficialBot(botId int) (Bot, error) {
	botIdStr := string(rune(botId))
	k := constant.OfficialBotPrefix + botIdStr
	resBot, err := redisUtils.GetStruct[Bot](k)
	return resBot, err
}
