package service

import (
	"math"
	"mini-gpt/api"
	"mini-gpt/constant"
	"mini-gpt/dto"
	"mini-gpt/models"
	"mini-gpt/setting"
	utils "mini-gpt/utils/jwt"
	"mini-gpt/utils/redisUtils"
	"strconv"
)

//var logger = setting.GetLogger()
//正常来说应该全局变量的 但是由于代码的先后执行问题先放到下面的函数中

// 最初始的调用方式
func LoadingChat(dto *dto.ChatDTO) (*models.GenerateMessage, error) {
	m := models.Message{
		Role:    constant.UserRole,
		Content: dto.InputPrompt,
	}
	var msgs []models.Message
	msgs = append(msgs, m)

	apiRequestMessage := &models.ChatCompletionRequest{
		CompletionRequest: models.CompletionRequest{
			MaxTokens: dto.MaxToken,
			Model:     dto.Model,
		},
		Messages: msgs,
		//TODO 组装消息记录
	}

	var logger = setting.GetLogger()
	completionResponse, err := api.Execute("-1", apiRequestMessage, dto.Model)
	if err != nil {
		logger.Error(err)
		return models.ErrorGeneration(), err
	}

	//modelType := api.ModelsMap.Models[dto.Model]
	//在这里判断它的类型 并根据不同类型的响应信息格式进行处理

	//这里的YAML配置读不进来 TODO
	//先默认是普通大模型
	generationMessage := api.SimplyMessage(completionResponse, constant.DefaultModel)
	return generationMessage, nil
}

// 将botConfig配置包装为调用api所需请求体
func botConfigToApiRequest(config *models.BotConfig) models.ApiRequestMessage {
	//TODO 改不动了 这里先默认是上下文大模型 并且只有一次性能用
	m := models.Message{
		Role:    constant.UserRole,
		Content: "",
	}
	var msgs []models.Message
	_ = append(msgs, m)

	return &models.ChatCompletionRequest{
		CompletionRequest: models.CompletionRequest{
			MaxTokens: constant.DefaultMaxToken,
			Model:     config.Model,
		},
		Messages: msgs,
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
	completionResponse, err := api.Execute(dto.UserId, botRequest, config.Model)
	if err != nil {
		return models.ErrorGeneration(), err
	}
	generationMessage := api.SimplyMessage(completionResponse, constant.DefaultModel)
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

func CreateChat(dto *dto.CreateChatDTO, tokenString string) (botId int, err error) {
	userId, err := utils.DecodeToId(tokenString)
	if err != nil {
		return constant.ZeroInt, err
	}
	return models.CreateNewChat(userId, dto.BotId)
}

// 使用默认的大模型
func defaultContextModel(askDTO *dto.AskDTO) *models.BotConfig {
	return &models.BotConfig{
		BotId:      0,
		InitPrompt: askDTO.Ask.Message,
		Model:      constant.DefaultModel,
		//Model: constant.InstructModel,
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
func ContextChat(ask *dto.AskDTO, tokenString string) (*models.GenerateMessage, error) {

	userId, err2 := utils.DecodeToId(tokenString)
	if err2 != nil {
		return nil, err2
	}
	ask.UserId = userId

	askInfo := ask.Ask
	botId := askInfo.BotId

	var botConfig *models.BotConfig
	//如果是0 则代表使用默认的大模型
	if botId == constant.DefaultContextModel {
		botConfig = defaultContextModel(ask)
	} else {
		botConfig, _ = otherContextModel(ask)
	}

	//从redis缓存 或mysql中获取历史记录
	history, err := models.GetChatHistory4DefaultContext(askInfo.ChatId)
	if err != nil {
		return models.ErrorGeneration(), err
	}

	//往redis中更新缓存
	//err = redisUtils.SetStructWithExpire(constant.ChatCache+strconv.Itoa(askInfo.ChatId), history, constant.ChatCacheExpire)
	if err != nil {
		return models.ErrorGeneration(), err
	}

	//根据已有权重（历史记录）更新上下文提示词
	botRequest := updateContextPrompt(history, botConfig)

	//更新引用功能
	//botConfig.InitPrompt = referenceToken(botConfig.InitPrompt, ask)

	//包装为请求体
	//botRequest := botConfigToApiRequest(botConfig)

	completionResponse, err := api.Execute(strconv.Itoa(ask.UserId), botRequest, botConfig.Model)
	if err != nil {
		return models.ErrorGeneration(), err
	}
	generationMessage := api.SimplyMessage(completionResponse, constant.DefaultModel)

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

	//err = redisUtils.SetStructWithExpire(constant.ChatCache+strconv.Itoa(askInfo.ChatId), history, constant.ChatCacheExpire)
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
	if referenceRecord != new(models.Record) && asks.ReferenceToken != constant.ZeroString {
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

// updateContextPrompt 针对上下文模型 包装信息且加载历史记录
func updateContextPrompt(history *[]*models.Record, botConfig *models.BotConfig) *models.ChatCompletionRequest {
	//这里先直接填充历史记录进入prompt
	var msgs []models.Message

	historyChat := *history
	currentlyHistoryWeight := math.Min(constant.ChatHistoryWeight, float64(len(historyChat)))

	var i float64
	i = 0
	for i < currentlyHistoryWeight {
		user := &models.Message{
			Role:    constant.UserRole,
			Content: historyChat[int(i)].ChatAsks.Message,
		}
		msgs = append(msgs, *user)
		asst := &models.Message{
			Role:    constant.GPTRole,
			Content: historyChat[int(i)].ChatGenerations.Message,
		}
		msgs = append(msgs, *asst)
		i++
	}

	last := &models.Message{
		Role:    constant.UserRole,
		Content: botConfig.InitPrompt,
	}

	msgs = append(msgs, *last)
	return &models.ChatCompletionRequest{
		CompletionRequest: models.CompletionRequest{
			MaxTokens: constant.DefaultMaxToken,
			Model:     botConfig.Model,
		},
		Messages: msgs,
	}
}

// updateTextCompletionPrompt 只对文本补全的引擎有用
func updateTextCompletionPrompt(history *[]*models.Record, prompt string) (initPrompt string) {
	//根据已有权重进行更新提示词 TODO  算法待补充

	//这里先直接填充历史记录进入prompt
	historyChat := *history

	if len(historyChat) == 0 {
		initPrompt = prompt
		return
	}

	initPrompt = constant.HistoryChatPrompt
	//直接将预处理话术和历史记录拼接的做法欠优 可能可以改进

	//初始化聊天记录 告诉gpt以下是我和你的聊天记录
	currentlyHistoryWeight := math.Min(constant.ChatHistoryWeight, float64(len(historyChat)))
	var i float64
	i = 0
	for i < currentlyHistoryWeight {
		initPrompt += constant.UserRole + historyChat[int(i)].ChatAsks.Message + "\n"
		initPrompt += constant.GPTRole + historyChat[int(i)].ChatGenerations.Message + "\n"

		if len(initPrompt) > constant.JumpOutToken {
			break
		}

		i++
	}

	initPrompt += constant.NowAsk + prompt
	return
}

func GetChatHistory(chatId int) (*[]*models.Record, error) {
	return models.GetChatHistory(chatId)
}

// 在进行一次聊天之后 返回对应的chat的标题
func UpdateInitTitle(historyDTO *dto.TitleDTO) (*dto.TitleDTO, error) {
	messages := historyDTO.Messages
	concludedTitle, err := api.ConcludeTitle(&messages)

	err = models.UpdateChatTitle(historyDTO.ChatId, concludedTitle)

	//没问题就换异步 TODO
	//go asyncUpdateTitle(history.ChatId, concludedTitle)

	if err != nil {
		return nil, err
	}
	return &dto.TitleDTO{Title: concludedTitle}, nil
}

// 更新现有的标题
func UpdateCurrentTitle(currentTitleDTO *dto.TitleDTO) error {
	err := models.UpdateChatTitle(currentTitleDTO.ChatId, currentTitleDTO.Title)
	if err != nil {
		return err
	}
	return nil
}

func asyncUpdateTitle(chatId int, concludedTitle string) {
	_ = models.UpdateChatTitle(chatId, concludedTitle)
}
