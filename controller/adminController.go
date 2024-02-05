package controller

import (
	"github.com/gin-gonic/gin"
	"mini-gpt/constant"
	"mini-gpt/dto"
	"mini-gpt/service"
	"net/http"
	"strconv"
)

// 管理员获取机器人信息
func AdminGetBot(c *gin.Context) {
	botIdStr := c.Param("botId")
	isOfficialStr := c.Param("isOfficial")
	botId, errBotId := strconv.Atoi(botIdStr)
	isOfficial, errIsOfficial := strconv.Atoi(isOfficialStr)

	// 检查参数解析是否出错
	resultDTO := dto.ResultDTO{}
	if errBotId != nil || errIsOfficial != nil {
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.AdminGetBotError, "参数解析失败", nil))
		return
	}

	bot, err := service.GetBot(botId, isOfficial)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.AdminGetBotError, "管理员获取机器人失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.AdminGetBotSuccess, "管理员获取机器人成功", bot))
	}
}

// 管理员创建新机器人 (默认为官方机器人)
func AdminSaveNewBot(c *gin.Context) {
	var createBot dto.CreateBotDTO

	resultDTO := dto.ResultDTO{}
	if err := c.BindJSON(&createBot); err != nil {
		// 解析请求体失败，返回400
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.AdminSaveNewBotError, "请求参数解析失败", nil))
		return
	}

	err := service.AdminCreateBot(createBot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.AdminSaveNewBotError, "管理员创建新机器人失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.AdminSaveNewBotSuccess, "管理员创建新机器人成功", nil))
	}
}

// 管理员更新机器人
func AdminModifyBot(c *gin.Context) {
	var bot *dto.UpdateBotDTO

	resultDTO := dto.ResultDTO{}
	if err := c.BindJSON(bot); err != nil {
		// 检查参数解析是否出错
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.AdminModifyBotError, "管理员修改机器人失败", nil))
	}

	err := service.AdminUpdateBot(bot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.AdminModifyBotError, "管理员修改机器人失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.FailResp(constant.AdminModifyBotSuccess, "管理员修改机器人成功", nil))
	}
}
