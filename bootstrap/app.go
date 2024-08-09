package bootstrap

import (
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/database/mongo"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/database/mysql"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/database/redis"
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

func App() Application {
	app := &Application{}
	app.Env = NewEnv()

	app.Databases.Mongo = NewMongoDatabase(app.Env)
	app.Databases.Redis = NewRedisDatabase(app.Env)
	app.Databases.Mysql = NewMysqlDatabase(app.Env)

	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Databases.Mongo)
}
