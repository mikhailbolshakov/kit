package v8

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"net/http"
)

type processorImpl struct {
}

var NoResponseFn = func(code int, data io.ReadCloser) error {
	return nil
}

var ExistsResponseCodes = []int{
	http.StatusOK,
	http.StatusNotFound,
}

type DoFn func() (*esapi.Response, error)
type ResponseFn func(code int, data io.ReadCloser) error

func (p *processorImpl) Do(ctx context.Context, fn DoFn, rs ResponseFn, validStatusCodes ...int) error {
	if fn == nil {
		return ErrExecuteFuncEmpty(ctx)
	}
	if rs == nil {
		return ErrExecuteResponseFuncEmpty(ctx)
	}
	res, err := fn()
	if err != nil {
		return ErrExecute(ctx, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		// maybe this code is valid
		for _, code := range validStatusCodes {
			if res.StatusCode == code {
				goto response
			}
		}
		bodyBytes, _ := io.ReadAll(res.Body)
		return ErrExecuteResponseCode(ctx, res.StatusCode, string(bodyBytes))
	}
response:
	if err = rs(res.StatusCode, res.Body); err != nil {
		return ErrExecuteResponseProcessing(ctx, err)
	}
	return nil
}
