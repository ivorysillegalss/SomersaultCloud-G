package bootstrap

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RunIpConfig(app *IpConfigApplication) {
	app.Dispatcher.Handle()
	app.DataHandler.Handle()
	s := server.Default(server.WithHostPorts(":6789"))
	s.GET("/ip/list", app.Api.GetInfoList)
	s.Spin()
}
