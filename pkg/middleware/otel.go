package middleware

import (
	"github.com/yurikilian/bills/pkg/server"
	"go.opentelemetry.io/otel"
)

func Otel() server.Middleware {
	return func(next server.HttpMethodHandler) server.HttpMethodHandler {
		return func(rCtx server.IHttpContext) error {
			traceCtx, span := otel.Tracer("").
				Start(rCtx.Request().Context(), rCtx.Request().URL.Path)

			newReq := rCtx.Request().WithContext(traceCtx)
			rCtx.SetRequest(newReq)

			// Should not end span giving error. The error handler should disable it
			err := next(rCtx)
			if err == nil {
				defer span.End()
			}
			return err
		}
	}
}
