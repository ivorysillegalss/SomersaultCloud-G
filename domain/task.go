package domain

type AskContextData struct {
	ChatId         int
	UserId         int
	Message        string
	BotId          int
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
