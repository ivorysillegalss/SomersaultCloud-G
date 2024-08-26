package main

import (
	route "SomersaultCloud/api/route"
)

func main() {
	app, err := InitializeApp()
	if err != nil {
		return
	}
	defer app.CloseDBConnection()

	setup := route.Setup(app.Controllers, app.Executor)
	//gin.SetMode(gin.DebugMode)
	setup.Run(app.Env.ServerAddress)
}
