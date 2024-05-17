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
	chat, err := models.GetChatInfo(chatId)
	chat.Records = chatHistory
	if err != nil {
		return constant.ZeroString, err
	}

	uuidV4 := uuid.New()
	u := uuidV4.String()
	duration := constant.DefaultShareSecretKeyDestroyTime
	//表示保存到数据库中限时3天

	//异步保存到Redis中 方便下方获取并写入表中 先返回对应sk 表示分享成功
	go setHistory(u, chat, duration)
	return u, nil
}

func setHistory(u string, chat *models.Chat, duration time.Duration) {
	//存到redis中方便想要取的对方
	//存进去的时候 A为Record指针 B为Chat指针
	//_ = redisUtils.SetStructWithExpire(u, *tuple.NewTuple(*history, chat), duration)
	_ = redisUtils.SetStructWithExpire(u, chat, duration)
}

// 解码对应的sk 并赋值为对应的userId
func DecodeSk(skStr string) (*models.Chat, error) {
	chatValue, err := redisUtils.GetStruct[*models.Chat](skStr)
	if err != nil {
		return nil, err
	}
	return chatValue, nil
}

// 依据分享的chat继续聊下去
// 从redis取出原历史记录 + 制作备份（更新对应的userId）+ 返回备份后的新chatId
func UpdateSharedChat(tokenString string, skStr string) (int, error) {
	userId, err := utils.DecodeToId(tokenString)
	chatValue, err := redisUtils.GetStruct[*models.Chat](skStr)
	if err != nil {
		return constant.FalseInt, err
	}
	//此处返回备份后的新chatId
	cloneChatId, err := copyChatHistory(userId, *chatValue.Records)
	if err != nil {
		return constant.FalseInt, err
	}
	return cloneChatId, nil
}

// 复制对应的一份历史记录 获取分享的历史记录的时候可以用到
func copyChatHistory(userId int, chatHistory []*models.Record) (int, error) {
	//从第0项获取这一chat的botId
	cloneChatId, err := models.CreateNewChat(userId, chatHistory[0].ChatAsks.BotId)
	if err != nil {
		return constant.FalseInt, err
	}

	//赋值聊天记录 除了chatId和userId以外其他都是和原来一样的
	count := len(chatHistory)
	for i := 0; i < count; i++ {

		//重新赋值问和答的chatID
		chatHistory[i].ChatAsks.ChatId = cloneChatId
		chatHistory[i].ChatGenerations.ChatId = cloneChatId

		err := models.SaveRecord(chatHistory[i], cloneChatId)
		if err != nil {
			return constant.FalseInt, err
		}
	}
	return cloneChatId, nil
}
