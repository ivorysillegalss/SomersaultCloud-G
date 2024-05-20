package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"mini-gpt/constant"
	"mini-gpt/dto"
	"mini-gpt/service"
	"net/http"
	"strconv"
)

// 创建新chat
func CreateChat(c *gin.Context) {
	var chatMessage dto.ChatDTO
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
	tokenString := c.Request.Header.Get("token")
	if err := c.BindJSON(&createChatDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.StartChatError, "请求参数解析失败", nil))
		return
	}
	chatId, err := service.CreateChat(&createChatDTO, tokenString)
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
	tokenString := c.Request.Header.Get("token")
	if err := c.BindJSON(&askDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.StartChatError, "请求参数解析失败", nil))
		return
	}
	generateMessage, err := service.ContextChat(&askDTO, tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.StartChatError, "开启聊天失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.StartChatSuccess, "开启聊天成功", generateMessage))
	}
}

// 主页面渲染chat记录
func InitChatHistory(c *gin.Context) {

	resultDTO := dto.ResultDTO{}
	tokenString := c.Request.Header.Get("token")
	if tokenString == constant.ZeroString {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.ShowChatHistoryError, "请求参数解析失败", nil))
		return
	}

	chats, err := service.InitMainPage(tokenString)
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
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.UserGetHistorySuccess, "获取聊天记录失败", nil))
	} else {
		// 获取历史记录成功，返回200状态码
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.UserGetHistoryError, "获取聊天记录成功", history))
	}
}

// 限时密钥形式分享 （生成对应的分享密钥）
func ShareHistoryWithSk(c *gin.Context) {
	chatIdStr := c.Param("chatId")
	chatId, errChatID := strconv.Atoi(chatIdStr)
	// 检查参数解析是否出错
	resultDTO := dto.ResultDTO{}
	if errChatID != nil {
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UserShareHistoryError, "参数解析失败", nil))
		return
	}

	secretKey, err := service.ShareChatHistory(chatId)
	// 生成分享的密钥返回
	if err != nil {
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UserShareHistoryError, "分享历史记录错误", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.UserShareHistorySuccess, "分享历史记录成功", secretKey))
	}
}

// 根据对应的密钥解析 并获取历史记录
func GetSharedHistoryWithSk(c *gin.Context) {
	skStr := c.Param("sk")
	resultDTO := dto.ResultDTO{}
	chatValue, err := service.DecodeSk(skStr)
	// 生成分享的密钥返回
	if errors.Is(err, redis.Nil) {
		c.JSON(http.StatusOK, resultDTO.FailResp(constant.UserShareHistoryNil, "分享密钥不存在或已过期", nil))
	} else if err != nil {
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UserShareHistoryError, "分享历史记录错误", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.UserShareHistorySuccess, "分享历史记录成功", chatValue))
	}
}

// 在原有的历史记录上继续聊天
func ContinueSharedChat(c *gin.Context) {
	skStr := c.Param("sk")
	tokenString := c.Request.Header.Get("token")
	resultDTO := dto.ResultDTO{}
	if skStr == constant.ZeroString || tokenString == constant.ZeroString {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UserGetSharedHistoryError, "请求参数解析失败", nil))
		return
	}
	cloneChatId, err := service.UpdateSharedChat(tokenString, skStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UserGetHistoryError, "使用分享历史记录失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.UserGetSharedHistorySuccess, "使用分享历史记录成功", dto.ShareDTO{CloneChatId: cloneChatId}))
	}
}

// 根据第一次chat的内容更新标题
func InitialTitle(c *gin.Context) {
	var titleDTO dto.TitleDTO
	resultDTO := dto.ResultDTO{}
	if err := c.BindJSON(&titleDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UpdateTitleError, "更新标题失败", nil))
		return
	}
	title, err := service.UpdateInitTitle(&titleDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UpdateTitleError, "更新标题失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.UpdateTitleSuccess, "更新标题成功", title))
	}
}

// 根据用户的输出更新标题
func UpdateTitle(c *gin.Context) {
	var titleDTO dto.TitleDTO
	resultDTO := dto.ResultDTO{}
	if err := c.BindJSON(&titleDTO); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UpdateTitleError, "更新标题失败", nil))
		return
	}
	err := service.UpdateCurrentTitle(&titleDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UpdateTitleError, "更新标题失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.UpdateTitleSuccess, "更新标题成功", nil))
	}
}
