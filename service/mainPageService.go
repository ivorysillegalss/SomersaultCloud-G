package service

import (
	"github.com/google/uuid"
	"mini-gpt/models"
	utils "mini-gpt/utils/jwt"
	"mini-gpt/utils/redisUtils"
	"time"
)

// 渲染主页chat标题等
func InitMainPage(tokenString string) ([]*models.Chat, error) {
	userId, err2 := utils.DecodeToId(tokenString)
	if err2 != nil {
		return nil, err2
	}
	chatTitles, err := models.ShowChatTitle(userId)
	return chatTitles, err
}

func ShareChatHistory(chatId int, day int) string {
	uuidV4 := uuid.New()
	u := uuidV4.String()
	duration := time.Duration(day) * 24 * time.Hour
	//表示保存到数据库中限时3天

	//异步保存 先返回对应sk
	go setHistory(chatId, u, duration)
	return u
}

func setHistory(chatId int, u string, duration time.Duration) {
	//TODO 这里补充代码查询相关历史记录复制一份
	_ = redisUtils.SetStructWithExpire(u, chatId, duration)
}
