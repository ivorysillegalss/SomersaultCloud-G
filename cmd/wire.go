//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"SomersaultCloud/api/controller"
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/cron"
	"SomersaultCloud/internal/tokenutil"
	"SomersaultCloud/repository"
	"SomersaultCloud/task"
	"SomersaultCloud/usecase"
	"github.com/google/wire"
)

var appSet = wire.NewSet(
	bootstrap.NewEnv,
	bootstrap.NewDatabases,
	bootstrap.NewRedisDatabase,
	bootstrap.NewMysqlDatabase,
	bootstrap.NewMongoDatabase,
	bootstrap.NewPoolFactory,
	bootstrap.NewChannel,
	bootstrap.NewControllers,

	repository.NewGenerationRepository,
	repository.NewChatRepository,
	repository.NewBotRepository,

	cron.NewExecutor,
	cron.NewAsyncService,

	usecase.NewChatUseCase,

	task.NewAskChatTask,

	controller.NewChatController,

	tokenutil.NewTokenUtil,

	wire.Struct(new(bootstrap.Application), "*"),
)

// InitializeApp init application.
func InitializeApp() (*bootstrap.Application, error) {
	wire.Build(appSet)
	return &bootstrap.Application{}, nil
}