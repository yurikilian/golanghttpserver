package logger

import (
	"github.com/yurikilian/bills/pkg/logger"
	"os"
)

type Provider struct {
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ProvideLog() logger.Logger {
	return logger.NewStandardLog(&logger.StandardLogOptions{
		Output:    os.Stderr,
		Level:     logger.DebugLevel,
		Formatter: logger.NewStandardFormatter(),
	})
}
