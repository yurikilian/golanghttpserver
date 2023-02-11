package server

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

var (
	emptyHandlerFunc = func(ctx IHttpContext) error { return nil }
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
		routeMap Routes
	}
	type args struct {
		path        string
		handlerFunc HttpMethodHandler
	}

	routeMap := Routes{"/transactions": map[string]HttpMethodHandler{"GET": emptyHandlerFunc}}
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
			want: &RestRouter{routes: routeMap},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestRouter{
				routes: tt.fields.routeMap,
			}
			if got := r.Get(tt.args.path, tt.args.handlerFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, expectedHandler %v", got, tt.want)
			}
		})
	}
}

func TestRestRouter_load(t *testing.T) {

	routeMap := Routes{
		"/transactions":                        map[string]HttpMethodHandler{"GET": emptyHandlerFunc},
		"/transactions/:id/product/:productId": map[string]HttpMethodHandler{"GET": trnProductWithIdFunc},
	}
	type fields struct {
		routeMap Routes
	}
	type args struct {
		path       string
		httpMethod string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		expectedHandler HttpMethodHandler
		expectedStatus  LoadStatus
	}{
		{
			name:   "Should return not found given not registered path",
			fields: fields{routeMap: routeMap},
			args: args{
				path:       "/invalid",
				httpMethod: "GET",
			},
			expectedHandler: nil,
			expectedStatus:  PathNotFound,
		},
		{
			name:   "Should return method not allowed exception given invalid method registered",
			fields: fields{routeMap: routeMap},
			args: args{
				path:       "/transactions",
				httpMethod: http.MethodPut,
			},
			expectedHandler: nil,
			expectedStatus:  MethodNotAllowed,
		},
		{
			name:   "Should return handler given valid path and valid method",
			fields: fields{routeMap: routeMap},
			args: args{
				path:       "/transactions",
				httpMethod: http.MethodGet,
			},
			expectedStatus: Matched,
		},
		{
			name:   "Should match url /transactions/:id/product/:productId - URL with Path variables",
			fields: fields{routeMap: routeMap},
			args: args{
				path:       "/transactions/1/product/1",
				httpMethod: http.MethodGet,
			},
			expectedStatus: Matched,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestRouter{
				routes: tt.fields.routeMap,
			}
			handler, status := r.load(tt.args.path, tt.args.httpMethod)
			reflect.DeepEqual(handler, tt.expectedHandler)
			assert.Equal(t, status, tt.expectedStatus)
		})
	}
}

func TestRestRouter_register(t *testing.T) {

	type fields struct {
		routeMap Routes
	}
	type args struct {
		path        string
		httpMethod  string
		handlerFunc HttpMethodHandler
		expected    HttpMethodHandler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "Should register to map given path and method",
			fields: fields{Routes{}},
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
				routes: tt.fields.routeMap,
			}
			r.register(tt.args.path, tt.args.httpMethod, tt.args.handlerFunc)
			reflect.DeepEqual(r.routes[tt.args.path][tt.args.httpMethod](nil), tt.args.expected(nil))
		})
	}
}

func Test_MatchPatchVariables(t *testing.T) {

	router := &RestRouter{
		routes: Routes{
			"/transactions":                        map[string]HttpMethodHandler{"GET": emptyHandlerFunc},
			"/transactions/:id/product/:productId": map[string]HttpMethodHandler{"GET": trnProductWithIdFunc},
		},
	}

	handler, _ := router.load("/transactions?name=yuri", http.MethodGet)
	assert.Equal(t, emptyHandlerFunc(nil), handler(nil))

	handler, _ = router.load("/transactions/:id/product/:productId", http.MethodGet)
	assert.Equal(t, trnProductWithIdFunc(nil), handler(nil))

	handler, _ = router.load("/transactions/:id/product/:productId", http.MethodGet)
	assert.Equal(t, trnProductWithIdFunc(nil), handler(nil))

}

func trnProductWithIdFunc(ctx IHttpContext) error {
	return errors.New("just a test")
}
