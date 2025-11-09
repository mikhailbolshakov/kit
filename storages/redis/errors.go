package redis

import (
	"context"
	"github.com/mikhailbolshakov/kit"
)

const (
	ErrCodeRedisPingErr                   = "RDS-001"
	ErrCodeRedisPriorityQueuePushErr      = "RDS-002"
	ErrCodeRedisPriorityQueuePopErr       = "RDS-003"
	ErrCodeRedisPriorityQueuePopRemoveErr = "RDS-004"
)

var (
	ErrRedisPingErr = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeRedisPingErr, "").Wrap(cause).Err()
	}
	ErrRedisPriorityQueuePushErr = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeRedisPriorityQueuePushErr, "priority queue: push").Wrap(cause).Err()
	}
	ErrRedisPriorityQueuePopErr = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeRedisPriorityQueuePopErr, "priority queue: pop").Wrap(cause).Err()
	}
	ErrRedisPriorityQueuePopRemoveErr = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeRedisPriorityQueuePopRemoveErr, "priority queue: pop and remove").Wrap(cause).Err()
	}
)
