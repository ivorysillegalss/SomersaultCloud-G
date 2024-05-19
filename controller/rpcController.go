package controller

import (
	"github.com/gin-gonic/gin"
	"mini-gpt/constant"
	"mini-gpt/dto"
	"mini-gpt/service"
	"net/http"
)

//此controller为需要与java模块进行rpc调用所使用到的接口

func Rpc4Title(c *gin.Context) {
	// 使用map[string]interface{}来动态接收数据
	var history dto.TitleDTO

	resultDTO := dto.ResultDTO{}
	if err := c.BindJSON(&history); err != nil {
		// 检查参数解析是否出错
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.Rpc4TitleError, "远程调用获取标题失败", nil))
	}

	titleDTO, err := service.GetTitle(&history)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.Rpc4TitleError, "远程调用获取标题失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.FailResp(constant.Rpc4TitleSuccess, "远程调用获取标题成功", titleDTO))
	}
}
