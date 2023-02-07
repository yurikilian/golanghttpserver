package server

import "github.com/yurikilian/bills/pkg/logger"

type Options struct {
	BindAddress string        `validate:"required"`
	Log         logger.Logger `validate:"required"`
}

func NewRestServerOptions(bindAddress string, log logger.Logger) *Options {
	return &Options{BindAddress: bindAddress, Log: log}
}
