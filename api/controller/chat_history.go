package controller

import (
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/request"
	"SomersaultCloud/domain"
	"github.com/gin-gonic/gin"
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
	if tokenString == common.ZeroString {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "请求参数解析失败", Code: request.ShowChatHistoryError})
		return
	}

	chats, err := hmc.chatUseCase.InitMainPage(c, tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "渲染聊天记录失败", Code: request.ShowChatHistoryError})
	} else {
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "渲染聊天记录成功", Code: request.ShowChatHistorySuccess, Data: chats})
	}
}

func (hmc *HistoryMessageController) GetChatHistory(c *gin.Context) {
	chatIdStr := c.Param("chatId")
	chatId, errChatID := strconv.Atoi(chatIdStr)
	// 检查参数解析是否出错
	if errChatID != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "参数解析失败", Code: request.UserGetHistoryError})
		return
	}

	history, err := hmc.chatUseCase.GetChatHistory(c, chatId)
	if err != nil {
		// 获取历史记录失败，返回500状态码
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "获取聊天记录失败", Code: request.UserGetHistoryError})
	} else {
		// 获取历史记录成功，返回200状态码
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "获取聊天记录成功", Code: request.UserGetHistorySuccess, Data: history})
	}
}
