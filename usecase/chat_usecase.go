package usecase

import (
	"SomersaultCloud/domain"
	"SomersaultCloud/repository"
	"context"
)

type chatUseCase struct {
	chatRepository domain.ChatRepository
}

func NewChatUseCase() domain.ChatUseCase {
	return &chatUseCase{chatRepository: repository.NewChatRepository()}
}

func (c chatUseCase) InitChat(ctx context.Context) int {
	id := c.chatRepository.CacheGetNewestChatId(ctx)
	c.chatRepository.CacheInsertNewChat(ctx, id)
	//	TODO 线程不安全 lua保证原子性
	return id + 1
}
