package domain

type AskContextData struct {
	ChatId         int
	UserId         int
	Message        string
	BotId          int
	History        *[]*Record
	Prompt         string
	Model          string
	HistoryMessage *[]Message
	Executor       LanguageModelExecutor
	Conn           ConnectionConfig
	Resp           GenerationResponse
	ParsedResponse ParsedResponse
}

func (a *AskContextData) Data() {
}
