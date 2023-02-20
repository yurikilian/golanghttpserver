package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yurikilian/bills/pkg/exception"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"sync"
)

type Middleware func(next HttpMethodHandler) HttpMethodHandler
type HttpMethodHandler func(ctx IHttpContext) error
type HandlersByPath = map[string]HttpMethodHandler
type Routes map[string]HandlersByPath
type RestServerConfiguration struct {
	Address string
}

type RestServer struct {
	mux         *http.ServeMux
	router      *RestRouter
	server      *http.Server
	middlewares []func(next HttpMethodHandler) HttpMethodHandler
	binder      *Binder
	ctxPool     sync.Pool

	configuration *RestServerConfiguration
	options       *Options
}

func NewRestServer(options *Options) *RestServer {

	srv := &RestServer{
		mux:     http.NewServeMux(),
		server:  &http.Server{Addr: options.BindAddress},
		binder:  NewBinder(),
		options: options,
	}

	srv.ctxPool.New = func() interface{} {
		return NewHttpContext(nil, nil, srv.options.Log, srv.binder)
	}

	return srv
}

func (srv *RestServer) Router(router *RestRouter) *RestServer {
	srv.router = router
	return srv
}

func (srv *RestServer) Use(middleware Middleware) *RestServer {
	srv.middlewares = append(srv.middlewares, middleware)
	return srv
}

func (srv *RestServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	httpContext := srv.ctxPool.Get().(IHttpContext)
	httpContext.reset(w, req)

	handler := srv.getHandler(req)

	err := srv.applyMiddlewares(handler)(httpContext)

	if err != nil {
		e := srv.errorHandler(srv.handleError(err))(httpContext)
		if e != nil {
			println("unexpected")
		}
	}

	srv.ReleaseContext(httpContext)

}

func (srv *RestServer) process(handler HttpMethodHandler, httpContext IHttpContext) {
	err := srv.applyMiddlewares(handler)(httpContext)

	if err != nil {
		e := srv.errorHandler(srv.handleError(err))(httpContext)
		if e != nil {
			println("unexpected")
		}
	}
}

func (srv *RestServer) getHandler(req *http.Request) HttpMethodHandler {
	httpMethodHandler, status := srv.router.load(req.URL.Path, req.Method)

	if status == Matched {
		return httpMethodHandler
	} else if status == PathNotFound {
		return srv.errorHandler(exception.NewRouteNotFound(req.URL.Path))
	} else {
		return srv.errorHandler(exception.NewMethodNotAllowed(req.URL.Path, req.Method))
	}
}

func (srv *RestServer) AcquireContext() IHttpContext {
	return srv.ctxPool.Get().(IHttpContext)
}

func (srv *RestServer) ReleaseContext(httpContext IHttpContext) {
	srv.ctxPool.Put(httpContext)
}

func (srv *RestServer) errorHandler(lErr exception.Problem) func(ctx IHttpContext) error {
	return func(ctx IHttpContext) error {

		span := trace.SpanFromContext(ctx.Request().Context())
		span.RecordError(lErr)
		defer span.End()

		srv.writeException(ctx.Writer(), lErr)
		return nil
	}
}

func (srv *RestServer) applyMiddlewares(handlerByMethod HttpMethodHandler) HttpMethodHandler {
	fnc := handlerByMethod

	for i := len(srv.middlewares) - 1; i >= 0; i-- {
		fnc = srv.middlewares[i](fnc)
	}
	return fnc
}

func (srv *RestServer) Start(ctx context.Context) (exception.Problem, bool) {

	err := Validator.Validate(srv.options)
	if err != nil {
		return exception.NewInternalServerError(err.Error()), false
	}

	srv.mux.Handle("/", srv)
	srv.server.Handler = srv.mux

	srv.options.Log.Info(ctx, fmt.Sprintf("Starting server on %v address", srv.options.BindAddress))

	return srv.handleError(srv.server.ListenAndServe()), true
}

func (srv *RestServer) Shutdown(ctx context.Context) {
	err := srv.server.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}

func (srv *RestServer) handleError(err error) exception.Problem {
	switch err.(type) {
	case exception.Problem:
		return err.(exception.Problem)
	default:
		return exception.NewInternalServerError(err.Error())
	}
}

func (srv *RestServer) getHttpServer() *http.Server {
	return srv.server
}

func (srv *RestServer) writeException(w http.ResponseWriter, ex exception.Problem) {
	w.WriteHeader(ex.Code)

	marshal, err := json.Marshal(ex)
	if err != nil {
		//TODO print error
	}
	w.Write(marshal)

	///TODO log for errors unwanted
}
