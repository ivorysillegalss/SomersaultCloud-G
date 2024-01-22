package controller

import (
	"github.com/gin-gonic/gin"
	"mini-gpt/models"
	"mini-gpt/service"
	"net/http"
)

func CreateChat(c *gin.Context) {
	var chatMessage models.PromptMessage
	c.BindJSON(&chatMessage)
	generateMessage, err := service.LoadingChat(chatMessage)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, generateMessage)
	}
}
