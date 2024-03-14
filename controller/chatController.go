package controller

import (
	"github.com/gin-gonic/gin"
	"mini-gpt/constant"
	"mini-gpt/dto"
	"mini-gpt/models"
	"mini-gpt/service"
	"net/http"
	"strconv"
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
	generateMessage, err := service.LoadingChat(&chatMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.StartChatError, "开启聊天失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.StartChatSuccess, "开启聊天成功", generateMessage))
	}
}

// 初始化新chat
func InitNewChat(c *gin.Context) {
	var createChatDTO dto.CreateChatDTO

	resultDTO := dto.ResultDTO{}
	if err := c.BindJSON(&createChatDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.StartChatError, "请求参数解析失败", nil))
		return
	}
	chatId, err := service.CreateChat(&createChatDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.StartChatError, "开启聊天失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.StartChatSuccess, "开启聊天成功", chatId))
	}
}

// 使用chat
func CallContextChat(c *gin.Context) {
	var askDTO dto.AskDTO

	resultDTO := dto.ResultDTO{}
	if err := c.BindJSON(&askDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.StartChatError, "请求参数解析失败", nil))
		return
	}
	generateMessage, err := service.ContextChat(&askDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.StartChatError, "开启聊天失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.StartChatSuccess, "开启聊天成功", generateMessage))
	}
}

// 主页面渲染chat记录
func InitChatHistory(c *gin.Context) {
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

	generateMessage, err := service.DisposableChat(&botDTO)
	if err != nil {
		// 调用机器人失败，返回500状态码
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.CallBotError, "调用机器人失败", nil))
	} else {
		// 调用机器人成功，返回200状态码
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.CallBotSuccess, "调用机器人成功", generateMessage))
	}
}

// 获取聊天记录
func GetChatHistory(c *gin.Context) {
	chatIdStr := c.Param("chatId")
	chatId, errChatID := strconv.Atoi(chatIdStr)
	// 检查参数解析是否出错
	resultDTO := dto.ResultDTO{}
	if errChatID != nil {
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UserGetHistoryError, "参数解析失败", nil))
		return
	}

	history, err := service.GetChatHistory(chatId)
	if err != nil {
		// 获取历史记录失败，返回500状态码
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.UserGetHistorySuccess, "获取聊天记录成功", nil))
	} else {
		// 获取历史记录成功，返回200状态码
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.UserGetHistoryError, "获取聊天记录失败", history))
	}
}

// 限时密钥形式分享
func ShareHistory(c *gin.Context) {
	chatIdStr := c.Param("chatId")
	durationDayStr := c.Param("ddl")
	chatId, errChatID := strconv.Atoi(chatIdStr)
	duration, errDuration := strconv.Atoi(durationDayStr)
	// 检查参数解析是否出错
	resultDTO := dto.ResultDTO{}
	if errChatID != nil || errDuration != nil {
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UserShareHistoryError, "参数解析失败", nil))
		return
	}

	secretKey := service.ShareChatHistory(chatId, duration)
	// 生成分享的密钥返回
	c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.UserShareHistorySuccess, "分享历史记录错误", secretKey))
}
