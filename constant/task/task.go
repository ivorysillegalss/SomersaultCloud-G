package task

const (
	SuccessCode    = 0
	SuccessMessage = "操作成功"

	FailCode    = -1
	FailMessage = "操作失败"

	ExecutingCode    = 1
	ExecutingMessage = "责任链执行中"

	InvalidDataFormatMessage = "消息格式错误"
	InvalidDataMarshal       = "数据转码错误"
	HistoryRetrievalFailed   = "历史记录调取失败"
	BotRetrievalFailed       = "机器人相关信息调取失败"
	ReqDataMarshalFailed     = "请求序列化失败"
)
