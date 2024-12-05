package bootstrap

import (
	"SomersaultCloud/app/somersaultcloud-ipconfig/domain"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Api struct {
	Dispatcher domain.Dispatcher
}

func (a *Api) GetInfoList(ctx context.Context, rectx *app.RequestContext) {
	defer func() {
		if err := recover(); err != nil {
			rectx.JSON(consts.StatusBadRequest, utils.H{"err": err})
		}
	}()
	ipConfCtx := domain.BuildIpConfContext(ctx, rectx)

	eds := a.Dispatcher.Do(ipConfCtx)
	ipConfCtx.AppCtx.JSON(consts.StatusOK, domain.SuccessResp(eds))
}

func NewApi(d domain.Dispatcher) *Api {
	return &Api{Dispatcher: d}
}
