package bootstrap

import (
	"SomersaultCloud/database/mongo"
	"SomersaultCloud/database/mysql"
	"SomersaultCloud/database/redis"
	"SomersaultCloud/internal/tokenutil"
)

type Application struct {
	Env       *Env
	Databases Databases
}

type Databases struct {
	Mongo mongo.Client
	Redis redis.Client
	Mysql mysql.Client
}

// 依赖注入大本营！TODO 使用wire进行改造

func App() Application {
	app := &Application{}
	app.Env = NewEnv()

	app.Databases.Mongo = NewMongoDatabase(app.Env)
	app.Databases.Redis = NewRedisDatabase(app.Env)
	app.Databases.Mysql = NewMysqlDatabase(app.Env)

	tokenutil.NewInternalApplicationConfig(app.Env)

	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Databases.Mongo)
}
