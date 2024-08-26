package domain

import (
	"net/http"
	"time"
)

// 目前思路：所有模型的所有任一请求方式 都需要实现下方languagemodel
// 但是同一个厂商的不同请求方式 eg 文生文 or 文生图 可以将他的请求参数综合
// TODO 还是希望这里executor能跟domain解耦，看着有点难受，虽然不违法开发规范
type LanguageModelExecutor interface {
	AssemblePrompt(tc *AskContextData) *[]Message
	EncodeReq(tc *AskContextData) *http.Request
	// ConfigureProxy 非必实现 根据api是否被墙
	ConfigureProxy(tc *AskContextData) *http.Client
	Execute(tc *AskContextData)
	ParseResp(tc *AskContextData) (ParsedResponse, string)
}

type ParsedResponse interface {
	GetGenerateText() string
}

type OpenAIParsedResponse struct {
	GenerateText string
	FinishReason string
}

func (o *OpenAIParsedResponse) GetGenerateText() string {
	return o.GenerateText
}

type LanguageModelRequest interface {
	//MaxTokens int    `json:"max_tokens"`
	//Model string `json:"model"`
	Req()
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ConnectionConfig struct {
	Time    time.Time
	Client  *http.Client
	Request *http.Request
}

func NewConnection(client *http.Client, r *http.Request) *ConnectionConfig {
	return &ConnectionConfig{
		Time:    time.Time{},
		Client:  client,
		Request: r,
	}
}
