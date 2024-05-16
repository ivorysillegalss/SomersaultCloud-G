package service

import (
	"github.com/google/uuid"
	"mini-gpt/constant"
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

//点对点分享历史记录 TODO
//func ShareChatHistory(chatId int, day int, tokenString string) (string, error) {
//	userId, err := utils.DecodeToId(tokenString)
//	if err != nil {
//		return constant.ZeroString, err
//	}
//
//	chatHistory, err := GetChatHistory(chatId)
//	if err != nil {
//		return constant.ZeroString, err
//	}
//
//	//从第0项获取这一chat的botId
//	a := *chatHistory
//	cloneChatId, err := models.CreateNewChat(userId, a[0].ChatAsks.BotId)
//	if err != nil {
//		return constant.ZeroString, err
//	}
//
//	//赋值聊天记录 除了chatId和userId以外其他都是和原来一样的
//	count := len(*chatHistory)
//	for i := 0; i < count; i++ {
//
//		//重新赋值问和答的chatID
//		a[i].ChatAsks.ChatId = cloneChatId
//		a[i].ChatGenerations.ChatId = cloneChatId
//
//		err := models.SaveRecord(a[i], cloneChatId)
//		if err != nil {
//			return constant.ZeroString, err
//		}
//	}
//
//	uuidV4 := uuid.New()
//	u := uuidV4.String()
//	duration := time.Duration(day) * 24 * time.Hour * 3
//	//表示保存到数据库中限时3天
//	err = redisUtils.SetStructWithExpire(u, chatId, duration)
//	if err != nil {
//		return constant.ZeroString, err
//	}
//
//	//异步保存 先返回对应sk
//	//go setHistory(chatId, u, duration)
//	return u, nil
//}

// 密钥形式分享
func ShareChatHistory(chatId int) (string, error) {

	chatHistory, err := GetChatHistory(chatId)
	if err != nil {
		return constant.ZeroString, err
	}

	//从第0项获取这一chat的botId
	a := *chatHistory
	//这里暂时是用0来代表 接受分享的人的userId
	cloneChatId, err := models.CreateNewChat(constant.ZeroInt, a[0].ChatAsks.BotId)
	if err != nil {
		return constant.ZeroString, err
	}

	//赋值聊天记录 除了chatId和userId以外其他都是和原来一样的
	count := len(*chatHistory)
	for i := 0; i < count; i++ {

		//重新赋值问和答的chatID
		a[i].ChatAsks.ChatId = cloneChatId
		a[i].ChatGenerations.ChatId = cloneChatId

		err := models.SaveRecord(a[i], cloneChatId)
		if err != nil {
			return constant.ZeroString, err
		}
	}

	uuidV4 := uuid.New()
	u := uuidV4.String()
	duration := constant.DefaultShareSecretKeyDestroyTime
	//表示保存到数据库中限时3天

	//存的时候 k为uuid  v为对应克隆之后的chatId 需要将对应的chatId中对应的userId更换
	if err != nil {
		return constant.ZeroString, err
	}

	//异步保存记录 先返回对应sk 表示分享成功
	go setHistory(cloneChatId, u, duration)
	return u, nil
}

func setHistory(cloneChatId int, u string, duration time.Duration) {
	//TODO 补充日志
	_ = redisUtils.SetStructWithExpire(u, cloneChatId, duration)
}

// 解码对应的sk 并赋值为对应的userId
func DecodeSk(skStr string, tokenString string) error {
	cloneChatId, err := redisUtils.GetStruct[int](skStr)
	if err != nil {
		return err
	}
	userId, err := utils.DecodeToId(tokenString)
	if err != nil {
		return err
	}
	err = models.UpdateSharedHistoryUser(cloneChatId, userId)
	if err != nil {
		return err
	}
	return nil
}
