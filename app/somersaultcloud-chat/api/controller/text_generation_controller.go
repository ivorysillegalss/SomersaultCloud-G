package controller

import (
	"SomersaultCloud/app/somersaultcloud-chat/api/dto"
	"SomersaultCloud/app/somersaultcloud-chat/constant/common"
	"SomersaultCloud/app/somersaultcloud-chat/constant/request"
	"SomersaultCloud/app/somersaultcloud-chat/domain"
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gin-gonic/gin"
	"github.com/hertz-contrib/sse"
	"github.com/thoas/go-funk"
)

type ChatController struct {
	chatUseCase domain.ChatUseCase
}

func NewChatController(useCase domain.ChatUseCase) *ChatController {
	return &ChatController{chatUseCase: useCase}
}

func (cc *ChatController) InitNewChat(ctx context.Context, c *app.RequestContext) {
	var createChatDTO dto.CreateChatDTO

	tokenString := c.Request.Header.Get("token")
	if err := c.BindJSON(&createChatDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "请求参数解析失败", Code: request.StartChatError})
		return
	}
	chatId := cc.chatUseCase.InitChat(ctx, tokenString, createChatDTO.BotId)
	if chatId == common.FalseInt {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "开启聊天失败", Code: request.StartChatError})
	} else {
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "开启聊天成功", Code: request.StartChatSuccess, Data: chatId})
	}

}

func (cc *ChatController) ContextTextChat(ctx context.Context, c *app.RequestContext) {
	var askDTO dto.AskDTO
	tokenString := c.Request.Header.Get("token")
	if err := c.BindJSON(&askDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "请求参数解析失败", Code: request.StartChatError})
		return
	}
	isSuccess, parsedResponse, _ := cc.chatUseCase.ContextChat(ctx, tokenString, askDTO.Ask.BotId, askDTO.Ask.ChatId, askDTO.Ask.Message, askDTO.Adjustment, askDTO.Model)
	if isSuccess {
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "开启聊天成功", Code: request.StartChatSuccess, Data: parsedResponse})
	} else {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: parsedResponse.GetErrorCause(), Code: request.StartChatError})
	}
}

func (cc *ChatController) StreamContextTextChatSetup(ctx context.Context, c *app.RequestContext) {
	var askDTO dto.AskDTO
	tokenString := c.Request.Header.Get("token")
	if err := c.BindJSON(&askDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "请求参数解析失败", Code: request.StartChatError})
		return
	}
	isSuccess, parsedResponse, _ := cc.chatUseCase.StreamContextChatSetup(ctx, tokenString, askDTO.Ask.BotId, askDTO.Ask.ChatId, askDTO.Ask.Message, askDTO.Adjustment)
	if isSuccess {
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "开启流式聊天成功", Code: request.StartChatSuccess})
	} else {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: parsedResponse.GetErrorCause(), Code: request.StartChatError})
	}
}

func (cc *ChatController) StreamContextTextChatWorker(ctx context.Context, c *app.RequestContext) {
	token := c.Request.Header.Get("token")
	//TODO REMOVE:测试用
	//token := "eyJhbGciOiJIUzI1NiJ9.eyJ1aWQiOi0xLCJleHAiOjEwMDAwMTcyNTQ2NTUzN30.nlW5kKPgBZwqdxafrt_VTEPwVg7x9OWWOsKTM4Xk0B4"

	//SSE处理函数
	c.Response.Header.Set("Content-Type", "text/event-stream")
	c.Response.Header.Set("Cache-Control", "no-cache")
	c.Response.Header.Set("Connection", "keep-alive")

	// 必须在第一次调用之前设置状态代码和响应标头
	c.SetStatusCode(http.StatusOK)
	stream := sse.NewStream(c)
	cc.chatUseCase.StreamContextChatWorker(ctx, token, stream)
}

func (cc *ChatController) StreamContextChatStorage(ctx context.Context, c *app.RequestContext) {
	token := c.Request.Header.Get("Token")
	if funk.IsEmpty(token) {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "请求参数解析失败", Code: request.StartChatError})
		return
	}
	isSuccess := cc.chatUseCase.StreamContextStorage(ctx, token)
	if isSuccess {
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "流式信息存储成功", Code: request.StorageStreamTextSuccess})
	} else {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "流式信息存储失败", Code: request.StorageStreamTextFail})
	}
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
