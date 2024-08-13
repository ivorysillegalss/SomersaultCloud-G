package usecase

import (
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/constant/common"
	"SomersaultCloud/domain"
	"SomersaultCloud/internal/ioutil"
	"SomersaultCloud/internal/tokenutil"
	"SomersaultCloud/repository"
	"context"
	"time"
)

type chatUseCase struct {
	chatRepository domain.ChatRepository
	chatTask       domain.ChatTask
}

func NewChatUseCase() domain.ChatUseCase {
	return &chatUseCase{chatRepository: repository.NewChatRepository()}
}

func (c *chatUseCase) InitChat(ctx context.Context, token string, botId int) int {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(env.ContextTimeout))
	defer cancel()

	script, err := ioutil.LoadLuaScript("lua/increment.lua")
	if err != nil {
		return common.FalseInt
	}

	chatId, err := c.chatRepository.CacheLuaInsertNewChatId(ctx, script, cache.NewestChatIdKey)
	if err != nil {
		return common.FalseInt
	}

	id, err := tokenutil.DecodeToId(token)
	if err != nil {
		return common.FalseInt
	}
	go c.chatRepository.DbInsertNewChatId(ctx, id, botId)
	// TODO mq异步写入MYSQL

	return chatId
}
