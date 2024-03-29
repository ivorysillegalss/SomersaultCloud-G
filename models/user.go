package models

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"mini-gpt/dao"
	"strconv"
	"sync"
)

// 用户基本信息 TODO
type UserInfo struct {
	UserId   uint   `json:"id"  gorm:"primaryKey;AUTO_INCREMENT"`
	UserName string `gorm:"column:user_name" json:"username"`
	Password string `gorm:"column:pass_word" json:"password"`
}

// 用户Chat相关状态信息
type UserChat struct {
	UserId   int
	Question struct {
		Counter int64
		Doing   bool
	}
	Answer struct {
		Counter int64
		Mu      sync.Mutex
		Buffer  bytes.Buffer
	}
}

func NewUserChat(userId string) *UserChat {
	uid, _ := strconv.Atoi(userId)
	return &UserChat{
		UserId: uid,
		Question: struct {
			Counter int64
			Doing   bool
		}{
			Counter: 0,
			Doing:   false,
		},
		Answer: struct {
			Counter int64
			Mu      sync.Mutex
			Buffer  bytes.Buffer
		}{
			Counter: 0,
		},
	}
}

// 用于接受前端传来的：
type UserServiceReq struct {
	UserName string `gorm:"column:user_name" form:"username" json:"username" binding:"required,min=3,max=15" example:"MiniGpt"`
	Password string `gorm:"column:pass_word" form:"password" json:"password" binding:"required,min=5,max=16" example:"MiniGpt2024"`
	//UserName string `gorm:"column:user_name" json:"username"`
	//Password string `gorm:"column:pass_word" json:"password"`
}

// 响应的用户类型
type UserResp struct {
	UserId   uint   `gorm:"column:user_id" json:"id" form:"id" example:"1"`
	UserName string `gorm:"column:user_name" json:"username" form:"username" example:"MiniGpt"`
	//CreateAt int64  `json:"create_at" form:"create_at"`                  // 创建
}

// 响应携带token
type TokenData struct {
	UserInfo interface{} `json:"user"`
	Token    string      `json:"token"`
}

// 加密程度
const PasswordCost = 12

// 获取特定用户信息
func GetUserInfo1(UserName string) (*UserInfo, error) {
	var userInfo UserInfo
	find := dao.DB.Table("user_info").Where("user_name = ?", UserName).Find(&userInfo)
	fmt.Print("从数据库中获取到的值是：", find.Value)
	err := find.Error
	return &userInfo, err
}

// 创建用户
func CreateUserInfo(u *UserInfo) (*UserInfo, error) {
	var userInfo UserInfo
	err := dao.DB.Table("user_info").Create(&u).First(&userInfo).Error
	return &userInfo, err
}

// SetPassword 设置密码
func (userInfo *UserInfo) SetPassword(password string) error {
	//用来加密传输
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	//将用户的密码设置为处理后的密码
	userInfo.Password = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (userInfo *UserInfo) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(password))
	return err == nil
}

// 通过用户名获取用户id
func GetUserInfo2(UserName string, Password string) (*UserInfo, error) {
	var userInfo UserInfo
	err := dao.DB.Table("user_info").Where("user_name = ? AND pass_word=?", UserName).First(&userInfo).Error
	return &userInfo, err
}
