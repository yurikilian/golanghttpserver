package middleware

import (
	"github.com/yurikilian/bills/pkg/exception"
	"github.com/yurikilian/bills/pkg/server"
	"mime"
)

func Json() server.Middleware {
	return func(next server.RestRouteHandler) server.RestRouteHandler {
		return func(ctx *server.HttpContext) error {
			contentType := ctx.Request().Header.Get("Content-Type")
			mediaType, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				return exception.NewUnsupportedMediaType("Invalid Content-type")
			}

			if mediaType != "application/json" {
				return exception.NewBadRequestProblem("Content-Type header must be application/json")
			}

			return next(ctx)
		}
	}
}
