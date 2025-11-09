package goroutine

import (
	"context"
	"github.com/mikhailbolshakov/kit"
)

const (
	ErrCodeGoroutineNoLogger = "GORTN-001"
)

var (
	ErrGoroutineNoLogger = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeGoroutineNoLogger, "either logger or logger func must be specified").C(ctx).Err()
	}
)
