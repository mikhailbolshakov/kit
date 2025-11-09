package http

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"net/http"
)

const (
	ErrCodeHttpTest                          = "HTTP-000"
	ErrCodeHttpSrvListen                     = "HTTP-001"
	ErrCodeHttpDecodeRequest                 = "HTTP-002"
	ErrCodeHttpUrlVar                        = "HTTP-003"
	ErrCodeHttpCurrentUser                   = "HTTP-004"
	ErrCodeHttpUrlVarEmpty                   = "HTTP-005"
	ErrCodeHttpUrlFormVarEmpty               = "HTTP-006"
	ErrCodeHttpUrlFormVarNotInt              = "HTTP-007"
	ErrCodeHttpUrlFormVarNotTime             = "HTTP-008"
	ErrCodeHttpMultipartParseForm            = "HTTP-009"
	ErrCodeHttpMultipartEmptyContent         = "HTTP-010"
	ErrCodeHttpMultipartNotMultipart         = "HTTP-011"
	ErrCodeHttpMultipartParseMediaType       = "HTTP-012"
	ErrCodeHttpMultipartWrongMediaType       = "HTTP-013"
	ErrCodeHttpMultipartMissingBoundary      = "HTTP-014"
	ErrCodeHttpMultipartEofReached           = "HTTP-015"
	ErrCodeHttpMultipartNext                 = "HTTP-016"
	ErrCodeHttpMultipartFormNameFileExpected = "HTTP-017"
	ErrCodeHttpMultipartFilename             = "HTTP-018"
	ErrCodeHttpCurrentClient                 = "HTTP-019"
	ErrCodeHttpUrlFormVarNotFloat            = "HTTP-020"
	ErrCodeHttpUrlFormVarNotBool             = "HTTP-021"
	ErrCodeHttpUrlWrongSortFormat            = "HTTP-022"
	ErrCodeHttpUrlVarInvalidUUID             = "HTTP-023"
	ErrCodeHttpUrlMaxPageSizeExceeded        = "HTTP-024"
	ErrCodeHttpCurrentPartner                = "HTTP-025"
	ErrCodeHttpFileHeaderEmpty               = "HTTP-026"
	ErrCodeHttpFileHeaderInvalidJson         = "HTTP-027"
	ErrCodeHttpFileHeaderInvalidUUID         = "HTTP-028"
	ErrCodeHttpProxyFileNewRequest           = "HTTP-029"
	ErrCodeHttpProxyFileInvalidContext       = "HTTP-030"
	ErrCodeHttpProxyFileClientDo             = "HTTP-031"
	ErrCodeHttpProxyFileCreatePart           = "HTTP-032"
	ErrCodeHttpProxyFileCopyFile             = "HTTP-033"
	ErrCodeHttpProxyFileWriteField           = "HTTP-034"
	ErrCodeHttpProxyFileReadResponse         = "HTTP-035"
	ErrCodeHttpProxyFileJsonUnmarshal        = "HTTP-036"
	ErrCodeHttpRequestEmptyUrl               = "HTTP-037"
	ErrCodeHttpRequestEmptyMethod            = "HTTP-038"
	ErrCodeHttpRequestEmptyClient            = "HTTP-039"
	ErrCodeHttpRequestEmptyResponseFunc      = "HTTP-040"
	ErrCodeHttpRequestInvalidUrl             = "HTTP-041"
	ErrCodeHttpRequestInvalidStatusCode      = "HTTP-042"
	ErrCodeHttpRequestReadAll                = "HTTP-043"
	ErrCodeHttpRequestResponseFuncFailed     = "HTTP-044"
	ErrCodeHttpRequestJsonMarshal            = "HTTP-045"
	ErrCodeHttpRequestJsonUnmarshal          = "HTTP-046"
	ErrCodeHttpRequestWriterWriteField       = "HTTP-047"
	ErrCodeHttpRequestWriterCreate           = "HTTP-048"
	ErrCodeHttpRequestWriterCopy             = "HTTP-049"
	ErrCodeHttpRequestWriterClose            = "HTTP-050"
	ErrCodeHttpRequestEmptyFileData          = "HTTP-051"
	ErrCodeHttpRequestNew                    = "HTTP-052"
	ErrCodeHttpRequestDo                     = "HTTP-053"
	ErrCodeHttpUrlInvalidPagingSortSet       = "HTTP-054"
	ErrCodeHttpCalculateLengthCopy           = "HTTP-055"
	ErrCodeHttpValidationRequest             = "HTTP-056"
	ErrCodeHttpSanitizeRequest               = "HTTP-057"
)

var (
	ErrHttpTest = func() error {
		return kit.NewAppErrBuilder(ErrCodeHttpTest, "").Business().Err()
	}
	ErrHttpSrvListen = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpSrvListen, "").Wrap(cause).Err()
	}
	ErrHttpDecodeRequest = func(cause error, ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpDecodeRequest, "invalid request").Wrap(cause).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpValidationRequest = func(cause error, ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpValidationRequest, "validation").Wrap(cause).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpSanitizeRequest = func(cause error, ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpSanitizeRequest, "sanitizing").Wrap(cause).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlVar = func(ctx context.Context, v string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpUrlVar, "invalid or empty URL parameter").F(kit.KV{"var": v}).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpCurrentUser = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpCurrentUser, `cannot obtain current user`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlVarEmpty = func(ctx context.Context, v string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpUrlVarEmpty, `URL parameter is empty`).Business().F(kit.KV{"var": v}).C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlVarInvalidUUID = func(ctx context.Context, v string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpUrlVarInvalidUUID, `invalid UUID`).Business().F(kit.KV{"var": v}).C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlFormVarEmpty = func(ctx context.Context, v string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpUrlFormVarEmpty, `URL form value is empty`).Business().F(kit.KV{"var": v}).C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlFormVarNotInt = func(cause error, ctx context.Context, v string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpUrlFormVarNotInt, "form value must be of int type").Wrap(cause).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlFormVarNotFloat = func(cause error, ctx context.Context, v string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpUrlFormVarNotFloat, "form value must be of float type").Wrap(cause).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlFormVarNotBool = func(cause error, ctx context.Context, v string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpUrlFormVarNotBool, "form value must be of bool type").Wrap(cause).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlFormVarNotTime = func(cause error, ctx context.Context, v string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpUrlFormVarNotTime, "form value must be of time type in RFC-3339 format").Wrap(cause).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpFileHeaderEmpty = func(ctx context.Context, v string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpFileHeaderEmpty, `file header is empty`).Business().F(kit.KV{"var": v}).C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpFileHeaderInvalidUUID = func(ctx context.Context, v string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpFileHeaderInvalidUUID, `invalid UUID`).Business().F(kit.KV{"var": v}).C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpFileHeaderInvalidJson = func(ctx context.Context, v string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpFileHeaderInvalidJson, `file header json is invalid`).Business().F(kit.KV{"var": v}).C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartParseForm = func(cause error, ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpMultipartParseForm, "parse multipart form").Wrap(cause).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartEmptyContent = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpMultipartEmptyContent, `content is empty`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartNotMultipart = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpMultipartNotMultipart, `content isn't multipart`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartParseMediaType = func(cause error, ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpMultipartParseMediaType, "parse media type").Wrap(cause).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartWrongMediaType = func(ctx context.Context, mt string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpMultipartWrongMediaType, `wrong media type %s`, mt).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartMissingBoundary = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpMultipartMissingBoundary, `missing boundary`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartEofReached = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpMultipartEofReached, `no parts found`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartNext = func(cause error, ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpMultipartNext, "reading part").Wrap(cause).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartFormNameFileExpected = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpMultipartFormNameFileExpected, `correct part must have name="file" param`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartFilename = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpMultipartFilename, `filename is empty`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpCurrentClient = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpCurrentClient, `cannot obtain current client`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpCurrentPartner = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpCurrentPartner, `cannot obtain current partner`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlWrongSortFormat = func(ctx context.Context, v string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpUrlWrongSortFormat, "wrong sort format").Business().F(kit.KV{"var": v}).C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlMaxPageSizeExceeded = func(ctx context.Context, maxPageSize int) error {
		return kit.NewAppErrBuilder(ErrCodeHttpUrlMaxPageSizeExceeded, "max page size (%d) exceeded", maxPageSize).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlInvalidPagingSortSet = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpUrlInvalidPagingSortSet, "invalid sort set").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpProxyFileNewRequest = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpProxyFileNewRequest, "new request").Wrap(cause).C(ctx).Err()
	}
	ErrHttpProxyFileInvalidContext = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpProxyFileInvalidContext, "invalid context").Wrap(cause).C(ctx).Err()
	}
	ErrHttpProxyFileClientDo = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpProxyFileClientDo, "client do failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpProxyFileCreatePart = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpProxyFileCreatePart, "create part failed").Wrap(cause).Err()
	}
	ErrHttpProxyFileCopyFile = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpProxyFileCopyFile, "copy failed").Wrap(cause).Err()
	}
	ErrHttpProxyFileWriteField = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpProxyFileWriteField, "write field failed").Wrap(cause).Err()
	}
	ErrHttpProxyFileReadResponse = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpProxyFileReadResponse, "read response failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpProxyFileJsonUnmarshal = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpProxyFileJsonUnmarshal, "unmarshall failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpRequestEmptyUrl = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestEmptyUrl, "request: empty url").C(ctx).Err()
	}
	ErrHttpRequestEmptyMethod = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestEmptyMethod, "request: empty method").C(ctx).Err()
	}
	ErrHttpRequestEmptyClient = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestEmptyClient, "request: empty client").C(ctx).Err()
	}
	ErrHttpRequestEmptyResponseFunc = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestEmptyResponseFunc, "request: empty response function").C(ctx).Err()
	}
	ErrHttpRequestInvalidUrl = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestInvalidUrl, "request: invalid url").Wrap(cause).C(ctx).Err()
	}
	ErrHttpRequestInvalidStatusCode = func(ctx context.Context, code int, detail string) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestInvalidStatusCode, "invalid status code").F(kit.KV{"code": code, "detail": detail}).C(ctx).Err()
	}
	ErrHttpRequestReadAll = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestReadAll, "io: read all failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpRequestNew = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestNew, "request create failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpRequestDo = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestDo, "request execution failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpRequestJsonMarshal = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestJsonMarshal, "json: marshal failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpCalculateLengthCopy = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpCalculateLengthCopy, "io: copy failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpRequestJsonUnmarshal = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestJsonUnmarshal, "json: unmarshal failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpRequestWriterWriteField = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestWriterWriteField, "writer: write failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpRequestWriterCreate = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestWriterCreate, "writer: create failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpRequestWriterCopy = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestWriterCopy, "writer: copy failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpRequestWriterClose = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestWriterClose, "writer: close failed").Wrap(cause).C(ctx).Err()
	}
	ErrHttpRequestEmptyFileData = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestEmptyFileData, "empty file data").Business().C(ctx).Err()
	}
	ErrHttpRequestResponseFuncFailed = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeHttpRequestResponseFuncFailed, "response fn failed").Wrap(cause).C(ctx).Err()
	}
)
