package task

const (
	SuccessCode    = 0
	SuccessMessage = "操作成功"

	FailCode    = -1
	FailMessage = "操作失败"

	ExecutingCode    = 1
	ExecutingMessage = "责任链执行中"

	NoneNode = "没节点怎么执行老铁"

	InvalidDataFormatMessage = "消息格式错误"
	InvalidDataMarshal       = "数据转码错误"
	HistoryRetrievalFailed   = "历史记录调取失败"
	BotRetrievalFailed       = "机器人相关信息调取失败"
	ReqDataMarshalFailed     = "请求序列化失败"
	ReqUploadError           = "请求上传失败"
	ReqParsedError           = "请求转码失败"
	ReqCatchError            = "请求获取失败"
	ChatGenerationError      = "生成内容请求事变"
	ChatGenerationDelError   = "删除生成缓存失败"
	RespParedError           = "生成内容转码失败"

	ExecuteChatAskType = "ExecuteChatAsk"
	ExecuteChatAskCode = 10
	ChatAskExecutorId  = 1

	ExecuteChatVisionAskType = "ExecuteChatVisionAsk"
	ExecuteChatVisionAskCode = 20
	ChatVisionAskExecutorId  = 2

	ExecuteTitleAskType    = "ExecuteTitleAsk"
	ExecuteTitleAskCode    = 30
	ChatTitleAskExecutorId = 3

	ExecuteChatAskTypeDs      = "ExecuteChatAskDs"
	ExecuteChatAskCodeDs      = 40
	DeepSeekChatAskExecutorId = 4

	StorageStreamType = "StorageStreamCode"
	StorageStreamCode = 40
)
