package v8

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gitlab.com/algmib/kit"
	"io"
	"net/http"
)

type documentImpl struct {
	*esImpl
}

// Index indexes a document
func (s *documentImpl) Index(ctx context.Context, index string, id string, data DataBody) error {
	s.l().C(ctx).Mth("index").F(kit.KV{"index": index, "id": id}).Dbg()

	if err := s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Index(
			index,
			data.Reader(),
			s.client.Index.WithDocumentID(id),
			s.client.Index.WithContext(ctx),
		)
	}, NoResponseFn); err != nil {
		return ErrDocIndex(ctx, err, index, id)
	}
	// refresh
	if s.config.Refresh {
		return s.Instance().Index().Refresh(ctx, index)
	}
	return nil
}

// Update updates an existing document in the specified index with the provided data.
func (s *documentImpl) Update(ctx context.Context, index string, id string, data DataBody) error {
	s.l().C(ctx).Mth("update").F(kit.KV{"index": index, "id": id}).Dbg()

	if err := s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Update(
			index,
			id,
			data.Reader(),
			s.client.Update.WithContext(ctx),
		)
	}, NoResponseFn); err != nil {
		return ErrDocUpdate(ctx, err, index, id)
	}
	// refresh
	if s.config.Refresh {
		return s.Instance().Index().Refresh(ctx, index)
	}
	return nil
}

// Exists checks if a document exists in the index
func (s *documentImpl) Exists(ctx context.Context, index, id string) (bool, error) {
	s.l().C(ctx).Mth("exists").F(kit.KV{"index": index, "id": id}).Dbg()

	var exists bool
	if err := s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Exists(index, id, s.client.Exists.WithContext(ctx))
	}, func(code int, data io.ReadCloser) error {
		exists = code == http.StatusOK
		return nil
	}, ExistsResponseCodes...); err != nil {
		return false, ErrDocExists(ctx, err, index, id)
	}

	return exists, nil
}

// Delete deletes a document
func (s *documentImpl) Delete(ctx context.Context, index string, id string) error {
	s.l().C(ctx).Mth("delete").F(kit.KV{"index": index, "id": id}).Dbg()

	if err := s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Delete(index, id, s.client.Delete.WithContext(ctx))
	}, NoResponseFn); err != nil {
		return ErrDocDelete(ctx, err, index, id)
	}
	// refresh
	if s.config.Refresh {
		return s.Instance().Index().Refresh(ctx, index)
	}
	return nil
}
