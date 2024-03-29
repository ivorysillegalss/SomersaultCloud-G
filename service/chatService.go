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
	generationMessage := simplyMessage(completionResponse)
	return generationMessage, nil
}

// 将调用api结果包装为返回用户的结果
func simplyMessage(completionResponse *models.CompletionResponse) *models.GenerateMessage {
	//openAI返回的json中请求体中的文本是一个数组 暂取第0项
	args := completionResponse.Choices
	if args == nil {
		return models.ErrorGeneration()
	}
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

// 获取bot 有无历史记录通用代码
func getBot(botId int) (*models.BotConfig, error) {
	//获取redis中所保存的所有官方机器人的id集合
	list, err2 := redisUtils.GetList(constant.OfficialBotIdList)
	if err2 != nil {
		return models.ErrorBotConfig(), err2
	}
	var config *models.BotConfig
	var err error
	//在list中 由于每一次更新新的机器人都是从rpush的 所以可以直接通过大小进行比较 从左往右即为从小到大
	//如果目标的比所需的大 即表明没有这个官方机器人 这里可以用哈希表优化
	for i := range list {
		eachOfficialBotId, _ := strconv.Atoi(list[i])
		if eachOfficialBotId > botId {
			config, err = models.GetBotConfig(botId)
			if err != nil {
				return models.ErrorBotConfig(), err
			}
			break
		} else if eachOfficialBotId == botId {
			//如果有这个官方机器人 就需要从redis中取它的配置
			bot, err2 := models.GetOfficialBot(eachOfficialBotId)
			config = bot.BotConfig
			if err2 != nil {
				return models.ErrorBotConfig(), err
			}
			break
		}
	}
	config.BotId = botId
	return config, nil
}

// 一次性使用的bot调用方式 (没有历史记录功能的调用方法)
func DisposableChat(dto *dto.ExecuteBotDTO) (*models.GenerateMessage, error) {
	botId := dto.BotId
	botPromptConfigs := dto.Configs
	config, err2 := getBot(botId)
	if err2 != nil {
		return models.ErrorGeneration(), err2
	}
	//修改自定义配置
	config.InitPrompt = updateCustomizeConfig(config.InitPrompt, botPromptConfigs)
	//包装为请求体
	botRequest := botConfigToApiRequest(config)
	completionResponse, err := api.Execute(dto.UserId, botRequest)
	if err != nil {
		return models.ErrorGeneration(), err
	}
	generationMessage := simplyMessage(completionResponse)
	return generationMessage, nil
}

// 替换为自定义prompt
// 此处可优化为Boyer-Moore-Horspool 或优化KMP算法
func updateCustomizeConfig(defaultPrompt string, customize []string) string {
	replaced := ""
	placeholderIndex := 0

	for _, runeValue := range defaultPrompt {
		if runeValue == constant.ReplaceCharFromDefaultToCustomize {
			if placeholderIndex < len(customize) {
				replaced += customize[placeholderIndex]
				placeholderIndex++
			} else {
				// 如果替换内容用尽，保留原始字符
				replaced += string(runeValue)
			}
		} else {
			replaced += string(runeValue)
		}
	}

	return replaced
}

func CreateChat(dto *dto.CreateChatDTO) (botId int, err error) {
	return models.CreateNewChat(dto.UserId, dto.BotId)
}

// 使用默认的大模型
func defaultContextModel(askDTO *dto.AskDTO) *models.BotConfig {
	return &models.BotConfig{
		BotId:      0,
		InitPrompt: askDTO.Ask.Message,
		Model:      constant.DefaultModel,
	}
}

// 使用机器人上下文模型
func otherContextModel(ask *dto.AskDTO) (*models.BotConfig, error) {
	askInfo := ask.Ask
	botId := askInfo.BotId

	botConfig, err := getBot(botId)
	if err != nil {
		return models.ErrorBotConfig(), err
	}

	return botConfig, nil
}

// 具有上下文的chat方式
func ContextChat(ask *dto.AskDTO) (*models.GenerateMessage, error) {
	askInfo := ask.Ask
	botId := askInfo.BotId

	var botConfig *models.BotConfig
	//如果是0 则代表使用默认的大模型
	if botId == 0 {
		botConfig = defaultContextModel(ask)
	} else {
		botConfig, _ = otherContextModel(ask)
	}

	//从redis缓存 或mysql中获取历史记录
	history, err := models.GetChatHistoryForChat(askInfo.ChatId)
	if err != nil {
		return models.ErrorGeneration(), err
	}

	//往redis中更新缓存
	err = redisUtils.SetStructWithExpire(constant.ChatCache+strconv.Itoa(askInfo.ChatId), history, constant.ChatCacheExpire)
	if err != nil {
		return models.ErrorGeneration(), err
	}

	//根据已有权重（历史记录）更新上下文提示词
	botConfig.InitPrompt = updateContextPrompt(history, botConfig.InitPrompt)

	//包装为请求体
	botRequest := botConfigToApiRequest(botConfig)

	completionResponse, err := api.Execute(strconv.Itoa(ask.UserId), botRequest)
	if err != nil {
		return models.ErrorGeneration(), err
	}
	generationMessage := simplyMessage(completionResponse)

	botConfig.InitPrompt = referenceToken(botConfig.InitPrompt, ask)

	//将生成的记录存放入数据库当中
	err = models.SaveRecord(&models.Record{
		ChatAsks: askInfo,
		ChatGenerations: &models.ChatGeneration{
			RecordId: askInfo.RecordId,
			ChatId:   askInfo.ChatId,
			Message:  generationMessage.GenerateText,
		},
	}, askInfo.ChatId)

	if err != nil {
		return models.ErrorGeneration(), err
	}

	//更新权重的时机：新chat初始化 进行了新的问答（新record） 下方为第二种
	//异步更新权重算法TBD
	//go calculateContextWeights(history)

	return generationMessage, nil
}

// 查看是否有引用字段 若有则加入prompt
func referenceToken(beforePrompt string, asks *dto.AskDTO) string {
	referenceRecord := asks.ReferenceRecord
	if !(referenceRecord != new(models.Record) && asks.ReferenceToken != constant.ZeroString) {
		beforePrompt += constant.ReferenceRecordPrompt
		beforePrompt += constant.UserRole + referenceRecord.ChatAsks.Message + "\n"
		beforePrompt += constant.GPTRole + referenceRecord.ChatGenerations.Message + "\n"
	}
	return beforePrompt
}

// 新思路：取数据的时候手动分配权重或更新权重
//分配的时间复杂度On 根据权重进行上下文处理需要遍历 时间复杂度也为On 两者时间复杂度很高 效率相当低

// 也许可以在用户调用的时候 在将历史记录存入数据库的时候 异步分配 更新权重
// 总结就是 异步对权重进行分布处理
// 把权重存进表里存储也许更合适 专门的一张表
func calculateContextWeights(history *[]*models.Record) {
	//需要根据已有记录的数量进行动态的内存分配
}

func updateContextPrompt(history *[]*models.Record, prompt string) (initPrompt string) {
	//根据已有权重进行更新提示词 TODO  算法待补充

	//这里先直接填充历史记录进入prompt
	historyChat := *history
	initPrompt = prompt
	if len(historyChat) != 0 {
		initPrompt = constant.HistoryChatPrompt
		//直接将预处理话术和历史记录拼接的做法欠优 可能可以改进

		//初始化聊天记录 告诉gpt以下是我和你的聊天记录
		i := 0
		for i < constant.ChatHistoryWeight {
			initPrompt += constant.UserRole + historyChat[i].ChatAsks.Message + "\n"
			initPrompt += constant.GPTRole + historyChat[i].ChatGenerations.Message + "\n"
			i++
		}

	} else if len(historyChat) == 0 {
		//这里可以分为两种情况 第一次聊天的时候就不用以上处理
		//这里的做法很粗糙 一旦这个机器人的功能具有上下文 并且需要预处理
		//需要在原始的initPrompt的基础上进行改进
		return
	}

	return
}

func GetChatHistory(chatId int) ([]*models.Record, error) {
	return models.GetChatHistory(chatId)
}
