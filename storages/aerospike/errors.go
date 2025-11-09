package aerospike

import (
	"context"
	"github.com/mikhailbolshakov/kit"
)

var (
	ErrCodeAeroConn           = "AERO-001"
	ErrCodeAeroClosed         = "AERO-002"
	ErrCodeAeroNewKey         = "AERO-003"
	ErrCodeAeroInvalidBinType = "AERO-004"
	ErrCodeAeroBodyBytes      = "AERO-005"
	ErrCodeAeroUnmarshal      = "AERO-006"
)

var (
	ErrAeroConn = func(cause error, ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeAeroConn, "").Wrap(cause).C(ctx).Err()
	}
	ErrAeroNewKey = func(cause error, ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeAeroNewKey, "").Wrap(cause).C(ctx).Err()
	}
	ErrAeroClosed = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeAeroClosed, "dealing with closed instance").C(ctx).Err()
	}
	ErrAeroInvalidBinType = func(ctx context.Context, bin string) error {
		return kit.NewAppErrBuilder(ErrCodeAeroInvalidBinType, "invalid bin type").F(kit.KV{"bin": bin}).C(ctx).Err()
	}
	ErrAeroBodyBytes = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeAeroBodyBytes, "aero bytes").Wrap(cause).C(ctx).Err()
	}
	ErrAeroUnmarshal = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeAeroUnmarshal, "unmarshal failed").Wrap(cause).C(ctx).Err()
	}
)
