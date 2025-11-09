package v7

import (
	"context"
	"github.com/olivere/elastic/v7"
	"gitlab.com/algmib/kit"
)

func ToSortRequestEs(ctx context.Context, request *kit.SortRequest) (*elastic.SortInfo, error) {

	if request.Field == "" {
		return nil, ErrEsSortRequestFieldEmpty(ctx)
	}

	res := EsSortRequestMissingFirst
	if request.NullsLast {
		res = EsSortRequestMissingLast
	}

	return &elastic.SortInfo{
		Field:     request.Field,
		Ascending: !request.Desc,
		Missing:   res,
	}, nil
}
