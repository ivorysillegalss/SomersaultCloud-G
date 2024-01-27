package models

// 用户基本信息 TODO
type UserInfo struct {
	UserId   int    `json:"id"  gorm:"primaryKey"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

// 用户Chat相关信息 TODO
type UserChat struct {
	UserId int `json:"id"`
	ChatId int `json:"chat_id"  gorm:"primaryKey"`
}
