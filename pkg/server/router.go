package server

import (
	"github.com/yurikilian/bills/pkg/matcher"
	"net/http"
	"strings"
)

type LoadStatus int

const (
	PathNotFound LoadStatus = iota
	MethodNotAllowed
	Matched
)

type RestRouter struct {
	routes Routes
}

func NewRestRouter() *RestRouter {
	return &RestRouter{routes: Routes{}}
}

func (r *RestRouter) Get(path string, handlerFunc HttpMethodHandler) *RestRouter {
	r.register(path, http.MethodGet, handlerFunc)
	return r
}

func (r *RestRouter) POST(path string, handlerFunc HttpMethodHandler) *RestRouter {
	r.register(path, http.MethodPost, handlerFunc)
	return r
}

func (r *RestRouter) register(path string, httpMethod string, handlerFunc HttpMethodHandler) {
	_, pathExists := r.routes[path]
	if !pathExists {
		r.routes[path] = map[string]HttpMethodHandler{}
	}

	r.routes[path][httpMethod] = handlerFunc
}

func (r *RestRouter) load(path, method string) (HttpMethodHandler, LoadStatus) {
	handlersByPath, ok := r.matchHandlers(path)
	if !ok {
		return nil, PathNotFound
	}

	httpMethodHandler, ok := handlersByPath[method]
	if !ok {
		return nil, MethodNotAllowed
	}

	return httpMethodHandler, Matched
}

func (r *RestRouter) loadPathHandlers(path string) (HandlersByPath, bool) {
	return r.matchHandlers(path)
}

func (r *RestRouter) matchHandlers(path string) (HandlersByPath, bool) {
	pathParts := strings.Split(path, "/")
	for pattern, methods := range r.routes {

		if matcher.MatchPath(pathParts, pattern) {
			return methods, true
		}
	}

	return nil, false
}
