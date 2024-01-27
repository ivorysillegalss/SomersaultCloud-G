package controller

import (
	"github.com/gin-gonic/gin"
	"mini-gpt/dto"
	"mini-gpt/models"
	"mini-gpt/service"
	"net/http"
)

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
