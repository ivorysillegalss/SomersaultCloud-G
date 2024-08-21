package domain

import (
	"context"
	"net/http"
)

// GenerationResponse 程序响应结构体
type GenerationResponse struct {
	Resp   *http.Response
	ChatId int
	Err    error
}

func NewGenerationResponse(response *http.Response, chatId int, err error) *GenerationResponse {
	return &GenerationResponse{
		Resp:   response,
		ChatId: chatId,
		Err:    err,
	}
}

type GenerationRepository interface {
	// CacheLuaPollHistory 由于htttp.response等字段不可序列化 暂且将消费缓存的缓冲map in memory
	CacheLuaPollHistory(ctx context.Context, generationResp GenerationResponse)
	// InMemoryPollHistory 内存存储
	InMemoryPollHistory(ctx context.Context, response *GenerationResponse)
}
