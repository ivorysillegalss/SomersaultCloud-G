package handler

import (
	"SomersaultCloud/app/bootstrap"
	task2 "SomersaultCloud/app/constant/task"
	"SomersaultCloud/app/domain"
	"SomersaultCloud/app/infrastructure/log"
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
		log.GetTextLogger().Fatal("illegal llm executor id")
	}

	log.GetJsonLogger().WithFields("choose executor", true, "executor type", executorType)
	return executor
}
