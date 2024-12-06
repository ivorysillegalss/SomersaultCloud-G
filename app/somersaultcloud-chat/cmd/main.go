package main

import (
	"SomersaultCloud/app/somersaultcloud-chat/api/grpc"
	"SomersaultCloud/app/somersaultcloud-chat/api/route"
)

func main() {
	app, err := InitializeApp()
	if err != nil {
		return
	}
	defer app.CloseDBConnection()

	setup := route.Setup(app.Controllers, app.Executor)
	grpc.Setup(app.Env.Grpc)
	//gin.SetMode(gin.DebugMode)
	setup.Run(app.Env.ServerAddress)
}
