package channel

import "net/http"

// GenerationResponse
// TODO 分割domain和基础 GenerationResponse迁移domain Response作为基础开一个
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
