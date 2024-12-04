package domain

import (
	"context"
	"net/http"
)

// GenerationResponse 程序响应结构体
// 非Stream式的直接使用*http.Response对数据进行传输
// 而流式输出因为结尾数据格式不一样 并且穿插空行 需要先对输出后的数据进行前置处理 此处处理成为	StreamRespData string 格式
type GenerationResponse struct {
	Resp           *http.Response
	StreamRespData string
	ChatId         int
	Err            error
	ExecutorId     int
	UserId         int
}

type StreamGenerationReadyStorageData struct {
	Records     *[]*Record
	UserContent string
	ChatId      int
	UserId      int
	BotId       int
	Title       string
}

func NewGenerationResponse(response *http.Response, chatId int, err error) *GenerationResponse {
	return &GenerationResponse{Resp: response, ChatId: chatId, Err: err}
}

func NewStreamGenerationResponse(streamData string, chatId int, err error, executorId int, userId int) *GenerationResponse {
	return &GenerationResponse{StreamRespData: streamData, ChatId: chatId, Err: err, ExecutorId: executorId, UserId: userId}
}

type GenerationRepository interface {
	// CacheLuaPollHistory 由于http.response等字段不可序列化 暂且将消费缓存的缓冲map in memory
	CacheLuaPollHistory(ctx context.Context, generationResp GenerationResponse)
	// InMemoryPollHistory 内存存储
	InMemoryPollHistory(ctx context.Context, response *GenerationResponse)

	ReadyStreamDataStorage(ctx context.Context, ready StreamGenerationReadyStorageData)
	GetStreamDataStorage(ctx context.Context, userId int) *AskContextData
}

type GenerationCron interface {
	AsyncPollerGeneration()
}
