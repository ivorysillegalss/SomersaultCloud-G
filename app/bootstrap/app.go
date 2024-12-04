package bootstrap

import (
	controller2 "SomersaultCloud/app/api/controller"
	"SomersaultCloud/app/domain"
	"SomersaultCloud/app/executor"
	"SomersaultCloud/app/infrastructure/mongo"
	"SomersaultCloud/app/infrastructure/mysql"
	"SomersaultCloud/app/infrastructure/pool"
	"SomersaultCloud/app/infrastructure/redis"
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
	RpcRes       chan *domain.GenerationResponse
	Stop         chan bool
	StreamRpcRes chan *domain.GenerationResponse
}

type Controllers struct {
	ChatController           *controller2.ChatController
	HistoryMessageController *controller2.HistoryMessageController
}

type Executor struct {
	CronExecutor    *executor.CronExecutor
	ConsumeExecutor *executor.ConsumeExecutor
	DataExecutor    *executor.DataExecutor
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Databases.Mongo)
}

func NewControllers(chatController *controller2.ChatController, messageController *controller2.HistoryMessageController) *Controllers {
	return &Controllers{ChatController: chatController, HistoryMessageController: messageController}
}

func NewExecutors(ce *executor.CronExecutor, cse *executor.ConsumeExecutor, de *executor.DataExecutor) *Executor {
	return &Executor{
		CronExecutor:    ce,
		ConsumeExecutor: cse,
		DataExecutor:    de,
	}
}
