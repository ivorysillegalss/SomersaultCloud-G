package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mini-gpt/constant"
	"mini-gpt/dto"
	"mini-gpt/models"
	"mini-gpt/service"
	"net/http"
)

// 注册  ..用户名是自定义但不能重复，密码也是自定义，但会进行复杂性校验，还有防止sql注入
func Register(c *gin.Context) {
	var reqUser models.UserServiceReq
	resultDTO := dto.ResultDTO{}
	if err := c.ShouldBind(&reqUser); err != nil {
		// 解析请求体失败，返回400状态码
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UserGetError, "请求解析失败", nil))
		return
	}
	fmt.Println("注册的用户是：", reqUser.UserName, reqUser.Password)

	if reqUser.UserName == "" || reqUser.Password == "" {
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UserExistNull, "用户名或密码为空", nil))
		return
	}

	userService := service.GetUserService()
	flag, err := userService.Register(&reqUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.RegisterError, "注册失败", flag))
		return
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.RegisterSuccess, "注册成功", flag))
	}
}

// 登录
func Login(c *gin.Context) {
	fmt.Println("进入到了登录的接口了")
	var reqUser models.UserServiceReq
	resultDTO := dto.ResultDTO{}
	if err := c.ShouldBind(&reqUser); err != nil {
		fmt.Println("用户发过来的格式是有误的，解析不了")
		fmt.Println("错误是：", err)
		c.JSON(http.StatusBadRequest, resultDTO.FailResp(constant.UserGetError, "用户名或密码解析失败", nil))
		return
	}
	fmt.Println("前端刚传过来的是：", reqUser.UserName, reqUser.Password)
	userService := service.GetUserService()
	fmt.Println("成功获取到用户", reqUser)
	//userRespData, err := userService.Login(reqUser)
	_, err := userService.Login(&reqUser)
	fmt.Println("成功调用login")
	if err != nil {
		fmt.Println("登录时发生了错误:", err)
		c.JSON(http.StatusInternalServerError, resultDTO.FailResp(constant.LoginError, "登录失败", nil))
	} else {
		c.JSON(http.StatusOK, resultDTO.SuccessResp(constant.LoginSuccess, "登录成功", "success"))
	}
}
