package service

import "mini-gpt/models"

// 渲染主页chat标题等
func InitMainPage(userId int) ([]*models.Chat, error) {
	chatTitles, err := models.ShowChatTitle(userId)
	return chatTitles, err
}
