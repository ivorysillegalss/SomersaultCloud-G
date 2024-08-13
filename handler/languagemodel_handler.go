package handler

import (
	"SomersaultCloud/api/middleware/taskchain"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/domain"
	"SomersaultCloud/internal/requtil"
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
)

func NewLanguageModelExecutor(id int) domain.LanguageModelExecutor {
	return &OpenaiChatLanguageChatModelExecutor{}
}

func newOpenaiChatLanguageChatModelRequest(message *[]domain.Message, model string) *openaiChatLanguageChatModelRequest {
	return &openaiChatLanguageChatModelRequest{
		Message: *message,
		Model:   model,
	}
}

type openaiChatLanguageChatModelRequest struct {
	Message []domain.Message `json:"messages"`
	Model   string           `json:"model"`
	//TODO 此处可丰富详细参数 见openai api doc
}

func (o *openaiChatLanguageChatModelRequest) Req() {}

type OpenaiChatLanguageChatModelExecutor struct {
}

// AssemblePrompt 组装prompt
// TODO 这个架构上或许可以改进 现在每调用一次都需要 转5次消息格式
//
//	一思路是将他转成哈希
func (o OpenaiChatLanguageChatModelExecutor) AssemblePrompt(tc *taskchain.TaskContextData) *[]domain.Message {
	var msgs []domain.Message
	historyChat := *tc.History
	var i float64
	i = 0
	for i < common.HistoryDefaultWeight {
		user := &domain.Message{
			Role:    common.UserRole,
			Content: historyChat[int(i)].ChatAsks.Message,
		}
		msgs = append(msgs, *user)
		asst := &domain.Message{
			Role:    common.GPTRole,
			Content: historyChat[int(i)].ChatGenerations.Message,
		}
		msgs = append(msgs, *asst)
		i++
	}
	last := &domain.Message{
		Role:    common.UserRole,
		Content: tc.Prompt,
	}
	msgs = append(msgs, *last)
	return &msgs
}

func (o OpenaiChatLanguageChatModelExecutor) EncodeReq(tc *taskchain.TaskContextData) *http.Request {

	jsonData, err := json.Marshal(newOpenaiChatLanguageChatModelRequest(tc.HistoryMessage, tc.Model))
	if err != nil {
		return nil
	}

	request, err := http.NewRequest(http.MethodPost, common.ApiServerOpenAI, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil
	}

	key := fmt.Sprintf("Bearer %s", env.ApiOpenaiSecretKey)
	request.Header.Set("Authorization", key) // 请确保使用你自己的API密钥
	request.Header.Set("Content-Type", "application/json")
	return request
}

func (o OpenaiChatLanguageChatModelExecutor) ConfigureProxy(tc *taskchain.TaskContextData) *http.Client {
	return requtil.SetProxy()
}

func (o OpenaiChatLanguageChatModelExecutor) Execute(tc *taskchain.TaskContextData) domain.LanguageModelResponse {
	//TODO implement me
	panic("implement me")
}

func (o OpenaiChatLanguageChatModelExecutor) ParseResp(tc *taskchain.TaskContextData) domain.ParsedResponse {
	//TODO implement me
	panic("implement me")
}
