package models

import "time"

// CompletionRequest 定义请求结构体
type CompletionRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

// CompletionResponse 定义响应结构体
type CompletionResponse struct {
	Id         string `json:"id"`
	Object     string `json:"object"`
	CreateTime int64  `json:"created"`
	Model      string `json:"model"`
	//openAI返回的json中请求体中的文本是一个数组
	Choices []struct {
		Text         string `json:"text"`
		Index        int64  `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int64 `json:"prompt_tokens"`
		CompletionTokens int64 `json:"completion_tokens"`
		TotalTokens      int64 `json:"total_tokens"`
	} `json:"usage"`
}

func (c CompletionResponse) ToString() {
	//return fmt.Sprintf()
}

// GenerateMessage 生成的返回数据的格式 暂时只有两种信息
type GenerateMessage struct {
	GenerateText string
	FinishReason string
}

// ApiRequestMessage 前端发请求json格式
type ApiRequestMessage struct {
	InputPrompt string `json:"prompt"`
	Model       string `json:"model"`
	MaxToken    int    `json:"max_token"`
}

// 异常返回空对象
func ErrorGeneration() *GenerateMessage {
	return &GenerateMessage{
		GenerateText: "",
		FinishReason: "error",
	}
}

func ErrorCompletionResponse() *CompletionResponse {
	return &CompletionResponse{
		CreateTime: time.Now().Unix(),
		Choices:    nil,
	}
}
