package controller

import (
	"SomersaultCloud/api/dto"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/request"
	"SomersaultCloud/domain"
	"SomersaultCloud/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatController struct {
	chatUseCase domain.ChatUseCase
}

func NewChatController() *ChatController {
	return &ChatController{chatUseCase: usecase.NewChatUseCase()}
}

func (cc *ChatController) InitNewChat(c *gin.Context) {
	var createChatDTO dto.CreateChatDTO

	tokenString := c.Request.Header.Get("token")
	if err := c.BindJSON(&createChatDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "请求参数解析失败", Code: request.StartChatError})
		return
	}
	chatId := cc.chatUseCase.InitChat(c, tokenString, createChatDTO.BotId)
	if chatId == common.FalseInt {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "开启聊天失败", Code: request.StartChatError})
	} else {
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "开启聊天成功", Code: request.StartChatSuccess, Data: chatId})
	}

}
