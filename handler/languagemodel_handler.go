package handler

import (
	"SomersaultCloud/api/middleware/taskchain"
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/sys"
	"SomersaultCloud/domain"
	"SomersaultCloud/internal/requtil"
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
)

func NewLanguageModelExecutor(env *bootstrap.Env, channels *bootstrap.Channels, id int) domain.LanguageModelExecutor {
	return &OpenaiChatLanguageChatModelExecutor{env: env, res: channels}
}

type OpenaiChatLanguageChatModelExecutor struct {
	env *bootstrap.Env
	res *bootstrap.Channels
}

func (o *openaiChatLanguageChatModelRequest) Req() {}

type openaiChatLanguageChatModelRequest struct {
	Message []domain.Message `json:"messages"`
	Model   string           `json:"model"`
	//TODO 此处可丰富详细参数 见openai api doc
}

func newOpenaiChatLanguageChatModelRequest(message *[]domain.Message, model string) *openaiChatLanguageChatModelRequest {
	return &openaiChatLanguageChatModelRequest{
		Message: *message,
		Model:   model,
	}
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

	key := fmt.Sprintf("Bearer %s", o.env.ApiOpenaiSecretKey)
	request.Header.Set("Authorization", key) // 请确保使用你自己的API密钥
	request.Header.Set("Content-Type", "application/json")
	return request
}

func (o OpenaiChatLanguageChatModelExecutor) ConfigureProxy(tc *taskchain.TaskContextData) *http.Client {
	return requtil.SetProxy()
}

// TODO 接入rabbitMQ
func (o OpenaiChatLanguageChatModelExecutor) Execute(tc *taskchain.TaskContextData) {
	conn := tc.Conn
	response, err := conn.Client.Do(conn.Request)
	defer response.Body.Close()

	generationResponse := domain.NewGenerationResponse(response, tc.ChatId, err)

	rpcRes := o.res.RpcRes
	if rpcRes == nil {
		rpcRes = make(chan *domain.GenerationResponse, sys.GenerationResponseChannelBuffer)
	}
	rpcRes <- generationResponse
}

// ParseResp 关于
// go-channel方案 设计一个异步线程 始终轮询rpcRes channel  并将轮询所得结果存到map当中 此处只需要GET MAP就可以了
func (o OpenaiChatLanguageChatModelExecutor) ParseResp(tc *taskchain.TaskContextData) domain.ParsedResponse {
	resp := tc.Resp
	body, err := io.ReadAll(resp.Resp.Body)
	if err != nil {
		return nil
	}

	var data *ChatCompletionResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil
	}
	//openAI返回的json中请求体中的文本是一个数组 暂取第0项

	args := data.Choices
	if args == nil {
		return nil
	}
	textBody := args[0]
	generateMessage := domain.OpenAIParsedResponse{
		GenerateText: textBody.Message.Content,
		FinishReason: textBody.FinishReason,
	}
	return &generateMessage
}

type ChatCompletionResponse struct {
	Id         string `json:"id"`
	Object     string `json:"object"`
	CreateTime int64  `json:"created"`
	Model      string `json:"model"`
	Usage      struct {
		PromptTokens     int64 `json:"prompt_tokens"`
		CompletionTokens int64 `json:"completion_tokens"`
		TotalTokens      int64 `json:"total_tokens"`
	} `json:"usage"`

	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}
