package domain

import (
	"SomersaultCloud/api/middleware/taskchain"
	"SomersaultCloud/constant/common"
)

// 目前思路：所有模型的所有任一请求方式 都需要实现下方languagemodel
// 但是同一个厂商的不同请求方式 eg 文生文 or 文生图 可以将他的请求参数综合
// TODO 处理返回的数据格式
type LanguageModelExecutor interface {
	AssemblePrompt(tc *taskchain.TaskContextData) *[]Message
	EncodeReq(tc *taskchain.TaskContextData) LanguageModelRequest
	Execute(tc *taskchain.TaskContextData) LanguageModelResponse
	// ParseResp TODO
	ParseResp(tc *taskchain.TaskContextData) ParsedResponse
}

func NewLanguageModelExecutor(id int) LanguageModelExecutor {
	return &OpenaiChatLanguageChatModelExecutor{}
}

type OpenaiChatLanguageChatModelExecutor struct {
}

// AssemblePrompt 组装prompt
// TODO 这个架构上或许可以改进 现在每调用一次都需要 转5次消息格式
//
//	一思路是将他转成哈希
func (o OpenaiChatLanguageChatModelExecutor) AssemblePrompt(tc *taskchain.TaskContextData) *[]Message {
	var msgs []Message
	historyChat := *tc.History
	var i float64
	i = 0
	for i < common.HistoryDefaultWeight {
		user := &Message{
			Role:    common.UserRole,
			Content: historyChat[int(i)].ChatAsks.Message,
		}
		msgs = append(msgs, *user)
		asst := &Message{
			Role:    common.GPTRole,
			Content: historyChat[int(i)].ChatGenerations.Message,
		}
		msgs = append(msgs, *asst)
		i++
	}
	last := &Message{
		Role:    common.UserRole,
		Content: tc.Prompt,
	}
	msgs = append(msgs, *last)
	return &msgs
}

func (o OpenaiChatLanguageChatModelExecutor) EncodeReq(tc *taskchain.TaskContextData) LanguageModelRequest {
	return &OpenaiChatLanguageChatModelRequest{
		Message: *tc.HistoryMessage,
		Model:   tc.Model,
	}
}

func (o OpenaiChatLanguageChatModelExecutor) Execute(tc *taskchain.TaskContextData) LanguageModelResponse {
	//TODO implement me
	panic("implement me")
}

func (o OpenaiChatLanguageChatModelExecutor) ParseResp(tc *taskchain.TaskContextData) ParsedResponse {
	//TODO implement me
	panic("implement me")
}

type ParsedResponse struct {
}

type LanguageModelResponse interface {
}

type LanguageModelRequest interface {
	//MaxTokens int    `json:"max_tokens"`
	//Model string `json:"model"`
	req()
}

type OpenaiChatLanguageChatModelRequest struct {
	Message []Message `json:"messages"`
	Model   string    `json:"model"`
	//TODO 此处可丰富详细参数 见openai api doc
}

func (o *OpenaiChatLanguageChatModelRequest) req() {}

//type ChatCompletionRequest struct {
//	LanguageModelRequest
//	Messages []Message `json:"messages"`
//}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
