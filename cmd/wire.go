//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"SomersaultCloud/api/controller"
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/consume"
	"SomersaultCloud/cron"
	"SomersaultCloud/executor"
	"SomersaultCloud/internal/tokenutil"
	"SomersaultCloud/repository"
	"SomersaultCloud/task"
	"SomersaultCloud/usecase"
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

	task.NewAskChatTask,
	task.NewChatTitleTask,
	task.NewConvertTask,

	controller.NewChatController,
	controller.NewHistoryMessageController,

	wire.Struct(new(bootstrap.Application), "*"),
)

// InitializeApp init application.
func InitializeApp() (*bootstrap.Application, error) {
	wire.Build(appSet)
	return &bootstrap.Application{}, nil
}
