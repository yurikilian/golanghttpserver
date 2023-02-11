package server

import (
	"context"
	"encoding/json"
	"github.com/yurikilian/bills/pkg/exception"
	"github.com/yurikilian/bills/pkg/logger"
	"net/http"
)

type IHttpContext interface {
	Writer() http.ResponseWriter
	Request() *http.Request
	ReqCtx() context.Context
	SetRequest(r *http.Request)
	WriteResponse(statusCode int, data interface{}) error
	Logger() logger.Logger
	ReadBody(bodyStruct interface{}) error
	reset(writer http.ResponseWriter, request *http.Request)
}

type HttpContext struct {
	writer  http.ResponseWriter
	request *http.Request
	log     logger.Logger
	binder  *Binder
}

func NewHttpContext(writer http.ResponseWriter, request *http.Request, log logger.Logger, binder *Binder) IHttpContext {
	return &HttpContext{writer: writer, request: request, log: log, binder: binder}
}

func (hCtx *HttpContext) reset(writer http.ResponseWriter, request *http.Request) {
	hCtx.request = request
	hCtx.writer = writer
}

func (hCtx *HttpContext) Writer() http.ResponseWriter {
	return hCtx.writer
}

func (hCtx *HttpContext) Request() *http.Request {
	return hCtx.request
}

func (hCtx *HttpContext) ReqCtx() context.Context {
	return hCtx.request.Context()
}

func (hCtx *HttpContext) SetRequest(r *http.Request) {
	hCtx.request = r
}

func (hCtx *HttpContext) WriteResponse(statusCode int, data interface{}) error {
	hCtx.writer.Header().Set("Content-Type", "application/json")
	hCtx.writer.WriteHeader(statusCode)

	err := json.NewEncoder(hCtx.Writer()).Encode(data)
	if err != nil {
		return exception.NewInternalServerError(err.Error())
	}

	return nil
}
func (hCtx *HttpContext) Logger() logger.Logger {
	return hCtx.log
}

func (hCtx *HttpContext) ReadBody(bodyStruct interface{}) error {
	return hCtx.binder.ReadBody(hCtx, bodyStruct)
}
