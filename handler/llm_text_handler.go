package handler

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/sys"
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/log"
	"SomersaultCloud/internal/requtil"
	"bufio"
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/thoas/go-funk"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type OpenaiChatModelExecutor struct {
	env           *bootstrap.Env
	res           *bootstrap.Channels
	generateEvent domain.GenerateEvent
}

func (o *openaiChatLanguageChatModelRequest) Req() {}

type openaiChatLanguageChatModelRequest struct {
	MaxToken int                  `json:"max_tokens"`
	Message  []domain.TextMessage `json:"messages"`
	Model    string               `json:"model"`
	Stream   bool                 `json:"stream"`
	//TODO 此处可丰富详细参数 见openai api doc
}

func newOpenaiChatLanguageChatModelRequest(message *[]domain.TextMessage, model string, stream bool) *openaiChatLanguageChatModelRequest {
	return &openaiChatLanguageChatModelRequest{
		MaxToken: 1000,
		Model:    model,
		Message:  *message,
		Stream:   stream,
	}
}

// AssemblePrompt 组装prompt
// TODO 这个架构上或许可以改进 现在每调用一次都需要 转5次消息格式
//
//	一思路是将他转成哈希
func (o OpenaiChatModelExecutor) AssemblePrompt(tc *domain.AskContextData) *domain.Message {
	var msgs []domain.TextMessage
	historyChat := *tc.History
	var i int
	i = 0
	first := &domain.TextMessage{
		Role:    common.SystemRole,
		Content: tc.SysPrompt,
	}
	msgs = append(msgs, *first)
	if funk.NotEmpty(historyChat) {
		for i < funk.MinInt([]int{len(historyChat), cache.HistoryDefaultWeight}) {
			user := &domain.TextMessage{
				Role:    common.UserRole,
				Content: historyChat[i].ChatAsks.Message,
			}
			msgs = append(msgs, *user)
			asst := &domain.TextMessage{
				Role:    common.GPTRole,
				Content: historyChat[i].ChatGenerations.Message,
			}
			msgs = append(msgs, *asst)
			i++
		}
	}
	last := &domain.TextMessage{
		Role:    common.UserRole,
		Content: tc.Message,
	}
	msgs = append(msgs, *last)
	return &domain.Message{TextMessage: &msgs}
}

func (o OpenaiChatModelExecutor) EncodeReq(tc *domain.AskContextData) *http.Request {

	jsonData, err := json.Marshal(newOpenaiChatLanguageChatModelRequest(tc.HistoryMessage.TextMessage, tc.Model, tc.Stream))
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

func (o OpenaiChatModelExecutor) ConfigureProxy(tc *domain.AskContextData) *http.Client {
	return requtil.SetProxy()
}

func (o OpenaiChatModelExecutor) Execute(tc *domain.AskContextData) {
	conn := tc.Conn
	response, err := conn.Client.Do(conn.Request)
	//若使用stream流式输出 则在发布到消息队列后 下发客户端前 进行消息格式的转换
	//若不使用流式输出 则在主线程中的channel中 再调用下方parse进行消息格式转换
	//why？ 消息队列网络传输需将数据序列化后传 而generationResponse中某些字段如http.Response不可进行序列化

	if response == nil {
		log.GetTextLogger().Warn("fail remote request ,response is nil, userId: " + strconv.Itoa(tc.UserId) + "   chatId: " + strconv.Itoa(tc.ChatId))
		return
	}

	if tc.Stream {
		// 使用 bufio.NewScanner 逐行读取 SSE 响应
		scanner := bufio.NewScanner(response.Body)

		//记录每一次循环所查询到的数据
		var jsonData string
		for scanner.Scan() {
			line := scanner.Text()

			//TODO	此处SSE的信令和前缀以OpenAI的为主，拓展模型可添加
			if line == sys.StreamOverSignal {
				return
			}
			// 过滤空行 并确保解析以 "data: " 开头的行
			if line == common.ZeroString || !strings.HasPrefix(line, sys.StreamPrefix) {
				continue
			}
			// 去除 "data: " 前缀并解析 JSON 数据
			jsonData = line[6:]

			generationResponse := domain.NewStreamGenerationResponse(jsonData, tc.ChatId, err, tc.ExecutorId, tc.UserId)
			o.res.StreamRpcRes <- generationResponse
		}

	} else {
		generationResponse := domain.NewGenerationResponse(response, tc.ChatId, err)
		rpcRes := o.res.RpcRes
		//TODO remove
		if rpcRes == nil {
			rpcRes = make(chan *domain.GenerationResponse, sys.GenerationResponseChannelBuffer)
		}
		rpcRes <- generationResponse
	}
}

// ParseResp 关于
// go-channel方案 设计一个异步线程 始终轮询rpcRes channel  并将轮询所得结果存到map当中 此处只需要GET MAP就可以了
func (o OpenaiChatModelExecutor) ParseResp(tc *domain.AskContextData) (domain.ParsedResponse, string) {

	streamData := tc.Resp.StreamRespData
	if tc.Stream {
		var data *StreamChatCompletionResponse
		err := json.Unmarshal([]byte(streamData), &data)
		if err != nil {
			return nil, common.ZeroString
		}

		args := data.Choices
		if args == nil {
			return nil, common.ZeroString
		}
		textBody := args[0]
		generateMessage := domain.OpenAIParsedResponse{
			//openAI返回的json中请求体中的文本是一个数组 暂取第0项
			//根据流式输出或否修改
			GenerateText: textBody.Delta.Content,
			FinishReason: textBody.FinishReason,
			UserId:       tc.Resp.UserId,
			ExecutorId:   tc.ExecutorId,
			ChatcmplId:   data.Id,
		}
		return &generateMessage, textBody.Delta.Content

	} else {
		resp := tc.Resp
		body, err := io.ReadAll(resp.Resp.Body)
		if err != nil {
			return nil, common.ZeroString
		}

		var data *ChatCompletionResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, common.ZeroString
		}
		//openAI返回的json中请求体中的文本是一个数组 暂取第0项

		args := data.Choices
		if args == nil {
			return nil, common.ZeroString
		}
		textBody := args[0]
		generateMessage := domain.OpenAIParsedResponse{
			GenerateText: textBody.Message.Content,
			FinishReason: textBody.FinishReason,
		}
		return &generateMessage, textBody.Message.Content
	}
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

type StreamChatCompletionResponse struct {
	Id                string      `json:"id"`
	Object            string      `json:"object"`
	Created           int         `json:"created"`
	Model             string      `json:"model"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
	Choices           []Choice    `json:"choices"`
}
type Choice struct {
	Index int `json:"index"`
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}
