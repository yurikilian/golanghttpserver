package logger

import (
	"context"
	"fmt"
	"github.com/yurikilian/bills/pkg/logger"
	"go.opentelemetry.io/otel/trace"
)

type OtelLogFormatter struct {
}

func NewOtelLogFormatter() *OtelLogFormatter {
	return &OtelLogFormatter{}
}

func (t *OtelLogFormatter) Format(ctx context.Context, fields logger.Fields, message string) string {

	//TODO: refactor fields to well-defined domain model
	span := trace.SpanFromContext(ctx)

	return fmt.
		Sprintf(
			"%v %v %v %v %v %v",
			fields["timestamp"],
			fields["severityText"],
			fields["severityNumber"],
			span.SpanContext().TraceID(),
			span.SpanContext().SpanID(),
			message,
		) //,

}

var _ logger.Formatter = (*OtelLogFormatter)(nil)
