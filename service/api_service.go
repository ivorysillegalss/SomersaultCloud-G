package service

import (
	"mini-gpt/api"
	"mini-gpt/models"
	"mini-gpt/setting"
)

//var logger = setting.GetLogger()
//正常来说应该全局变量的 但是由于代码的先后执行问题先放到下面的函数中

func LoadingChat(apiRequestMessage models.ApiRequestMessage) (models.GenerateMessage, error) {

	var logger = setting.GetLogger()

	completionResponse, err := api.Execute(apiRequestMessage)
	if err != nil {
		logger.Error(err)
		return models.GenerateMessage{}, err
	}

	//openAI返回的json中请求体中的文本是一个数组 暂取第0项
	args := completionResponse.Choices
	textBody := args[0]
	generateMessage := models.GenerateMessage{
		GenerateText: textBody.Text,
		FinishReason: textBody.FinishReason,
	}
	return generateMessage, nil
}
