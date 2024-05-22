package service

import (
	"mini-gpt/api"
	"mini-gpt/dto"
)

func GetTitle(getTitle *dto.TitleDTO) (*dto.TitleDTO, error) {
	concludedTitle, err := api.ConcludeTitle(&getTitle.Messages)
	if err != nil {
		return nil, err
	}
	//顺便直接将数据库对应的标题修改了
	go api.AsyncUpdateTitle(getTitle.ChatId, concludedTitle)
	if err != nil {
		return nil, err
	}
	return &dto.TitleDTO{Title: concludedTitle}, nil

}
