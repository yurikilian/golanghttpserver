package server

import (
	"context"
	"encoding/json"
	"github.com/yurikilian/bills/pkg/exception"
	"github.com/yurikilian/bills/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type Middleware func(next RestRouteHandler) RestRouteHandler
type RestRouteHandler func(ctx *HttpContext) error
type RouteMap map[string]map[string]RestRouteHandler

type RestServer struct {
	mux         *http.ServeMux
	router      *RestRouter
	server      *http.Server
	middlewares []func(next RestRouteHandler) RestRouteHandler
	log         logger.Logger
	binder      *Binder
}

func NewRestServer(options *Options) *RestServer {
	return &RestServer{
		mux:    http.NewServeMux(),
		server: &http.Server{Addr: options.BindAddress},
		log:    options.Log,
		binder: NewBinder(),
	}
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
	handler, lErr := srv.router.load(req.URL.Path, req.Method)

	if lErr != nil {
		handler = srv.errorHandler(lErr)
	}

	//TODO: check concurrency!
	httpContext := NewHttpContext(w, req, srv.log, srv.binder)

	err := srv.applyMiddlewares(handler)(httpContext)

	if err != nil {
		e := srv.errorHandler(srv.handleError(err))(httpContext)
		if e != nil {
			println("unexpected")
			return
		}
		return
	}

}

func (srv *RestServer) errorHandler(lErr *exception.Problem) func(ctx *HttpContext) error {
	return func(ctx *HttpContext) error {

		span := trace.SpanFromContext(ctx.Request().Context())
		span.RecordError(lErr)
		defer span.End()

		srv.writeException(ctx.writer, lErr)
		return nil
	}
}

func (srv *RestServer) applyMiddlewares(handlerByMethod RestRouteHandler) RestRouteHandler {
	fnc := handlerByMethod

	for i := len(srv.middlewares) - 1; i >= 0; i-- {
		fnc = srv.middlewares[i](fnc)
	}
	return fnc
}

func (srv *RestServer) Start() *exception.Problem {
	srv.mux.Handle("/", srv)
	srv.server.Handler = srv.mux
	err := srv.server.ListenAndServe()

	if err != nil {
		return srv.handleError(err)
	}
	return nil
}

func (srv *RestServer) Shutdown(ctx context.Context) {
	err := srv.server.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}

func (srv *RestServer) handleError(err error) *exception.Problem {
	switch err.(type) {
	case *exception.Problem:
		return err.(*exception.Problem)
	default:
		return exception.NewInternalServerError(err.Error())
	}
}

func (srv *RestServer) getHttpServer() *http.Server {
	return srv.server
}

func (srv *RestServer) writeException(w http.ResponseWriter, ex *exception.Problem) {
	w.WriteHeader(ex.Code)

	marshal, err := json.Marshal(ex)
	if err != nil {
		//TODO print error
	}
	w.Write(marshal)

	///TODO log for errors unwanted
}
