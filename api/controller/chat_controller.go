package controller

import (
	"SomersaultCloud/domain"
	"SomersaultCloud/usecase"
)

type ChatController struct {
	chatUseCase domain.ChatUseCase
}

func NewChatController() *ChatController {
	return &ChatController{chatUseCase: usecase.NewChatUseCase()}
}

func (c *ChatController) InitNewChat() {
	c.chatUseCase.InitChat()
}
