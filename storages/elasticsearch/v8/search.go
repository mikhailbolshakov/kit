package v8

import (
	"context"
	"io"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gitlab.com/algmib/kit"
)

type searchImpl struct {
	*esImpl
}

func (s *searchImpl) NewBuilder() QueryBuilder {
	return newQueryBuilder()
}

func (s *searchImpl) Search(ctx context.Context, index string, body QueryBody) (*SearchResponse, error) {
	s.l().C(ctx).Mth("search").F(kit.KV{"index": index}).Dbg()

	var response SearchResponse

	if err := s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Search(
			s.client.Search.WithContext(ctx),
			s.client.Search.WithIndex(index),
			s.client.Search.WithBody(body.Reader()),
		)
	}, func(code int, data io.ReadCloser) error {
		return kit.NewDecoder(data).Decode(&response)
	}); err != nil {
		return nil, ErrSearch(ctx, err, index)
	}

	return &response, nil
}
