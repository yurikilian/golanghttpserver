package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type StandardLogOptions struct {
	Output    io.Writer
	Level     Level
	Formatter Formatter
}

type StandardLog struct {
	stdLogger *log.Logger
	options   *StandardLogOptions
	exitCall  func()
}

func NewStandardLog(options *StandardLogOptions) *StandardLog {
	if options == nil {
		options = &StandardLogOptions{
			Output:    os.Stderr,
			Level:     InfoLevel,
			Formatter: NewStandardFormatter(),
		}
	}

	standardLog := StandardLog{
		stdLogger: log.New(options.Output, "", log.Lmsgprefix),
		options:   options,
	}

	standardLog.SetExitCall(func() {
		os.Exit(1)
	})

	return &standardLog
}

func (l *StandardLog) Trace(ctx context.Context, message string) {
	l.log(ctx, TraceLevel, message)
}

func (l *StandardLog) Debug(ctx context.Context, message string) {
	l.log(ctx, DebugLevel, message)
}

func (l *StandardLog) Info(ctx context.Context, message string) {
	l.log(ctx, InfoLevel, message)
}

func (l *StandardLog) Warn(ctx context.Context, message string) {
	l.log(ctx, WarnLevel, message)
}

func (l *StandardLog) Error(ctx context.Context, message string) {
	l.log(ctx, ErrorLevel, message)
}

func (l *StandardLog) Fatal(ctx context.Context, message string) {
	l.log(ctx, FatalLevel, message)
}

func (l *StandardLog) exit() {
	l.exitCall()
}

func (l *StandardLog) SetExitCall(exitCall func()) {
	l.exitCall = exitCall
}

func (l *StandardLog) log(ctx context.Context, level Level, message string) {
	if l.options.Level == NoOp {
		return
	}

	line := l.options.Formatter.
		Format(ctx, map[string]string{
			"timestamp":      time.Now().Format(time.RFC3339),
			"severityText":   LevelLabelMap[level],
			"severityNumber": strconv.Itoa(int(level)),
		}, message)

	if isLevelWritableAndInRange(level) && level >= l.options.Level {

		if level >= FatalLevel {
			l.stdLogger.Output(2, line)
			l.exit()
		} else {
			l.stdLogger.Println(line)
		}
	}
}

func isLevelWritableAndInRange(level Level) bool {
	return level >= TraceLevel && level <= FatalLevel
}

type StandardFormatter struct {
}

func NewStandardFormatter() Formatter {
	return &StandardFormatter{}
}

func (t *StandardFormatter) Format(ctx context.Context, fields Fields, message string) string {
	return fmt.
		Sprintf(
			"%v %v %v %v",
			fields["timestamp"],
			fields["severityText"],
			fields["severityNumber"],
			message,
		) //,

}
