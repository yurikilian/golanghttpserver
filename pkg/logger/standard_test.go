package logger

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func Test_LogLevel(t *testing.T) {

	type options struct {
		level Level
	}

	type fields struct {
		options options
	}

	type args struct {
		level Level
		s     string
	}
	tests := []struct {
		name     string
		fields   fields
		args     []args
		expected string
	}{
		{
			name: "Should write message to all logs given trace log Level",
			fields: fields{
				options: options{
					level: TraceLevel,
				},
			},
			args: []args{
				{
					level: TraceLevel,
					s:     "This is a trace message",
				},
				{
					level: DebugLevel,
					s:     "This is a debug message",
				},
				{
					level: InfoLevel,
					s:     "This is an info message",
				},
				{
					level: WarnLevel,
					s:     "This is a warn message",
				},
				{
					level: ErrorLevel,
					s:     "This is a error message",
				},
				{
					level: FatalLevel,
					s:     "This is a fatal message",
				},
			},
			expected: "This is a trace message\n" +
				"This is a debug message\n" +
				"This is an info message\n" +
				"This is a warn message\n" +
				"This is a error message\n" +
				"This is a fatal message\n" +
				"Exit call was called\n",
		},
		{
			name: "Should write message to debug, info, warn, error and fatal given debug log Level",
			fields: fields{
				options: options{
					level: DebugLevel,
				},
			},
			args: []args{
				{
					level: TraceLevel,
					s:     "Should not write",
				},
				{
					level: DebugLevel,
					s:     "This is a debug message",
				},
				{
					level: InfoLevel,
					s:     "This is an info message",
				},
				{
					level: WarnLevel,
					s:     "This is a warn message",
				},
				{
					level: ErrorLevel,
					s:     "This is a error message",
				},
				{
					level: FatalLevel,
					s:     "This is a fatal message",
				},
			},
			expected: "This is a debug message\n" +
				"This is an info message\n" +
				"This is a warn message\n" +
				"This is a error message\n" +
				"This is a fatal message\n" +
				"Exit call was called\n",
		},
		{
			name: "Should write message to  info, warn, error and fatal given info log Level",
			fields: fields{
				options: options{
					level: InfoLevel,
				},
			},
			args: []args{
				{
					level: TraceLevel,
					s:     "Should not write",
				},
				{
					level: DebugLevel,
					s:     "Should not write",
				},
				{
					level: InfoLevel,
					s:     "This is an info message",
				},
				{
					level: WarnLevel,
					s:     "This is a warn message",
				},
				{
					level: ErrorLevel,
					s:     "This is a error message",
				},
				{
					level: FatalLevel,
					s:     "This is a fatal message",
				},
			},
			expected: "This is an info message\n" +
				"This is a warn message\n" +
				"This is a error message\n" +
				"This is a fatal message\n" +
				"Exit call was called\n",
		},
		{
			name: "Should write message to warn, error and fatal given info warn Level",
			fields: fields{
				options: options{
					level: WarnLevel,
				},
			},
			args: []args{
				{
					level: TraceLevel,
					s:     "Should not write",
				},
				{
					level: DebugLevel,
					s:     "Should not write",
				},
				{
					level: InfoLevel,
					s:     "Should not write",
				},
				{
					level: WarnLevel,
					s:     "This is a warn message",
				},
				{
					level: ErrorLevel,
					s:     "This is a error message",
				},
				{
					level: FatalLevel,
					s:     "This is a fatal message",
				},
			},
			expected: "This is a warn message\n" +
				"This is a error message\n" +
				"This is a fatal message\n",
		},
		{
			name: "Should write message to error and fatal given info error Level",
			fields: fields{
				options: options{
					level: ErrorLevel,
				},
			},
			args: []args{
				{
					level: TraceLevel,
					s:     "Should not write",
				},
				{
					level: DebugLevel,
					s:     "Should not write",
				},
				{
					level: InfoLevel,
					s:     "Should not write",
				},
				{
					level: WarnLevel,
					s:     "Should not write",
				},
				{
					level: ErrorLevel,
					s:     "This is a error message",
				},
				{
					level: FatalLevel,
					s:     "This is a fatal message",
				},
			},
			expected: "This is a error message\n" +
				"This is a fatal message\n",
		},
		{
			name: "Should write message to  fatal given fatal error Level",
			fields: fields{
				options: options{
					level: FatalLevel,
				},
			},
			args: []args{
				{
					level: TraceLevel,
					s:     "Should not write",
				},
				{
					level: DebugLevel,
					s:     "Should not write",
				},
				{
					level: InfoLevel,
					s:     "Should not write",
				},
				{
					level: WarnLevel,
					s:     "Should not write",
				},
				{
					level: ErrorLevel,
					s:     "Should not write",
				},
				{
					level: FatalLevel,
					s:     "This is a fatal message",
				},
			},
			expected: "This is a fatal message\n",
		},
		{
			name: "Should disable logging giving noop Level",
			fields: fields{
				options: options{
					level: NoOp,
				},
			},
			args: []args{
				{
					level: TraceLevel,
					s:     "Should not write",
				},
				{
					level: DebugLevel,
					s:     "Should not write",
				},
				{
					level: InfoLevel,
					s:     "Should not write",
				},
				{
					level: WarnLevel,
					s:     "Should not write",
				},
				{
					level: ErrorLevel,
					s:     "Should not write",
				},
				{
					level: FatalLevel,
					s:     "Should not write",
				},
			},
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r, w, _ := os.Pipe()
			output := w

			l := newLogToTest(output, tt.fields.options.level)

			for _, arg := range tt.args {
				l.log(nil, arg.level, arg.s)
			}

			output.Close()
			all, _ := io.ReadAll(r)
			assert.Equal(t, tt.expected, string(all))
		})
	}
}

func newLogToTest(w io.Writer, level Level) *StandardLog {
	log := NewStandardLog(StandardLogOptions{
		Output:    w,
		Level:     level,
		Formatter: NewTestFormatter(),
	})

	log.SetExitCall(func() {
		log.Info(context.TODO(), "Exit call was called")
	})
	return log
}

type TestFormatter struct {
}

func NewTestFormatter() *TestFormatter {
	return &TestFormatter{}
}

func (t *TestFormatter) Format(ctx context.Context, fields Fields, message string) string {
	return fmt.
		Sprintf(
			"%v",
			message,
		) //,

}
