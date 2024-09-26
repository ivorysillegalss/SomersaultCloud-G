package controller

import (
	"SomersaultCloud/api/dto"
	"SomersaultCloud/constant/request"
	"SomersaultCloud/domain"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"net/http"
	"strconv"
)

type HistoryMessageController struct {
	chatUseCase domain.ChatUseCase
}

func NewHistoryMessageController(useCase domain.ChatUseCase) *HistoryMessageController {
	return &HistoryMessageController{chatUseCase: useCase}
}

func (hmc *HistoryMessageController) HistoryTitle(c *gin.Context) {
	tokenString := c.Request.Header.Get("token")
	botIdStr := c.Param("botId")
	botId, errBotId := strconv.Atoi(botIdStr)
	if funk.IsEmpty(tokenString) || funk.NotEmpty(errBotId) {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "请求参数解析失败", Code: request.ShowChatHistoryError})
		return
	}

	chats, err := hmc.chatUseCase.InitMainPage(c, tokenString, botId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "渲染聊天记录失败", Code: request.ShowChatHistoryError})
	} else {
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "渲染聊天记录成功", Code: request.ShowChatHistorySuccess, Data: chats})
	}
}

func (hmc *HistoryMessageController) GetChatHistory(c *gin.Context) {
	chatIdStr := c.Param("chatId")
	botIdStr := c.Param("botId")
	tokenString := c.Request.Header.Get("token")
	chatId, errChatID := strconv.Atoi(chatIdStr)
	botId, errBotID := strconv.Atoi(botIdStr)
	// 检查参数解析是否出错
	if funk.NotEmpty(errBotID) || funk.NotEmpty(errChatID) || funk.IsEmpty(tokenString) {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "参数解析失败", Code: request.UserGetHistoryError})
		return
	}

	history, err := hmc.chatUseCase.GetChatHistory(c, chatId, botId, tokenString)
	if err != nil {
		// 获取历史记录失败，返回500状态码
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "获取聊天记录失败", Code: request.UserGetHistoryError})
	} else {
		// 获取历史记录成功，返回200状态码
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "获取聊天记录成功", Code: request.UserGetHistorySuccess, Data: history})
	}
}

func (hmc *HistoryMessageController) UpdateInitTitle(c *gin.Context) {
	var titleDTO dto.TitleDTO
	tokenString := c.Request.Header.Get("token")
	if err := c.BindJSON(&titleDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "参数处理错误 更新标题失败", Code: request.UpdateTitleError})
		return
	}
	title, err := hmc.chatUseCase.GenerateUpdateTitle(c, &titleDTO.Messages, tokenString, titleDTO.ChatId)
	if err != nil {
		// 获取历史记录失败，返回500状态码
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "获取聊天记录失败", Code: request.UpdateTitleError})
	} else {
		// 获取历史记录成功，返回200状态码
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "获取聊天记录成功", Code: request.UpdateTitleSuccess, Data: &dto.TitleDTO{Title: title}})
	}
}

func (hmc *HistoryMessageController) InputTitle(c *gin.Context) {
	var titleDTO dto.TitleDTO
	tokenString := c.Request.Header.Get("token")
	if err := c.BindJSON(&titleDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "参数处理错误 更新标题失败", Code: request.UpdateTitleError})
		return
	}
	isSuccess := hmc.chatUseCase.InputUpdateTitle(c, titleDTO.Title, tokenString, titleDTO.ChatId, titleDTO.BotId)
	if !isSuccess {
		// 获取历史记录失败，返回500状态码
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "获取聊天记录失败", Code: request.UpdateTitleError})
	} else {
		// 获取历史记录成功，返回200状态码
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "获取聊天记录成功", Code: request.UpdateTitleSuccess})
	}
}
