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

	setup := route.Setup(app.Env.ServerAddress, app.Controllers, app.Executor)
	grpc.Setup(app.Env.Grpc)
	setup.Run()
}
