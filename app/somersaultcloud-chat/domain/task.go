package domain

type AskContextData struct {
	ChatId         int
	UserId         int
	Message        string
	BotId          int
	Adjustment     bool
	Stream         bool //是否开启流式输出
	History        *[]*Record
	SysPrompt      string
	Model          string
	HistoryMessage *Message
	ExecutorId     int
	Executor       LanguageModelExecutor
	Conn           ConnectionConfig
	Resp           GenerationResponse
	ParsedResponse ParsedResponse
}

func (a *AskContextData) Data() {}
