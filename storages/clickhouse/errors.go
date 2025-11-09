package clickhouse

import (
	"gitlab.com/algmib/kit"
)

const (
	ErrCodeClickOpen                 = "CH-003"
	ErrCodeClickGetVer               = "CH-008"
	ErrCodeClickLockTableCreation    = "CH-009"
	ErrCodeClickLockLifeViewCreation = "CH-010"
	ErrCodeClickLockTimeout          = "CH-011"
	ErrCodeClickPing                 = "CH-012"
)

var (
	ErrClickOpen = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeClickOpen, "").Wrap(cause).Err()
	}
	ErrClickGetVer = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeClickGetVer, "").Wrap(cause).Err()
	}
	ErrClickPing = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeClickPing, "ping").Wrap(cause).Err()
	}
)
