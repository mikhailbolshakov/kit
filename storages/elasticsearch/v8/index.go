package v8

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/mikhailbolshakov/kit"
	"io"
	"net/http"
)

type indexImpl struct {
	*esImpl
}

func (s *indexImpl) NewBuilder() IndexBuilder {
	return newIndexBuilder(s.indexImpl, s.config, s.logger)
}

func (s *indexImpl) Refresh(ctx context.Context, index string) error {
	s.l().C(ctx).Mth("refresh").F(kit.KV{"index": index}).Dbg()

	if err := s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Indices.Refresh(
			s.client.Indices.Refresh.WithContext(ctx),
			s.client.Indices.Refresh.WithIndex(index),
		)
	}, NoResponseFn); err != nil {
		return ErrIndexRefresh(ctx, err, index)
	}

	return nil
}

func (s *indexImpl) UpdateAliases(ctx context.Context, action *AliasAction) error {
	s.l().C(ctx).Mth("update-aliases").Dbg()

	body, err := kit.JsonEncode(action)
	if err != nil {
		return err
	}

	if err = s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Indices.UpdateAliases(
			bytes.NewReader(body),
			s.client.Indices.UpdateAliases.WithContext(ctx),
		)
	}, NoResponseFn); err != nil {
		return ErrIndexUpdateAliases(ctx, err)
	}

	return nil
}

// GetIndices response example
//
//	{
//		"c6id8uqjrpf68ni87z7yddrpnr-idx-1" : {
//			"aliases" : {
//				"c6id8uqjrpf68ni87z7yddrpnr" : {
//					"is_write_index" : true
//				}
//			}
//		},
//		"c6id8uqjrpf68ni87z7yddrpnr-idx-2" : {
//			"aliases" : {
//				"c6id8uqjrpf68ni87z7yddrpnr" : { }
//			}
//		}
//	}
func (s *indexImpl) GetIndices(ctx context.Context, alias string) (*GetIndicesResponse, error) {
	s.l().C(ctx).Mth("get-indices").F(kit.KV{"alias": alias}).Dbg()

	var response GetIndicesResponse

	if err := s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Indices.GetAlias(
			s.client.Indices.GetAlias.WithName(alias),
			s.client.Indices.GetAlias.WithContext(ctx),
		)
	}, func(code int, data io.ReadCloser) error {
		return json.NewDecoder(data).Decode(&response.Indices)
	}); err != nil {
		return nil, ErrIndexGetIndices(ctx, err, alias)
	}

	return &response, nil

}

func (s *indexImpl) PutMapping(ctx context.Context, index string, mapping MappingBody) error {
	s.l().C(ctx).Mth("put-mapping").F(kit.KV{"index": index}).Dbg()

	if err := s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Indices.PutMapping(
			[]string{index},
			mapping.Reader(),
			s.client.Indices.PutMapping.WithContext(ctx),
		)
	}, NoResponseFn); err != nil {
		return ErrIndexPutMapping(ctx, err, index)
	}

	return nil
}

func (s *indexImpl) GetMapping(ctx context.Context, index string) (*GetMappingResponse, error) {
	s.l().C(ctx).Mth("get-mapping").F(kit.KV{"index": index}).Dbg()

	var response GetMappingResponse

	if err := s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Indices.GetMapping(
			s.client.Indices.GetMapping.WithContext(ctx),
			s.client.Indices.GetMapping.WithIndex(index),
		)
	}, func(code int, data io.ReadCloser) error {
		return json.NewDecoder(data).Decode(&response.Mappings)
	}); err != nil {
		return nil, ErrIndexGetMapping(ctx, err, index)
	}

	return &response, nil
}

func (s *indexImpl) Create(ctx context.Context, index string, mapping MappingBody) error {
	s.l().C(ctx).Mth("create").F(kit.KV{"index": index}).Dbg()

	if err := s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Indices.Create(
			index,
			s.client.Indices.Create.WithContext(ctx),
			s.client.Indices.Create.WithBody(mapping.Reader()),
		)
	}, NoResponseFn); err != nil {
		return ErrIndexCreate(ctx, err, index)
	}

	return nil
}

func (s *indexImpl) Exists(ctx context.Context, index string) (bool, error) {
	s.l().C(ctx).Mth("exists").F(kit.KV{"index": index}).Dbg()

	var exists bool
	if err := s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Indices.Exists(
			[]string{index},
			s.client.Indices.Exists.WithContext(ctx),
		)
	}, func(code int, data io.ReadCloser) error {
		exists = code == http.StatusOK
		return nil
	}, ExistsResponseCodes...); err != nil {
		return false, ErrIndexExists(ctx, err, index)
	}

	return exists, nil
}

func (s *indexImpl) Delete(ctx context.Context, indices []string) error {
	s.l().C(ctx).Mth("delete").F(kit.KV{"indices-len": len(indices)}).Dbg()

	if err := s.Do(ctx, func() (*esapi.Response, error) {
		return s.client.Indices.Delete(
			indices,
			s.client.Indices.Delete.WithContext(ctx),
		)
	}, NoResponseFn); err != nil {
		return ErrIndexDelete(ctx, err, indices)
	}

	return nil
}

func (s *indexImpl) DeleteIndices(ctx context.Context, alias string) error {
	s.l().C(ctx).Mth("delete-indices").F(kit.KV{"alias": alias}).Dbg()

	rs, err := s.GetIndices(ctx, alias)
	if err != nil {
		return err
	}
	if len(rs.Indices) == 0 {
		// nothing to do, break
		return nil
	}

	return s.Delete(ctx, kit.MapKeys(rs.Indices))
}
