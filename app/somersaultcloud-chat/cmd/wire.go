//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
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

var appSet = wire.NewSet(
	bootstrap.NewEnv,
	tokenutil.NewTokenUtil,
	bootstrap.NewDatabases,
	bootstrap.NewRedisDatabase,
	bootstrap.NewMysqlDatabase,
	bootstrap.NewMongoDatabase,
	bootstrap.NewPoolFactory,
	bootstrap.NewChannel,
	bootstrap.NewRabbitConnection,
	bootstrap.NewControllers,
	bootstrap.NewExecutors,

	repository.NewGenerationRepository,
	repository.NewChatRepository,
	repository.NewBotRepository,

	consume.NewStorageEvent,
	consume.NewGenerateEvent,
	consume.NewMessageHandler,

	cron.NewGenerationCron,

	executor.NewCronExecutor,
	executor.NewConsumeExecutor,
	executor.NewDataExecutor,

	usecase.NewChatUseCase,

	task2.NewAskChatTask,
	task2.NewChatTitleTask,
	task2.NewConvertTask,

	controller2.NewChatController,
	controller2.NewHistoryMessageController,

	wire.Struct(new(bootstrap.Application), "*"),
)

// InitializeApp init application.
func InitializeApp() (*bootstrap.Application, error) {
	wire.Build(appSet)
	return &bootstrap.Application{}, nil
}
