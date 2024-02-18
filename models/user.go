package models

import (
	"bytes"
	"strconv"
	"sync"
)

// 用户基本信息 TODO
type UserInfo struct {
	UserId   int    `json:"id"  gorm:"primaryKey"`
	UserName string `json:"username"`
	Password string `json:"password"`
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
