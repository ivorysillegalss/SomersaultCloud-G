package domain

import (
	"net/http"
	"time"
)

// 目前思路：所有模型的所有任一请求方式 都需要实现下方languagemodel
// 但是同一个厂商的不同请求方式 eg 文生文 or 文生图 可以将他的请求参数综合
// TODO 还是希望这里executor能跟domain解耦，看着有点难受，虽然不违法开发规范
type LanguageModelExecutor interface {
	AssemblePrompt(tc *AskContextData) *Message
	EncodeReq(tc *AskContextData) *http.Request
	// ConfigureProxy 非必实现 根据api是否被墙
	ConfigureProxy(tc *AskContextData) *http.Client
	Execute(tc *AskContextData)
	ParseResp(tc *AskContextData) (ParsedResponse, string)
}

// ParsedResponse 转码的历史记录抽象接口
type ParsedResponse interface {
	GetGenerateText() string
	GetErrorCause() string
	GetIndex() int
	GetIdentity() int
	GetFinishReason() string
	GetExecutorId() int
}

// OpenAIParsedResponse 目前的实现在生成的时候 若出现错误直接将错误置为GenerateText
// 也就是说两者方法实现上本质上一样 只是长得不一样
//
//	流式生成的时候 通过UserId作为标识 保证时序性
type OpenAIParsedResponse struct {
	UserId       int
	GenerateText string
	FinishReason string
	Index        int //索引 流式输出的索引
	ExecutorId   int // 执行器ID 用于标识发送数据时反序列化进度
}

func (o *OpenAIParsedResponse) GetGenerateText() string { return o.GenerateText }

func (o *OpenAIParsedResponse) GetErrorCause() string { return o.GenerateText }

func (o *OpenAIParsedResponse) GetIndex() int { return o.Index }

func (o *OpenAIParsedResponse) GetIdentity() int { return o.UserId }

func (o *OpenAIParsedResponse) GetFinishReason() string { return o.FinishReason }

// TODO
func (o *OpenAIParsedResponse) GetExecutorId() int { return o.ExecutorId }

type LanguageModelRequest interface {
	Req()
}

type Message struct {
	TextMessage *[]TextMessage
	TypeMessage *[]TypeMessage
}

type TextMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TypeMessage struct {
	Role    string     `json:"role"`
	Content []TypeInfo `json:"content"`
}

//type Message interface {getContent()}
//func (t *TextMessage) getContent() {}
//func (y *TypeMessage) getContent() {}

type TypeInfo interface {
	GetType()
}

type TextType struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (t *TextType) GetType() {}

type ImageUrlType struct {
	Type     string    `json:"type"`
	ImageUrl *ImageUrl `json:"image_url"`
}

type ImageUrl struct {
	Url    string `json:"url"`
	Detail string `json:"detail"`
}

func (i *ImageUrlType) GetType() {}

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
