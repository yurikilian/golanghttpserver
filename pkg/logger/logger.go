package logger

import "context"

type Level int

const (
	NoOp Level = iota
	TraceLevel
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type Logger interface {
	Trace(context.Context, string)
	Debug(context.Context, string)
	Info(context.Context, string)
	Warn(context.Context, string)
	Error(context.Context, string)
	Fatal(context.Context, string)
	SetExitCall(exitCall func())
}

var LevelLabelMap = map[Level]string{
	TraceLevel: "TRACE",
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
}

type Fields map[string]string

type (
	Options struct {
		Format Format
	}
	Format string
)
