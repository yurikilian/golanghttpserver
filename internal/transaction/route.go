package transaction

import (
	"github.com/yurikilian/bills/pkg/exception"
	"github.com/yurikilian/bills/pkg/server"
)

type Route struct {
	service *Service
}

func (r *Route) Find(ctx server.IHttpContext) error {

	ctx.Logger().Debug(ctx.ReqCtx(), "Entering find function")

	var request FindRequest
	if bErr := ctx.ReadBody(&request); bErr != nil {
		return bErr
	}

	trn, err := r.service.Find(request.id)
	if err != nil {
		ctx.Logger().Error(ctx.ReqCtx(), err.Error())
		return exception.NewInternalServerError(err.Error())
	}
	return ctx.WriteResponse(200, trn)
}

func (r *Route) Create(ctx server.IHttpContext) error {
	var request CreationRequest
	if bErr := ctx.ReadBody(&request); bErr != nil {
		return bErr
	}

	if err := r.service.Create(request); err != nil {
		ctx.Logger().Debug(ctx.ReqCtx(), err.Error())
		return exception.NewInternalServerError(err.Error())
	}
	return ctx.WriteResponse(204, nil)
}
