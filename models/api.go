package models

import (
	"encoding/json"
	"mini-gpt/constant"
)

type BaseModel interface {
	completion()
}

func (t *TextCompletionResponse) completion() {}
func (c *ChatCompletionResponse) completion() {}

type ApiRequestMessage interface {
	request()
}

func (t *TextCompletionRequest) request() {}
func (c *ChatCompletionRequest) request() {}

// ChatCompletionRequest 定义大语言请求结构体
type ChatCompletionRequest struct {
	CompletionRequest
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// TextCompletionRequest 定义文本补全请求结构体
type TextCompletionRequest struct {
	CompletionRequest
	Prompt string `json:"prompt"`
}

type CompletionRequest struct {
	MaxTokens int    `json:"max_tokens"`
	Model     string `json:"model"`
}

// CompletionResponse 定义响应结构体
type CompletionResponse struct {
	Id         string `json:"id"`
	Object     string `json:"object"`
	CreateTime int64  `json:"created"`
	Model      string `json:"model"`
	Usage      struct {
		PromptTokens     int64 `json:"prompt_tokens"`
		CompletionTokens int64 `json:"completion_tokens"`
		TotalTokens      int64 `json:"total_tokens"`
	} `json:"usage"`
}

// openAI返回的json中请求体中的文本是一个数组
// 文本补全 响应数据格式
type TextCompletionResponse struct {
	CompletionResponse
	Choices []struct {
		Text         string `json:"text"`
		Index        int64  `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// chat模型响应数据格式
type ChatCompletionResponse struct {
	//CompletionResponse

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

func (c CompletionResponse) ToString() {
	//return fmt.Sprintf()
}

// GenerateMessage 生成的返回数据的格式 暂时只有两种信息
type GenerateMessage struct {
	GenerateText string
	FinishReason string
}

//// ApiRequestMessage 前端发请求json格式
//type ApiRequestMessage struct {
//	InputPrompt string `json:"prompt"`
//	Model       string `json:"model"`
//	MaxToken    int    `json:"max_token"`
//}

// 异常返回空对象
func ErrorGeneration() *GenerateMessage {
	return &GenerateMessage{
		GenerateText: "",
		FinishReason: "error",
	}
}

// 加载bot配置错误
func ErrorBotConfig() *BotConfig {
	return &BotConfig{
		BotId:      0,
		InitPrompt: "",
		Model:      "",
	}
}

// 调用api之后的状态
type ExecuteStatus struct {
	Status     string
	StatusCode int
}

// 系统错误
func ErrorCompletionResponse() BaseModel {
	return nil
}

// 运行错误 用户的锅
func ExceptionCompletionResponse(exceptionMessage string) BaseModel {
	//这里草草包装了一下 未测试 感觉还可以优化
	return nil
}

// 错误请求信息
//func ErrorApiRequestMessage() ApiRequestMessage {
//	return new(ApiRequestMessage)
//}

// 文本补全模型包装信息 gpt-3.5-turbo-instruct
func EncodeTextCompletionJsonData(reqMessage ApiRequestMessage) ([]byte, error) {
	var textCompletion *TextCompletionRequest
	if t, ok := reqMessage.(*TextCompletionRequest); !ok {
		return nil, nil
		//TODO 错误处理
	} else {
		textCompletion = t
	}
	//这个判空的过程可以优化在结构体 models层中
	if textCompletion.MaxTokens == constant.ZeroInt {
		textCompletion.MaxTokens = constant.DefaultMaxToken
	}

	// 构造请求体
	data := TextCompletionRequest{
		CompletionRequest: CompletionRequest{
			MaxTokens: textCompletion.MaxTokens,
			Model:     textCompletion.Model,
		},
		//Model:     "gpt-3.5-turbo-instruct", // 替换为当前可用的模型
		Prompt: textCompletion.Prompt,
	}
	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func EncodeChatCompletionJsonData(reqMessage ApiRequestMessage) ([]byte, error) {
	var chatCompletion *ChatCompletionRequest
	if c, ok := reqMessage.(*ChatCompletionRequest); !ok {
		//TODO 错误处理
		return nil, nil
	} else {
		chatCompletion = c
	}
	//这个判空的过程可以优化在结构体 models层中
	if chatCompletion.MaxTokens == constant.ZeroInt {
		chatCompletion.MaxTokens = constant.DefaultMaxToken
	}

	// 构造请求体
	data := ChatCompletionRequest{
		CompletionRequest: CompletionRequest{
			MaxTokens: chatCompletion.MaxTokens,
			Model:     chatCompletion.Model,
		},
		Messages: chatCompletion.Messages,
	}
	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
