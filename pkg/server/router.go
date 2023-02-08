package server

import (
	"github.com/yurikilian/bills/pkg/exception"
	"github.com/yurikilian/bills/pkg/matcher"
	"net/http"
	"strings"
)

type RestRouter struct {
	routes Routes
}

func NewRestRouter() *RestRouter {
	return &RestRouter{routes: Routes{}}
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
	_, pathExists := r.routes[path]
	if !pathExists {
		r.routes[path] = map[string]RestRouteHandler{}
	}

	r.routes[path][httpMethod] = handlerFunc
}

func (r *RestRouter) load(path string, httpMethod string) (RestRouteHandler, *exception.Problem) {
	handlersByPath, pathExists := r.matchHandlers(path)

	if !pathExists {
		return nil, exception.NewRouteNotFound(path)
	}

	handlerByMethod, methodMapExists := handlersByPath[httpMethod]
	if !methodMapExists {
		return nil, exception.NewMethodNotAllowed(path, httpMethod)

	}
	return handlerByMethod, nil
}

func (r *RestRouter) matchHandlers(path string) (HandlersByMethod, bool) {
	pathParts := strings.Split(path, "/")
	for pattern, methods := range r.routes {

		if matcher.MatchPath(pathParts, pattern) {
			return methods, true
		}
	}

	return nil, false
}
