package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/yurikilian/bills/internal/logger"
	"github.com/yurikilian/bills/pkg/exception"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestRestServer_Router(t *testing.T) {
	type fields struct {
		server *RestServer
	}
	type args struct {
		router *RestRouter
	}

	server := NewRestServer(&Options{BindAddress: ":8080"})

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RestServer
	}{
		{
			name:   "Should return rest server on Route method",
			fields: fields{server: server},
			args:   args{router: NewRestRouter()},
			want:   server,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fields.server
			if got := r.Router(tt.args.router); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Router() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestRestServer_ServeHTTP(t *testing.T) {

	server := NewRestServer(&Options{BindAddress: ":8080"})
	handlerFunc := func(httpContext *HttpContext) error {
		return httpContext.WriteResponse(200, "Test")
	}
	server.Router(NewRestRouter().Get("/test", handlerFunc))

	type args struct {
		w   *httptest.ResponseRecorder
		req *http.Request
	}

	type Response struct {
		statusCode int
		body       string
	}

	tests := []struct {
		name             string
		args             args
		expectedResponse *Response
	}{
		{
			name: "Should process request given valid path",
			args: args{
				w:   httptest.NewRecorder(),
				req: newRequest(http.MethodGet, "/test", nil),
			},
			expectedResponse: &Response{
				statusCode: http.StatusOK,
				body:       "\"Test\"\n",
			},
		},
		{
			name: "Should return not found given not registered path",
			args: args{
				w:   httptest.NewRecorder(),
				req: newRequest(http.MethodGet, "/notfound", nil),
			},
			expectedResponse: &Response{
				statusCode: http.StatusNotFound,
				body:       getJson(t, exception.NewRouteNotFound("/notfound")),
			},
		},
		{
			name: "Should return not found given not registered path and different method",
			args: args{
				w:   httptest.NewRecorder(),
				req: newRequest(http.MethodPut, "/notfound", nil),
			},
			expectedResponse: &Response{
				statusCode: http.StatusNotFound,
				body:       getJson(t, exception.NewRouteNotFound("/notfound")),
			},
		},
		{
			name: "Should return method not allowed given existent path but wrong method",
			args: args{
				w:   httptest.NewRecorder(),
				req: newRequest(http.MethodPut, "/test", nil),
			},
			expectedResponse: &Response{
				statusCode: http.StatusMethodNotAllowed,
				body:       getJson(t, exception.NewMethodNotAllowed("/test", http.MethodPut)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := tt.args.w
			server.ServeHTTP(recorder, tt.args.req)
			if tt.expectedResponse != nil {
				assert.Equal(t, tt.expectedResponse.statusCode, recorder.Code)
				assert.Equal(t, tt.expectedResponse.body, recorder.Body.String())
			}
		})
	}
}

func TestRestServer_Start(t *testing.T) {
	type fields struct {
		server *RestServer
		router *RestRouter
	}
	tests := []struct {
		name          string
		fields        fields
		expectedErr   *exception.Problem
		expectedStart bool
	}{
		{
			name: "Return exception given invalid server configuration",
			fields: fields{
				router: NewRestRouter(),
				server: NewRestServer(&Options{Log: logger.NewProvider().ProvideLog(), BindAddress: "-1"}),
			},
			expectedErr:   exception.NewInternalServerError("listen tcp: address -1: missing port in address"),
			expectedStart: false,
		},
		{
			name: "Return exception given invalid server configuration",
			fields: fields{
				router: NewRestRouter(),
				server: NewRestServer(&Options{BindAddress: "-1"}),
			},
			expectedErr:   exception.NewInternalServerError("Key: 'Options.Log' Error:Field validation for 'Log' failed on the 'required' tag"),
			expectedStart: false,
		},
		{
			name: "Should start server",
			fields: fields{
				router: NewRestRouter(),
				server: NewRestServer(&Options{BindAddress: ":0"}),
			},
			expectedStart: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fields.server.Router(tt.fields.router)

			if tt.expectedStart {

				go func() {
					time.Sleep(1 * time.Second)
					err := tt.fields.server.getHttpServer().Shutdown(context.Background())
					assert.NoError(t, err)
				}()

				assert.Error(t, errors.New("http: Server closed"))
			} else if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, r.Start(nil))
			}

		})
	}
}

func getJson(t *testing.T, ex *exception.Problem) string {
	marshal, err := json.Marshal(ex)
	assert.NoError(t, err)
	return string(marshal)
}

func newRequest(method string, path string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(method, path, body)
	return request
}

func BenchmarkRestServer_ServeHTTP(b *testing.B) {

	server := NewRestServer(&Options{BindAddress: ":0"})
	mutex := &sync.Mutex{}
	handlerFunc := func(ctx *HttpContext) error {
		mutex.Lock()
		ctx.WriteResponse(200, "Test")
		mutex.Unlock()
		return nil
	}
	writer := httptest.NewRecorder()

	server.Router(NewRestRouter().Get("/test", handlerFunc))
	request := newRequest(http.MethodGet, "/test", nil)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for i := 0; i < b.N; i++ {
			for pb.Next() {
				server.ServeHTTP(writer, request)
			}
		}
	})

}
