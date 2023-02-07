package server

import (
	"github.com/yurikilian/bills/pkg/exception"
	"net/http"
	"reflect"
	"testing"
)

var (
	emptyHandlerFunc = func(ctx *HttpContext) error { return nil }
)

func TestNewRestRouter(t *testing.T) {
	tests := []struct {
		name string
		want *RestRouter
	}{
		{name: "Should create a new rest router", want: NewRestRouter()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRestRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRestRouter() = %v, expectedHandler %v", got, tt.want)
			}
		})
	}
}

func TestRestRouter_Get(t *testing.T) {
	type fields struct {
		routeMap RouteMap
	}
	type args struct {
		path        string
		handlerFunc RestRouteHandler
	}

	routeMap := RouteMap{"/transactions": map[string]RestRouteHandler{"GET": emptyHandlerFunc}}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RestRouter
	}{
		{
			name: "Should return the router on GET call", fields: fields{routeMap: routeMap}, args: args{
				path:        "/transactions",
				handlerFunc: emptyHandlerFunc,
			},
			want: &RestRouter{routeMap: routeMap},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestRouter{
				routeMap: tt.fields.routeMap,
			}
			if got := r.Get(tt.args.path, tt.args.handlerFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, expectedHandler %v", got, tt.want)
			}
		})
	}
}

func TestRestRouter_load(t *testing.T) {

	routeMap := RouteMap{"/transactions": map[string]RestRouteHandler{"GET": emptyHandlerFunc}}
	type fields struct {
		routeMap RouteMap
	}
	type args struct {
		path       string
		httpMethod string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		expectedHandler RestRouteHandler
		expectedEx      *exception.Problem
	}{
		{
			name:   "Should return not found given not registered path",
			fields: fields{routeMap: routeMap},
			args: args{
				path:       "/invalid",
				httpMethod: "GET",
			},
			expectedHandler: nil,
			expectedEx:      exception.NewRouteNotFound("/invalid"),
		},
		{
			name:   "Should return method not allowed exception given invalid method registerd",
			fields: fields{routeMap: routeMap},
			args: args{
				path:       "/transactions",
				httpMethod: http.MethodPut,
			},
			expectedHandler: nil,
			expectedEx:      exception.NewMethodNotAllowed("/transactions", http.MethodPut),
		},
		{
			name:   "Should return handler given valid path and valid method",
			fields: fields{routeMap: routeMap},
			args: args{
				path:       "/transactions",
				httpMethod: http.MethodGet,
			},
			expectedHandler: emptyHandlerFunc,
			expectedEx:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestRouter{
				routeMap: tt.fields.routeMap,
			}
			got, got1 := r.load(tt.args.path, tt.args.httpMethod)
			if got != nil && tt.expectedHandler != nil {
				if !reflect.DeepEqual(got(nil), tt.expectedHandler(nil)) {
					t.Errorf("load() got = %v, expectedHandler %v", got, tt.expectedHandler)
				}
			}
			if !reflect.DeepEqual(got1, tt.expectedEx) {
				t.Errorf("load() got1 = %v, expectedHandler %v", got1, tt.expectedEx)
			}
		})
	}
}

func TestRestRouter_register(t *testing.T) {

	type fields struct {
		routeMap RouteMap
	}
	type args struct {
		path        string
		httpMethod  string
		handlerFunc RestRouteHandler
		expected    RestRouteHandler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "Should register to map given path and method",
			fields: fields{RouteMap{}},
			args: args{
				path:        "/transactions",
				httpMethod:  http.MethodPatch,
				handlerFunc: emptyHandlerFunc,
				expected:    emptyHandlerFunc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestRouter{
				routeMap: tt.fields.routeMap,
			}
			r.register(tt.args.path, tt.args.httpMethod, tt.args.handlerFunc)
			reflect.DeepEqual(r.routeMap[tt.args.path][tt.args.httpMethod](nil), tt.args.expected(nil))
		})
	}
}
