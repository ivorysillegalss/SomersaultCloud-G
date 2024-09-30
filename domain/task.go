package domain

type AskContextData struct {
	ChatId         int
	UserId         int
	Message        string
	BotId          int
	Adjustment     bool
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
