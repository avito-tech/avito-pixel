package extras

import (
	"context"
	"fmt"
)

type Logger struct{}

type key int

const (
	errorKey key = iota
)

func NewLogger() Logger {
	return Logger{}
}

func (l Logger) Debug(ctx context.Context, args ...interface{}) {
	fmt.Println("[DEBUG] ", args)
}

func (l Logger) Info(ctx context.Context, args ...interface{}) {
	fmt.Println("[INFO] ", args)
}

func (l Logger) Warning(ctx context.Context, args ...interface{}) {
	fmt.Println("[WARNING] ", args)
}

func (l Logger) Error(ctx context.Context, args ...interface{}) {
	fmt.Println("[ERROR] ", args)
}

func (l Logger) WithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errorKey, err)
}
