package bootstrap

import (
	"SomersaultCloud/app/somersaultcloud-ipconfig/dispatcher"
	"SomersaultCloud/app/somersaultcloud-ipconfig/source"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type Api struct {
	DataHandler *source.DataHandler
	Dispatcher  *dispatcher.Dispatcher
}

func (a *Api) GetInfoList(ctx context.Context, ctx2 *app.RequestContext) {

}

func NewApi(h *source.DataHandler, d *dispatcher.Dispatcher) *Api {
	return &Api{DataHandler: h, Dispatcher: d}
}
