package service

import (
	"mini-gpt/api"
	"mini-gpt/dto"
	"mini-gpt/models"
)

func GetTitle(getTitle *dto.TitleDTO) (*dto.TitleDTO, error) {
	concludedTitle, err := api.ConcludeTitle(&getTitle.Messages)
	if err != nil {
		return nil, err
	}
	//顺便直接将数据库对应的标题修改了
	err = models.UpdateChatTitle(getTitle.ChatId, concludedTitle)
	if err != nil {
		return nil, err
	}
	return &dto.TitleDTO{Title: concludedTitle}, nil

}
