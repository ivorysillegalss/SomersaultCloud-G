package domain

import (
	"SomersaultCloud/task"
	"net/http"
	"time"
)

// 目前思路：所有模型的所有任一请求方式 都需要实现下方languagemodel
// 但是同一个厂商的不同请求方式 eg 文生文 or 文生图 可以将他的请求参数综合
type LanguageModelExecutor interface {
	AssemblePrompt(tc *task.AskContextData) *[]Message
	EncodeReq(tc *task.AskContextData) *http.Request
	// ConfigureProxy 非必实现 根据api是否被墙
	ConfigureProxy(tc *task.AskContextData) *http.Client
	Execute(tc *task.AskContextData)
	ParseResp(tc *task.AskContextData) ParsedResponse
}

type ParsedResponse interface {
	parse()
}
type OpenAIParsedResponse struct {
	GenerateText string
	FinishReason string
}

func (o *OpenAIParsedResponse) parse() {
}

type LanguageModelResponse interface {
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
