package service

import (
	"mini-gpt/constant"
	"mini-gpt/dto"
	"mini-gpt/models"
	"mini-gpt/utils/redisUtils"
)

// 管理员获取机器人信息
func GetBot(botId int, isOfficial int) (*models.Bot, error) {
	var bot *models.Bot
	var err error
	if isOfficial == 0 {
		bot, err = models.GetUnofficialBot(botId)
	} else {
		bot, err = models.GetOfficialBot(botId)
	}

	bot.BotId = botId
	if err != nil {
		return models.ErrorBot(), err
	}
	return bot, nil
}

func adminCreateNewBot(botId int, dto dto.CreateBotDTO) *models.Bot {
	return &models.Bot{
		BotId:      botId,
		BotInfo:    dto.BotInfo,
		BotConfig:  dto.BotConfig,
		IsDelete:   false,
		IsOfficial: true,
	}
}

// 管理员创建新机器人
func AdminCreateBot(dto dto.CreateBotDTO) error {
	botId, err := models.CreateBot(true)
	newBot := adminCreateNewBot(botId, dto)
	if err != nil {
		return err
	}
	//存入redis当中
	return redisUtils.SetStruct(constant.OfficialBotPrefix, newBot)
}
