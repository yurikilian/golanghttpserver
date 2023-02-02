package logger

import (
	"context"
)

type Formatter interface {
	Format(ctx context.Context, fields Fields, message string) string
}
