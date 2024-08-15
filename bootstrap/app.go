package bootstrap

import (
	"SomersaultCloud/handler"
	"SomersaultCloud/infrastructure/channel"
	"SomersaultCloud/infrastructure/mongo"
	"SomersaultCloud/infrastructure/mysql"
	"SomersaultCloud/infrastructure/pool"
	"SomersaultCloud/infrastructure/redis"
	"SomersaultCloud/internal/tokenutil"
	"SomersaultCloud/usecase"
	"SomersaultCloud/usecase/task"
)

type Application struct {
	Env          *Env
	Databases    Databases
	PoolsFactory *PoolsFactory
	Channels     *Channels
}

type Databases struct {
	Mongo mongo.Client
	Redis redis.Client
	Mysql mysql.Client
}

// PoolsFactory k为pool业务号 v为poll详细配置信息
type PoolsFactory struct {
	Pools map[int]pool.Pool
}

type Channels struct {
	RpcRes      chan *channel.GenerationResponse
	asyncPoller chan *channel.GenerationResponse
}

// 依赖注入大本营！TODO 使用wire进行改造

func App() Application {
	app := &Application{}
	app.Env = NewEnv()

	app.Databases.Mongo = NewMongoDatabase(app.Env)
	app.Databases.Redis = NewRedisDatabase(app.Env)
	app.Databases.Mysql = NewMysqlDatabase(app.Env)
	app.PoolsFactory.Pools = NewPoolFactory()
	//TODO 可修改为定时任务模块
	//	这里耦合channel中的异步任务需使用到redis 耦合了
	app.Channels = NewChannel(app.Databases.Redis)

	tokenutil.NewInternalApplicationConfig(app.Env)
	usecase.NewUseCaseApplicationConfig(app.Env)
	task.NewUseCaseApplicationConfig(app.Env, app.PoolsFactory)
	handler.NewUseCaseApplicationConfig(app.Env, app.Channels)

	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Databases.Mongo)
}
