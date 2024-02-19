package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"mini-gpt/models"
	"mini-gpt/setting"
	util "mini-gpt/utils/jwt"
	"sync"
)

// 单例模式：懒汉式实例化
var UserServiceInstance *UserService
var UserServiceOnce sync.Once

type UserService struct {
}

func GetUserService() *UserService {
	UserServiceOnce.Do(func() {
		UserServiceInstance = &UserService{}
	})
	return UserServiceInstance
}

// 注册
func (us *UserService) Register(req *models.UserServiceReq) (flag bool, err error) {
	var logger = setting.GetLogger()

	u, err := models.GetUserInfo1(req.UserName)

	switch err {
	case gorm.ErrRecordNotFound: //库里没有记录

		u = &models.UserInfo{
			UserName: req.UserName,
		}

		//密码加密存储
		if err = u.SetPassword(req.Password); err != nil {
			logger.Error(err)
			return false, err
		}

		//创建用户
		_, err := models.CreateUserInfo(u)
		if err != nil {
			logger.Error(err)
			return false, err
		}

		return true, nil

	case nil:
		err = errors.New("用户已存在")
		logger.Error(err)
		return false, err
	default:
		return false, err
	}
}

// 登录
func (us *UserService) Login(req *models.UserServiceReq) (resp interface{}, err error) {
	var logger = setting.GetLogger()
	fmt.Print("获取到的用户名为：", req.UserName)
	user, err := models.GetUserInfo1(req.UserName)
	if err == gorm.ErrRecordNotFound {
		//err = errors.New("用户不存在")
		logger.Println(err)
		logger.Error(err)
		return "", err
	}

	//校验密码
	if !user.CheckPassword(req.Password) {
		err = errors.New("账号/密码错误")
		logger.Println(err)
		logger.Error(err)
		return "", err
	}
	//用户id的获取：

	token, err := util.GenerateToken(user.UserId, user.UserName, 0)
	if err != nil {
		logger.Println(err)
		logger.Error(err)
		return "", err
	}

	userResp := &models.UserResp{
		UserId:   user.UserId,
		UserName: user.UserName,
	}
	userRespData := &models.TokenData{
		UserInfo: userResp,
		Token:    token,
	}

	return userRespData, err
}
