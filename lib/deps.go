package lib

import (
	"context"
)

type Logger interface {
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Warning(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})

	WithError(ctx context.Context, err error) context.Context
}

type Metrics interface {
	Increment(key string)
	Timer(key string) func()
}
