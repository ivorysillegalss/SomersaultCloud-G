package api

import (
	"mini-gpt/constant"
	"mini-gpt/models"
	"mini-gpt/service"
)

func ConcludeTitle(msg *[]models.Message) (string, error) {

	//拼接提示词 & 聊天记录
	messages := *msg
	prompt4Title := models.Message{
		Role:    constant.UserRole,
		Content: constant.Conclude2TitlePrompt,
	}
	messages = append(messages, prompt4Title)

	//拼接提问模型
	history := &models.ChatCompletionRequest{
		Messages: messages,
		CompletionRequest: models.CompletionRequest{
			MaxTokens: constant.DefaultMaxToken,
			Model:     constant.DefaultModel,
		},
	}

	//执行获取标题内容
	baseModel, err := Execute(constant.DefaultAdminUID, history, constant.DefaultModel)
	if err != nil {
		return constant.ZeroString, err
	}

	response := baseModel.(*models.ChatCompletionResponse)
	completionMessage := service.SimplyChatCompletionMessage(response)
	return completionMessage.GenerateText, nil
}
