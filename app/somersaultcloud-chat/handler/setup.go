package handler

import (
	"SomersaultCloud/app/somersaultcloud-chat/bootstrap"
	task2 "SomersaultCloud/app/somersaultcloud-chat/constant/task"
	"SomersaultCloud/app/somersaultcloud-chat/domain"
	log2 "SomersaultCloud/app/somersaultcloud-common/log"
)

func NewLanguageModelExecutor(env *bootstrap.Env, channels *bootstrap.Channels, executorId int) domain.LanguageModelExecutor {
	var executor domain.LanguageModelExecutor
	var executorType string
	switch executorId {
	case task2.ChatAskExecutorId:
		executor = &OpenaiChatModelExecutor{env: env, res: channels}
		executorType = task2.ExecuteChatAskType
	case task2.ChatTitleAskExecutorId:
		executor = &OpenaiChatModelExecutor{env: env, res: channels}
		executorType = task2.ExecuteTitleAskType
	case task2.ChatVisionAskExecutorId:
		executor = &OpenaiVisionModelExecutor{env: env, res: channels}
		executorType = task2.ExecuteChatVisionAskType
	default:
		log2.GetTextLogger().Fatal("illegal llm executor id")
	}

	log2.GetJsonLogger().WithFields("choose executor", true, "executor type", executorType)
	return executor
}
