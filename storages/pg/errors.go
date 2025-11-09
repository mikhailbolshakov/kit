package pg

import "github.com/mikhailbolshakov/kit"

const (
	ErrCodePostgresOpen = "PG-001"
	ErrCodePgEmptyJsonb = "PG-002"
	ErrCodePgSetJsonb   = "PG-003"
	ErrCodePgGetJsonb   = "PG-004"
)

var (
	ErrPostgresOpen = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodePostgresOpen, "").Wrap(cause).Err()
	}
	ErrPgEmptyJsonb = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodePgEmptyJsonb, "empty JSONB conversion").Wrap(cause).Err()
	}
	ErrPgSetJsonb = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodePgSetJsonb, "set JSONB").Wrap(cause).Err()
	}
	ErrPgGetJsonb = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodePgGetJsonb, "get JSONB").Wrap(cause).Err()
	}
)
