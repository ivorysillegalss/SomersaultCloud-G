package controller

import (
	"github.com/gin-gonic/gin"
	"mini-gpt/dto"
	"mini-gpt/service"
	"net/http"
	"strconv"
)

func AdminGetBot(c *gin.Context) {
	botIdStr := c.Param("botId")
	isOfficialStr := c.Param("isOfficial")
	botId, _ := strconv.Atoi(botIdStr)
	isOfficial, _ := strconv.Atoi(isOfficialStr)
	bot, err := service.GetBot(botId, isOfficial)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, bot)
	}
}

func AdminSaveNewBot(c *gin.Context) {
	var createBot dto.CreateBotDTO
	_ = c.BindJSON(createBot)
	err := service.AdminCreateBot(createBot)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, nil)
	}
}
