package constant

import "time"

const (
	// DefaultMaxToken 不输入token时默认定义的token最大数量
	DefaultMaxToken = 1000
	// DefaultMaxLimitedTime 请求默认超时时间
	DefaultMaxLimitedTime = time.Minute / 2
	// ApiServerOpenAI OpenAI-API服务器默认网址
	ApiServerOpenAI = "https://api.openai.com/v1/completions"
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
)
