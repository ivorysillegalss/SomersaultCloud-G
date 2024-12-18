// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	controller2 "SomersaultCloud/app/somersaultcloud-chat/api/controller"
	"SomersaultCloud/app/somersaultcloud-chat/bootstrap"
	"SomersaultCloud/app/somersaultcloud-chat/consume"
	"SomersaultCloud/app/somersaultcloud-chat/cron"
	"SomersaultCloud/app/somersaultcloud-chat/executor"
	"SomersaultCloud/app/somersaultcloud-chat/internal/tokenutil"
	"SomersaultCloud/app/somersaultcloud-chat/repository"
	task2 "SomersaultCloud/app/somersaultcloud-chat/task"
	"SomersaultCloud/app/somersaultcloud-chat/usecase"
	"github.com/google/wire"
)

// Injectors from wire.go:

// InitializeApp init application.
func InitializeApp() (*bootstrap.Application, error) {
	env := bootstrap.NewEnv()
	databases := bootstrap.NewDatabases(env)
	poolsFactory := bootstrap.NewPoolFactory()
	channels := bootstrap.NewChannel()
	chatRepository := repository.NewChatRepository(databases, env)
	botRepository := repository.NewBotRepository(databases)
	connection := bootstrap.NewRabbitConnection(env)
	messageHandler := consume.NewMessageHandler(connection)
	storageEvent := consume.NewStorageEvent(chatRepository, messageHandler)
	askTask := task2.NewAskChatTask(botRepository, chatRepository, env, channels, poolsFactory, storageEvent)
	tokenUtil := tokenutil.NewTokenUtil(env)
	titleTask := task2.NewChatTitleTask(chatRepository, env, channels)
	generationRepository := repository.NewGenerationRepository(databases)
	generateEvent := consume.NewGenerateEvent(messageHandler, env, channels, generationRepository)
	convertTask := task2.NewConvertTask(generateEvent, generationRepository)
	chatUseCase := usecase.NewChatUseCase(env, chatRepository, botRepository, askTask, tokenUtil, storageEvent, titleTask, convertTask, generationRepository)
	chatController := controller2.NewChatController(chatUseCase)
	historyMessageController := controller2.NewHistoryMessageController(chatUseCase)
	controllers := bootstrap.NewControllers(chatController, historyMessageController)
	generationCron := cron.NewGenerationCron(generationRepository, channels, env, generateEvent)
	cronExecutor := executor.NewCronExecutor(generationCron)
	consumeExecutor := executor.NewConsumeExecutor(storageEvent, generateEvent)
	dataExecutor := executor.NewDataExecutor(databases.Redis)
	bootstrapExecutor := bootstrap.NewExecutors(cronExecutor, consumeExecutor, dataExecutor)
	application := &bootstrap.Application{
		Env:          env,
		Databases:    databases,
		PoolsFactory: poolsFactory,
		Channels:     channels,
		Controllers:  controllers,
		Executor:     bootstrapExecutor,
	}
	return application, nil
}

// wire.go:

var appSet = wire.NewSet(bootstrap.NewEnv, tokenutil.NewTokenUtil, bootstrap.NewDatabases, bootstrap.NewRedisDatabase, bootstrap.NewMysqlDatabase, bootstrap.NewMongoDatabase, bootstrap.NewPoolFactory, bootstrap.NewChannel, bootstrap.NewRabbitConnection, bootstrap.NewControllers, bootstrap.NewExecutors, repository.NewGenerationRepository, repository.NewChatRepository, repository.NewBotRepository, consume.NewStorageEvent, consume.NewGenerateEvent, consume.NewMessageHandler, cron.NewGenerationCron, executor.NewCronExecutor, executor.NewConsumeExecutor, usecase.NewChatUseCase, task2.NewAskChatTask, task2.NewChatTitleTask, task2.NewConvertTask, controller2.NewChatController, controller2.NewHistoryMessageController, wire.Struct(new(bootstrap.Application), "*"))
