package domain

import (
	"SomersaultCloud/api/middleware/taskchain"
	"net/http"
	"time"
)

// 目前思路：所有模型的所有任一请求方式 都需要实现下方languagemodel
// 但是同一个厂商的不同请求方式 eg 文生文 or 文生图 可以将他的请求参数综合
// TODO 处理返回的数据格式
type LanguageModelExecutor interface {
	AssemblePrompt(tc *taskchain.TaskContextData) *[]Message
	EncodeReq(tc *taskchain.TaskContextData) *http.Request
	// ConfigureProxy 非必实现 根据api是否被墙
	ConfigureProxy(tc *taskchain.TaskContextData) *http.Client
	Execute(tc *taskchain.TaskContextData) LanguageModelResponse
	// ParseResp TODO
	ParseResp(tc *taskchain.TaskContextData) ParsedResponse
}

type ParsedResponse struct {
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
	time    time.Time
	client  *http.Client
	request *http.Request
}

func NewConnection(client *http.Client, r *http.Request) *ConnectionConfig {
	return &ConnectionConfig{
		time:    time.Time{},
		client:  client,
		request: r,
	}
}
