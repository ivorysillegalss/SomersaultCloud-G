package bootstrap

import (
	"SomersaultCloud/api/controller"
	"SomersaultCloud/domain"
	"SomersaultCloud/executor"
	"SomersaultCloud/infrastructure/mongo"
	"SomersaultCloud/infrastructure/mysql"
	"SomersaultCloud/infrastructure/pool"
	"SomersaultCloud/infrastructure/redis"
)

type Application struct {
	Env          *Env
	Databases    *Databases
	PoolsFactory *PoolsFactory
	Channels     *Channels
	Controllers  *Controllers
	Executor     *Executor
}

type Databases struct {
	Mongo mongo.Client
	Redis redis.Client
	Mysql mysql.Client
}

// PoolsFactory k为pool业务号 v为poll详细配置信息
type PoolsFactory struct {
	Pools map[int]*pool.Pool
}

type Channels struct {
	RpcRes chan *domain.GenerationResponse
	Stop   chan bool
}

type Controllers struct {
	ChatController           *controller.ChatController
	HistoryMessageController *controller.HistoryMessageController
}

type Executor struct {
	CronExecutor    *executor.CronExecutor
	ConsumeExecutor *executor.ConsumeExecutor
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Databases.Mongo)
}

func NewControllers(chatController *controller.ChatController, messageController *controller.HistoryMessageController) *Controllers {
	return &Controllers{ChatController: chatController, HistoryMessageController: messageController}
}

func NewExecutors(ce *executor.CronExecutor, cse *executor.ConsumeExecutor) *Executor {
	return &Executor{
		CronExecutor:    ce,
		ConsumeExecutor: cse,
	}
}
