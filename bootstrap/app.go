package bootstrap

import (
	"SomersaultCloud/cron"
	"SomersaultCloud/domain"
	"SomersaultCloud/handler"
	"SomersaultCloud/infrastructure/mongo"
	"SomersaultCloud/infrastructure/mysql"
	"SomersaultCloud/infrastructure/pool"
	"SomersaultCloud/infrastructure/redis"
	"SomersaultCloud/internal/tokenutil"
	"SomersaultCloud/task"
	"SomersaultCloud/usecase"
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
	RpcRes chan *domain.GenerationResponse
	Stop   chan bool
}

// 依赖注入大本营！TODO 使用wire进行改造

func App() Application {
	app := &Application{}
	app.Env = NewEnv()

	app.Databases.Mongo = NewMongoDatabase(app.Env)
	app.Databases.Redis = NewRedisDatabase(app.Env)
	app.Databases.Mysql = NewMysqlDatabase(app.Env)
	app.PoolsFactory.Pools = NewPoolFactory()
	app.Channels = NewChannel()

	tokenutil.NewInternalApplicationConfig(app.Env)
	usecase.NewUseCaseApplicationConfig(app.Env)
	task.NewUseCaseApplicationConfig(app.Env, app.PoolsFactory)
	handler.NewUseCaseApplicationConfig(app.Env, app.Channels)
	cron.NewCronApplicationConfig(app.Channels)

	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Databases.Mongo)
}
