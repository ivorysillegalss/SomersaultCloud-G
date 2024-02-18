package service

import (
	"mini-gpt/api"
	"mini-gpt/constant"
	"mini-gpt/dto"
	"mini-gpt/models"
	"mini-gpt/setting"
	"mini-gpt/utils/redisUtils"
	"strconv"
)

//var logger = setting.GetLogger()
//正常来说应该全局变量的 但是由于代码的先后执行问题先放到下面的函数中

// 最初始的调用方式
func LoadingChat(apiRequestMessage *models.ApiRequestMessage) (*models.GenerateMessage, error) {

	var logger = setting.GetLogger()

	completionResponse, err := api.Execute("-1", apiRequestMessage)
	if err != nil {
		logger.Error(err)
		return models.ErrorGeneration(), err
	}
	generationMessage := completionResponseToGenerationMessage(completionResponse)
	return generationMessage, nil
}

// 将调用api结果包装为返回用户的结果
func completionResponseToGenerationMessage(completionResponse *models.CompletionResponse) *models.GenerateMessage {
	//openAI返回的json中请求体中的文本是一个数组 暂取第0项
	args := completionResponse.Choices
	textBody := args[0]
	generateMessage := models.GenerateMessage{
		GenerateText: textBody.Text,
		FinishReason: textBody.FinishReason,
	}
	return &generateMessage
}

// 将botConfig配置包装为调用api所需请求体
func botConfigToApiRequest(config *models.BotConfig) *models.ApiRequestMessage {
	return &models.ApiRequestMessage{
		InputPrompt: config.InitPrompt,
		Model:       config.Model,
		//暂定最大字符串不能修改
		MaxToken: constant.DefaultMaxToken,
	}
}

// 一次性使用的bot调用方式 (没有历史记录功能的调用方法)
func DisposableChat(dto *dto.ExecuteBotDTO) (*models.GenerateMessage, error) {
	botId := dto.BotId
	botPromptConfigs := dto.Configs
	//获取redis中所保存的所有官方机器人的id集合
	list, err2 := redisUtils.GetList(constant.OfficialBotIdList)
	if err2 != nil {
		return models.ErrorGeneration(), err2
	}
	var config *models.BotConfig
	var err error
	//在list中 由于每一次更新新的机器人都是从rpush的 所以可以直接通过大小进行比较 从左往右即为从小到大
	//如果目标的比所需的大 即表明没有这个官方机器人
	for i := range list {
		eachOfficialBotId, _ := strconv.Atoi(list[i])
		if eachOfficialBotId > botId {
			config, err = models.GetBotConfig(botId)
			if err != nil {
				return models.ErrorGeneration(), err
			}
		} else if eachOfficialBotId == botId {
			//如果有这个官方机器人 就需要从redis中取它的配置
			bot, err2 := models.GetOfficialBot(eachOfficialBotId)
			config = bot.BotConfig
			if err2 != nil {
				return models.ErrorGeneration(), err
			}
		}
	}
	//修改自定义配置
	config.InitPrompt = updateCustomizeConfig(config.InitPrompt, botPromptConfigs)
	//包装为请求体
	botRequest := botConfigToApiRequest(config)
	completionResponse, err := api.Execute(dto.UserId, botRequest)
	if err != nil {
		return models.ErrorGeneration(), err
	}
	generationMessage := completionResponseToGenerationMessage(completionResponse)
	return generationMessage, nil
}

// 替换为自定义prompt
// 此处可优化为Boyer-Moore-Horspool 或优化KMP算法
func updateCustomizeConfig(defaultPrompt string, customize []string) string {
	replaced := ""
	placeholderIndex := 0

	for i := 0; i < len(defaultPrompt); i++ {
		if defaultPrompt[i] == constant.ReplaceCharFromDefaultToCustomize {
			if placeholderIndex < len(customize) {
				replaced += customize[placeholderIndex]
				placeholderIndex++
			} else {
				// 如果替换内容用尽，保留原始字符
				replaced += string(defaultPrompt[i])
			}
		} else {
			replaced += string(defaultPrompt[i])
		}
	}

	return replaced
}
