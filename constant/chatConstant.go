package constant

import "time"

const (
	// JumpOutToken 跳出获取历史记录的token数量
	JumpOutToken = 1000
	// DefaultContextModel 默认的上下文大模型ID
	DefaultContextModel = 0
	// DefaultMaxToken 不输入token时默认定义的token最大数量
	DefaultMaxToken = 1000
	//DefaultAdminUID 默认的管理员uid （官方调试所用uid）
	DefaultAdminUID = "0"
	// DefaultMaxLimitedTime 请求默认超时时间 方便调试默认关闭
	DefaultMaxLimitedTime = time.Minute / 2
	// ApiServerOpenAI OpenAI-API服务器默认网址
	//ApiServerOpenAI = "https://api.openai.com/v1/completions"
	ApiServerOpenAI = "https://api.openai.com/v1/chat/completions"
	// InstructModel 初始模型
	InstructModel = "gpt-3.5-turbo-instruct"
	// DefaultModel 默认模型
	DefaultModel = "gpt-3.5-turbo-0125"
	// ReplaceCharFromDefaultToCustomize 自定义唯一标识符 选了个挺少见的 可优化算法
	ReplaceCharFromDefaultToCustomize = '¶'
	// OfficialBotPrefix 创建新机器人的前缀
	OfficialBotPrefix = "OfficialBot"
	// UserCachePrefix 用户chat缓存前缀
	UserCachePrefix = "UserCache"
	// OfficialBotIdList redis中存储官方机器人id 维护的便于id查找的list
	OfficialBotIdList = "OfficialBotIdList"
	// ChatCache redis中存储以往chat记录的缓存前缀
	ChatCache = "ChatCache"
	// ChatCacheExpire redis中存储chat记录的限时
	ChatCacheExpire = 30 * time.Minute
	// HistoryChatPrompt 告诉chatGPT以往聊天记录的prompt模板 可改进
	HistoryChatPrompt = "Here is the chat history which I have talked with you,please according to the history give me generation:"
	// SystemRole 系统角色
	SystemRole = "system"
	// UserRole 用户角色
	UserRole = "user"
	// GPTRole GPT角色
	GPTRole = "assistant"
	// NowAsk 当前的一次询问
	NowAsk = "And Here is my question this time:  "
	// Conclude2TitlePrompt 根据一次的历史记录总结出一个标题的提示词
	Conclude2TitlePrompt = "#Content#\n你是一个标题总结员，你总能很完美且精炼的将一段话的内容总结成一个标题。\n" +
		"#Objective# \n现在给你一段对话，请你将对话的内容总结成一个标题,标题要求能够让人知道这段对话的大致内容。直接输出一个标题，不用输出其他的内容。" +
		"\n#Style# \n言简意赅\n#Tone# \n正式\n#input# \n<一段对话>"
	// ChatHistoryWeight 发送上下文历史记录的权重设置
	ChatHistoryWeight = 3
	// APIExecuteSuccessStatus 执行API成功后返回的状态码
	APIExecuteSuccessStatus = 200
	// ReferenceRecordPrompt 告诉chatGPT要他回复回应中的某个部分
	ReferenceRecordPrompt = "Here is a record we have been talked,And I have confused about parts of your generation,please fairly and clearly explain about it and my question:"
	// DefaultShareSecretKeyDestroyTime 默认分享密钥存活时间
	DefaultShareSecretKeyDestroyTime = 24 * time.Hour * 3
)
