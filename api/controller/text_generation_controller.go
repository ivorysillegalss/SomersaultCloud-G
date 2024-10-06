package controller

import (
	"SomersaultCloud/api/dto"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/request"
	"SomersaultCloud/domain"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatController struct {
	chatUseCase domain.ChatUseCase
}

func NewChatController(useCase domain.ChatUseCase) *ChatController {
	return &ChatController{chatUseCase: useCase}
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

func (cc *ChatController) ContextTextChat(c *gin.Context) {
	var askDTO dto.AskDTO
	tokenString := c.Request.Header.Get("token")
	if err := c.BindJSON(&askDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "请求参数解析失败", Code: request.StartChatError})
		return
	}
	isSuccess, parsedResponse, _ := cc.chatUseCase.ContextChat(c, tokenString, askDTO.Ask.BotId, askDTO.Ask.ChatId, askDTO.Ask.Message, askDTO.Adjustment)
	if isSuccess {
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "开启聊天成功", Code: request.StartChatSuccess, Data: parsedResponse})
	} else {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: parsedResponse.GetErrorCause(), Code: request.StartChatError})
	}
}

func (cc *ChatController) StreamContextTextChatSetup(c *gin.Context) {
	var askDTO dto.AskDTO
	tokenString := c.Request.Header.Get("token")
	if err := c.BindJSON(&askDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "请求参数解析失败", Code: request.StartChatError})
		return
	}
	isSuccess, parsedResponse, _ := cc.chatUseCase.StreamContextChatSetup(c, tokenString, askDTO.Ask.BotId, askDTO.Ask.ChatId, askDTO.Ask.Message, askDTO.Adjustment)
	if isSuccess {
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "开启流式聊天成功", Code: request.StartChatSuccess})
	} else {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: parsedResponse.GetErrorCause(), Code: request.StartChatError})
	}
}

func (cc *ChatController) StreamContextTextChatWorker(c *gin.Context) {
	token := c.Request.Header.Get("token")
	//SSE处理函数
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// 获取支持刷新的接口
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "Streaming unsupported!")
		return
	}
	cc.chatUseCase.StreamContextChatWorker(context.Background(), token, c, flusher)
}

func (cc *ChatController) VisionChat(c *gin.Context) {
	var visionDTO dto.VisionDTO
	tokenString := c.Request.Header.Get("token")
	if err := c.BindJSON(&visionDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "请求参数解析失败", Code: request.StartChatError})
		return
	}
	isSuccess, parsedResponse, _ := cc.chatUseCase.DisposableVisionChat(c, tokenString, visionDTO.ChatId, visionDTO.BotId, visionDTO.Message, visionDTO.PicUrl)
	if isSuccess {
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "开启聊天成功", Code: request.StartChatSuccess, Data: parsedResponse})
	} else {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "开启聊天失败", Code: request.StartChatError})
	}
}
