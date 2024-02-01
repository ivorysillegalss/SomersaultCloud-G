package controller

import (
	"github.com/gin-gonic/gin"
	"mini-gpt/constant"
	"mini-gpt/dto"
	"mini-gpt/models"
	"mini-gpt/service"
	"net/http"
)

// 创建新chat
func CreateChat(c *gin.Context) {
	var chatMessage models.ApiRequestMessage

	resultDTO := dto.ResultDTO{}
	if err := c.BindJSON(&chatMessage); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.StartChatError, "请求参数解析失败", nil))
		return
	}
	generateMessage, err := service.LoadingChat(chatMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.StartChatError, "开启聊天失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.StartChatSuccess, "开启聊天成功", generateMessage))
	}
}

// 主页面渲染chat记录
func InitChat(c *gin.Context) {
	var initDTO dto.InitDTO

	resultDTO := dto.ResultDTO{}
	if err := c.BindJSON(&initDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.ShowChatHistoryError, "请求参数解析失败", nil))
		return
	}

	chats, err := service.InitMainPage(initDTO.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.ShowChatHistoryError, "渲染聊天记录失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.ShowChatHistorySuccess, "渲染聊天记录成功", chats))
	}
}

// 调用机器人功能
func CallBot(c *gin.Context) {
	var botDTO dto.ExecuteBotDTO

	resultDTO := dto.ResultDTO{}
	if err := c.BindJSON(&botDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.CallBotError, "请求参数解析失败", nil))
		return
	}

	generateMessage, err := service.DisposableChat(botDTO)
	if err != nil {
		// 调用机器人失败，返回500状态码
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.CallBotError, "调用机器人失败", nil))
	} else {
		// 调用机器人成功，返回200状态码
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.CallBotSuccess, "调用机器人成功", generateMessage))
	}
}
