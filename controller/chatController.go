package controller

import (
	"github.com/gin-gonic/gin"
	"mini-gpt/dto"
	"mini-gpt/models"
	"mini-gpt/service"
	"net/http"
)

// 创建新chat
func CreateChat(c *gin.Context) {
	var chatMessage models.ApiRequestMessage
	c.BindJSON(&chatMessage)
	generateMessage, err := service.LoadingChat(chatMessage)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, generateMessage)
	}
}

// 主页面渲染chat记录
func InitChat(c *gin.Context) {
	var initDTO dto.InitDTO
	c.BindJSON(&initDTO)
	chats, err := service.InitMainPage(initDTO.UserId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, chats)
	}
}

// 调用机器人功能
func CallBot(c *gin.Context) {
	var botDTO dto.BotDTO
	c.BindJSON(&botDTO)
	generateMessage, err := service.DisposableChat(botDTO)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, generateMessage)
	}
}
