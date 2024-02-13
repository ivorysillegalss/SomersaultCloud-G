package service

import (
	"mini-gpt/constant"
	"mini-gpt/dto"
	"mini-gpt/models"
	"mini-gpt/utils/redisUtils"
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

// 管理员更新机器人 初见反射
func AdminUpdateBot(updatedBotDTO *dto.UpdateBotDTO) error {
	updatedBot := updatedBotDTO.Bot

	if updatedBot.IsOfficial {
		beforeBot, err := models.GetOfficialBot(updatedBot.BotId)
		if err != nil {
			return err
		}

		//这里的反射代码其实可以用第三方structs类进行简化
		//通过反射部分更新目标结构体
		updateValue := reflect.ValueOf(updatedBot)
		//这个updateValue是需要部分更新的数据的元数据集 自己理解
		dataValue := reflect.ValueOf(&beforeBot).Elem()
		//Elem方法是获取到结构体字段或数组切片等数据结构底层字段值的方法
		//于是这里先通过结构体的指针获取 其中指向的字段的值

		for i := 0; i < updateValue.NumField(); i++ {
			field := updateValue.Type().Field(i)
			//类似java中的Field类型对象
			updateFieldValue := updateValue.Field(i)
			if updateFieldValue.IsZero() {
				continue
			}
			//如果是零值则跳过 不更新
			dataFieldValue := dataValue.FieldByName(field.Name)
			//通过属性的key获取返回的值
			if dataFieldValue.IsValid() && dataFieldValue.CanSet() {
				dataFieldValue.Set(updateFieldValue)
			}
		}
		//整个过程的思想是将传进来的updatedBot中需要部分更新的值通过反射 更新入旧的bot中
		//再将整个bot重新set一遍
		err = redisUtils.SetStruct(constant.OfficialBotPrefix+string(rune(beforeBot.BotId)), beforeBot)
		if err != nil {
			return err
		}
	} else {
		err := models.UpdateUnofficialBot(updatedBot)
		if err != nil {
			return err
		}
	}
	return nil
}
