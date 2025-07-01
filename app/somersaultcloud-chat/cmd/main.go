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

	setup := route.Setup(app.Env, app.Controllers, app.Executor, app.Databases.Redis)
	grpc.Setup(app.Env.Grpc)
	setup.Run()
}
