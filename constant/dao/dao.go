package dao

const (
	// DefaultTitle 默认标题
	DefaultTitle = "init"
	// DefaultData 默认数据
	DefaultData = "initData"
	// AsyncPollingFrequency 轮询打日志频率
	AsyncPollingFrequency = 100

	// RecordNotFoundError *old存法没找到历史记录
	RecordNotFoundError = "record not found"

	// DefaultModel 默认模型
	DefaultModel = "gpt-4o-mini"
	// DefaultModelBotId 默认botId
	DefaultModelBotId = 0

	// OriginTable 原始表
	OriginTable = "chat"
	// RefactorTable 重构后的表
	RefactorTable = "chat_re"

	/* -------其他模块的ModelId 见Java模块 -------- */
	CommentBotId   = 2
	MathSolveBotId = 10
)
