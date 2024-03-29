package models

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"mini-gpt/constant"
	"mini-gpt/dao"
	"mini-gpt/utils/redisUtils"
	"strconv"
	"time"
)

// 一次Chat
type Chat struct {
	//这里将chatId直接更名为ID 方便gorm进行主键回显
	ID             int       `json:"chat_id"  gorm:"primaryKey"`
	UserId         int       `json:"user_id"`
	BotId          int       `json:"bot_id"`
	Title          string    `json:"title"`
	LastUpdateTime int64     `json:"last_update_time"`
	IsDelete       bool      `json:"is_delete"`
	Records        *[]Record `json:"records"`
}

// Record 一次问答
type Record struct {
	RecordId        int             `json:"record_id"`
	ChatAsks        *ChatAsk        `json:"chat_asks"`
	ChatGenerations *ChatGeneration `json:"chat_generations"`
	//Weights         float64
}

// ChatAsk 一次问题
type ChatAsk struct {
	RecordId int    `json:"record_id"`
	ChatId   int    `json:"chat_id"`
	Message  string `json:"message"`
	BotId    int    `json:"bot_id"`
}

// ChatGeneration 一次生成
type ChatGeneration struct {
	RecordId int    `json:"record_id"`
	ChatId   int    `json:"chat_id"`
	Message  string `json:"message"`
}

// ShowChatTitle 主页面展示已有chat的标题
func ShowChatTitle(userId int) ([]*Chat, error) {
	var chats []*Chat
	err := dao.DB.Table("chat").Where("is_delete = ?", 0).Where("user_id = ?", userId).Find(&chats).Error
	return chats, err
}

// 创建新的chat初始化
func CreateNewChat(userId int, botId int) (int, error) {
	chat := &Chat{
		UserId:         userId,
		BotId:          botId,
		Title:          "init",
		LastUpdateTime: time.Now().Unix(),
		IsDelete:       false,
	}
	if err := dao.DB.Table("chat").Create(chat).Error; err != nil {
		return -1, err
	}
	return chat.ID, nil
}

// 获取聊天记录错误的时候返回
func ErrorRecord() *[]*Record {
	return new([]*Record)
}

func GetChatHistoryForChat(chatId int) (*[]*Record, error) {
	//返回一个存放record结构体的 指针的切片的 指针

	//var ask []*ChatAsk
	//dao.DB.Table("chat_ask").Where("chat_id = ?", chatId).Find(&ask).Order("recordId asc")
	//var generation []*ChatGeneration
	//dao.DB.Table("chat_generation").Where("chat_id = ?", chatId).Find(&generation).Order("recordId asc")
	var records []*Record

	records, err := redisUtils.GetStruct[[]*Record](constant.ChatCache + strconv.Itoa(chatId))
	//去redis里查

	//此处可以优化逻辑

	if errors.Is(redis.Nil, err) {

		//redis中查不到的时候去mysql里查
		//if err := dao.DB.Joins("JOIN chat_generation ON chat_ask.record_id = chat_generation.record_id").Where("chat_id = ?", chatId).
		//	Find(&records).Limit(10).Order("recordId asc").Error; err != nil {
		//	return ErrorRecord(), err
		//}

		err := dao.DB.Raw("SELECT * FROM chat_ask JOIN chat_generation ON chat_ask.record_id = chat_generation.record_id WHERE chat_ask.chat_id = ?", chatId).Scan(&records).Limit(10).Order("recordId asc").Error
		if err != nil {
			// 处理错误
		}

	} else if err != nil && !errors.Is(redis.Nil, err) {
		//出现了其他错误
		return ErrorRecord(), err
	}

	return &records, nil
}

// 写入数据库的聊天记录映射类
type recordToStruct struct {
	ID     int `gorm:"primaryKey column:record_id" `
	ChatId int `gorm:"column:chat_id"`
}

// 保存记录
func SaveRecord(record *Record, chatId int) error {
	r := &recordToStruct{
		ChatId: record.ChatAsks.ChatId,
	}
	if err := dao.DB.Table("chat_record_id").Save(r).Error; err != nil {
		return err
	}

	//由上方将recordId写入数据库 主键回显获得ID 赋值给ask及generation两张表

	record.ChatAsks.RecordId = record.RecordId
	record.ChatGenerations.RecordId = record.RecordId

	if err := dao.DB.Table("chat_ask").Save(record.ChatAsks).Error; err != nil {
		return err
	}
	if err := dao.DB.Table("chat_generation").Save(record.ChatGenerations).Error; err != nil {
		return err
	}
	if err := dao.DB.Table("chat").Where("chat_id = ?", chatId).Update("last_update_time", time.Now().Unix()).Error; err != nil {
		return err
	}
	return nil
}

// 获得聊天记录
func GetChatHistory(chatId int) ([]*Record, error) {
	var history []int
	if err := dao.DB.Table("record_info").Where("chat_id = ?", chatId).Find(&history).Error; err != nil {
		return nil, err
	}
	var records []*Record
	if err := dao.DB.Raw("SELECT * FROM chat_ask JOIN chat_generation ON chat_ask.record_id = chat_generation.record_id WHERE chat_ask.chat_id = ?", chatId).Scan(&records).Limit(10).Order("recordId asc").Error; err != nil {
		return nil, err
	}
	return records, nil
}
