package main

import (
	route "SomersaultCloud/api/route"
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/cron"
	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()

	env := app.Env

	//db := app.Databases
	defer app.CloseDBConnection()

	//timeout := time.Duration(env.ContextTimeout) * time.Second

	gin := gin.Default()

	route.Setup(env, gin)
	cron.Setup()

	gin.Run(env.ServerAddress)
}
