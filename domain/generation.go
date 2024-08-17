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
	CacheLuaPollHistory(ctx context.Context, generationResp GenerationResponse)
}
