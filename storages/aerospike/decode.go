package aerospike

import (
	"context"

	aero "github.com/aerospike/aerospike-client-go/v7"
	"gitlab.com/algmib/kit"
)

func Decode[T any](ctx context.Context, bin string, rec *aero.Record) (*T, error) {
	body, err := AsBytes(ctx, rec.Bins, bin)
	if err != nil {
		return nil, ErrAeroBodyBytes(ctx, err)
	}
	if body == nil {
		return nil, nil
	}
	var result *T
	err = kit.Unmarshal(body, &result)
	if err != nil {
		return nil, ErrAeroUnmarshal(ctx, err)
	}
	return result, nil
}
