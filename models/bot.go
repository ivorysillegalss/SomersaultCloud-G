package models

import (
	"mini-gpt/constant"
	"mini-gpt/dao"
	"mini-gpt/utils/redisUtils"
	"strconv"
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
	BotId      int    `json:"bot_id"`
	InitPrompt string `json:"init_prompt"`
	Model      string `json:"model"`
}

// bot相关操作映射到sql中的结构体 害怕Bot直接调用会有错误 专门整了一个内部使用
type BotToStruct struct {
	ID         int  `gorm:"primaryKey column:bot_id"` // 明确指定BotId为主键
	IsDelete   bool `gorm:"column:is_delete"`         // 导出字段，并可指定列名
	IsOfficial bool `gorm:"column:is_official"`       // 导出字段，并可指定列名
}

// 写入映射结构体对象中
func writeBotToStruct(isOfficial bool) *BotToStruct {
	return &BotToStruct{
		IsDelete:   false,
		IsOfficial: isOfficial,
	}
}

func CreateBot(isOfficial bool) (int, error) {
	botToStruct := writeBotToStruct(isOfficial)
	if err := dao.DB.Table("bot").Create(botToStruct).Error; err != nil {
		return -1, err
	}
	return botToStruct.ID, nil
}

func CreateBotInfo(botInfo *BotInfo) error {
	if err := dao.DB.Table("bot_info").Create(botInfo).Error; err != nil {
		return err
	}
	return nil
}

func CreateBotConfig(config *BotConfig) error {
	if err := dao.DB.Table("bot_config").Create(config).Error; err != nil {
		return err
	}
	return nil
}

// 获取特定bot的配置信息
func GetBotConfig(botId int) (*BotConfig, error) {
	var botConfig BotConfig
	err := dao.DB.Table("bot_config").Where("bot_id = ?", botId).Find(&botConfig).Error
	return &botConfig, err
}

// 获取特定bot一般信息
func GetBotInfo(botId int) (*BotInfo, error) {
	var botInfo BotInfo
	err := dao.DB.Table("bot_info").Where("bot_id = ?", botId).Find(&botInfo).Error
	return &botInfo, err
}

// 获取bot启动信息
func GetUnofficialBot(botId int) (*Bot, error) {
	var bot Bot
	var err error

	err = dao.DB.Table("bot").Where("ID = ?", botId).Find(&bot).Error

	botInfo, err1 := GetBotInfo(botId)
	botConfig, err2 := GetBotConfig(botId)
	//有点想要创建一个错误管道工具处理批量错误。。。
	if err != nil || err2 != nil || err1 != nil {
		//记录日志
		return ErrorBot(), err
	}
	bot.BotInfo = botInfo
	bot.BotConfig = botConfig
	return &bot, nil
}

func ErrorBot() *Bot {
	return &Bot{
		BotInfo:    nil,
		BotConfig:  nil,
		BotId:      0,
		IsDelete:   false,
		IsOfficial: false,
	}
}

// redis映射类
type BotToRedis struct {
	BotName        string `json:"bot_name"`
	BotAvatar      string `json:"bot_avatar"`
	BotDescription string `json:"bot_description"`
	InitPrompt     string `json:"init_prompt"`
	Model          string `json:"model"`
	BotId          int    `json:"bot_id"`
	IsDelete       bool   `json:"is_delete"`
	IsOfficial     bool   `json:"is_official"`
}

// 将官方机器人存到redis当中 如果调用的是官方的 直接从redis中取出
func GetOfficialBot(botId int) (*Bot, error) {
	botIdStr := strconv.Itoa(botId)
	k := constant.OfficialBotPrefix + botIdStr
	resBot, err := redisUtils.GetStruct[BotToRedis](k)
	resBot.BotId = botId
	resBot.IsOfficial = true
	//不知道为什么getStruct出来bool值是false 由于这里是官方的就直接写为true了
	bot := convertRedisBot(&resBot)
	return bot, err
}

func convertRedisBot(bot *BotToRedis) *Bot {
	return &Bot{
		BotInfo: &BotInfo{
			BotId:       bot.BotId,
			Name:        bot.BotName,
			Avatar:      bot.BotAvatar,
			Description: bot.BotDescription,
		},
		BotConfig: &BotConfig{
			BotId:      bot.BotId,
			InitPrompt: bot.InitPrompt,
			Model:      bot.Model,
		},
		BotId:      bot.BotId,
		IsDelete:   bot.IsDelete,
		IsOfficial: bot.IsOfficial,
	}
}

// 非官方放到mysql的数据
func UpdateUnofficialBot(bot *Bot) error {
	//这里可以根据部分更新需求优化 TBD
	err := dao.DB.Model(&bot).Where("bot_id = ?", bot.BotId).Updates(bot).Error
	return err
}

func redisBotConvert(beforeBot *Bot) *BotToRedis {
	return &BotToRedis{
		BotName:        beforeBot.Name,
		BotAvatar:      beforeBot.Avatar,
		BotDescription: beforeBot.Description,
		InitPrompt:     beforeBot.InitPrompt,
		Model:          beforeBot.Model,
		BotId:          beforeBot.BotId,
		IsDelete:       false,
		IsOfficial:     true,
	}
}

func SetOfficialBot(beforeBot *Bot) error {
	redisBot := redisBotConvert(beforeBot)
	return redisUtils.SetStruct(constant.OfficialBotPrefix+strconv.Itoa(beforeBot.BotId), redisBot)
}

func CreateOfficialBot(bot *Bot) *BotToRedis {
	return redisBotConvert(bot)
}
