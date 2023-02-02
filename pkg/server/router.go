package server

import (
	"github.com/yurikilian/bills/pkg/exception"
	"net/http"
)

type RestRouter struct {
	routeMap RouteMap
}

func NewRestRouter() *RestRouter {
	return &RestRouter{routeMap: RouteMap{}}
}

func (r *RestRouter) Get(path string, handlerFunc RestRouteHandler) *RestRouter {
	r.register(path, http.MethodGet, handlerFunc)
	return r
}

func (r *RestRouter) POST(path string, handlerFunc RestRouteHandler) *RestRouter {
	r.register(path, http.MethodPost, handlerFunc)
	return r
}

func (r *RestRouter) register(path string, httpMethod string, handlerFunc RestRouteHandler) {
	_, pathExists := r.routeMap[path]
	if !pathExists {
		r.routeMap[path] = map[string]RestRouteHandler{}
	}

	r.routeMap[path][httpMethod] = handlerFunc
}

func (r *RestRouter) load(path string, httpMethod string) (RestRouteHandler, *exception.Problem) {
	handlersByPath, pathExists := r.routeMap[path]

	if !pathExists {
		return nil, exception.NewRouteNotFound(path)
	}

	handlerByMethod, methodMapExists := handlersByPath[httpMethod]
	if !methodMapExists {
		return nil, exception.NewMethodNotAllowed(path, httpMethod)

	}
	return handlerByMethod, nil
}
