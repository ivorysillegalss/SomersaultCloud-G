package service

import (
	"mini-gpt/constant"
	"mini-gpt/dto"
	"mini-gpt/models"
	"mini-gpt/utils/redisUtils"
	"mini-gpt/utils/reflectUtils"
	"reflect"
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
	return redisUtils.SetStruct(constant.OfficialBotPrefix+string(rune(botId)), newBot)
}

// 将映射Map转换为对应model 很丑陋很丑陋 因为反射实现更新的d不下去才这么写的
func updateMapToModel(updatedBotMap *map[string]interface{}) *models.Bot {
	botMap := *updatedBotMap
	return &models.Bot{
		BotInfo: &models.BotInfo{
			//BotId:       botMap["bot_id"].(int),
			Name:        botMap["bot_name"].(string),
			Avatar:      botMap["bot_avatar"].(string),
			Description: botMap["bot_description"].(string),
		},
		BotConfig: &models.BotConfig{
			//BotId:      botMap["bot_id"].(int),
			InitPrompt: botMap["init_prompt"].(string),
			Model:      botMap["model"].(string),
		},
		//BotId:      botMap["bot_id"].(int),
		IsOfficial: botMap["is_official"].(bool),
	}
}

// 管理员更新机器人
func AdminUpdateBot(updatedBotMap *map[string]interface{}) error {
	m := *updatedBotMap

	isOfficial := m["is_official"].(bool)
	botIdFloat := m["bot_id"].(float64)
	botId := int(botIdFloat)

	if isOfficial {
		beforeBot, err := models.GetOfficialBot(botId)
		//beforeBot为*Bot
		if err != nil {
			return err
		}

		// 使用反射更新用户结构体
		reflectUtils.UpdateStruct(reflect.ValueOf(beforeBot), *updatedBotMap)

		err = redisUtils.SetStruct(constant.OfficialBotPrefix+string(rune(beforeBot.BotId)), beforeBot)
		if err != nil {
			return err
		}
	} else {
		//将需要更新的字段映射成为Map传进来 转成model
		updatedBot := updateMapToModel(updatedBotMap)
		err := models.UpdateUnofficialBot(updatedBot)
		if err != nil {
			return err
		}
	}
	return nil
}
