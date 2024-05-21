package api

import (
	"mini-gpt/constant"
	"mini-gpt/models"
	"mini-gpt/prompt"
)

//此类本意是想处理所有类中响应信息的 目前只有openai一个接口
//暂时只有这一接口的方法

func ConcludeTitle(msg *[]models.Message) (string, error) {

	//拼接提示词 & 聊天记录
	messages := *msg
	titleSystemPromptMessage := models.Message{
		Role: constant.SystemRole,
		//从存储提示词的map中取出对应的提示词
		Content: prompt.OpenaiPrompt[constant.Conclude2TitlePrompt],
	}

	//处理prompt格式
	titleInputPrompt := "<" + messages[0].Role + ":" + messages[0].Content +
		messages[1].Role + ":" + messages[1].Content + ">"

	titleInputPromptMessage := models.Message{
		Role:    constant.UserRole,
		Content: titleInputPrompt,
	}

	var titlePromptMessage []models.Message
	titlePromptMessage = append(titlePromptMessage, titleSystemPromptMessage, titleInputPromptMessage)

	//拼接提问模型
	history := &models.ChatCompletionRequest{
		Messages: titlePromptMessage,
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
	completionMessage := simplyChatCompletionMessage(response)
	return completionMessage.GenerateText, nil
}

// 将instruct模型调用api结果包装为返回用户的结果
func simplyTextCompletionMessage(completionResponse *models.TextCompletionResponse) *models.GenerateMessage {
	//openAI返回的json中请求体中的文本是一个数组 暂取第0项
	args := completionResponse.Choices
	if args == nil {
		return models.ErrorGeneration()
	}
	textBody := args[0]
	generateMessage := models.GenerateMessage{
		GenerateText: textBody.Text,
		FinishReason: textBody.FinishReason,
	}
	return &generateMessage
}

// 转换chat模型调用结果
func simplyChatCompletionMessage(completionResponse *models.ChatCompletionResponse) *models.GenerateMessage {
	//openAI返回的json中请求体中的文本是一个数组 暂取第0项
	args := completionResponse.Choices
	if args == nil {
		return models.ErrorGeneration()
	}
	textBody := args[0]
	generateMessage := models.GenerateMessage{
		GenerateText: textBody.Message.Content,
		FinishReason: textBody.FinishReason,
	}
	return &generateMessage
}

// 简化包装信息
func SimplyMessage(completionResponse models.BaseModel, modelType string) *models.GenerateMessage {
	var generationMessage *models.GenerateMessage
	//修改为类型绑定 TODO
	if modelType == constant.InstructModel {
		if textCompletionResponse, ok := completionResponse.(*models.TextCompletionResponse); ok {
			generationMessage = simplyTextCompletionMessage(textCompletionResponse)
		}
		//这里可以拓展 设计模式优化
	} else if modelType == constant.DefaultModel {
		if chatCompletionResponse, ok := completionResponse.(*models.ChatCompletionResponse); ok {
			generationMessage = simplyChatCompletionMessage(chatCompletionResponse)
		}
	}
	return generationMessage
}
