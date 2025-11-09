package v8

import (
	"context"
	"github.com/mikhailbolshakov/kit"
)

var (
	ErrCodeNewClient                                 = "ESv8-001"
	ErrCodeDocIndex                                  = "ESv8-002"
	ErrCodeDocDelete                                 = "ESv8-003"
	ErrCodeDocExists                                 = "ESv8-004"
	ErrCodeIndexRefresh                              = "ESv8-005"
	ErrCodeIndexCreate                               = "ESv8-006"
	ErrCodeExecuteResponseCode                       = "ESv8-007"
	ErrCodeIndexDelete                               = "ESv8-008"
	ErrCodeIndexExists                               = "ESv8-009"
	ErrCodeIndexUpdateAliases                        = "ESv8-010"
	ErrCodeSearch                                    = "ESv8-011"
	ErrCodeIndexPutMapping                           = "ESv8-012"
	ErrCodeIndexGetIndices                           = "ESv8-013"
	ErrCodeIndexGetMapping                           = "ESv8-014"
	ErrCodeExecuteFuncEmpty                          = "ESv8-015"
	ErrCodeExecuteResponseFuncEmpty                  = "ESv8-016"
	ErrCodeIndexBuilderNoMappingFound                = "ESv8-017"
	ErrCodeIndexAliasesNoAliasFound                  = "ESv8-018"
	ErrCodeIndexBuilderInvalidModel                  = "ESv8-019"
	ErrCodeIndexBuilderInvalidModelType              = "ESv8-020"
	ErrCodeIndexBuilderMappingSchemaNotExpected      = "ESv8-021"
	ErrCodeIndexBuilderMappingExistentFieldsModified = "ESv8-022"
	ErrCodeIndexBuilderAliasAndIndexEmpty            = "ESv8-023"
	ErrCodeIndexBuilderModelEmpty                    = "ESv8-024"
	ErrCodeIndexBuilderNoWriteIndexForAlias          = "ESv8-025"
	ErrCodeExecuteResponseProcessing                 = "ESv8-026"
	ErrCodeExecute                                   = "ESv8-027"
	ErrCodeDocUpdate                                 = "ESv8-028"
)

var (
	ErrNewClient = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeNewClient, "es: new client").Wrap(cause).Err()
	}
	ErrDocIndex = func(ctx context.Context, cause error, index, id string) error {
		return kit.NewAppErrBuilder(ErrCodeDocIndex, "es: document index").C(ctx).Wrap(cause).F(kit.KV{"index": index, "id": id}).Err()
	}
	ErrDocUpdate = func(ctx context.Context, cause error, index, id string) error {
		return kit.NewAppErrBuilder(ErrCodeDocUpdate, "es: document update").C(ctx).Wrap(cause).F(kit.KV{"index": index, "id": id}).Err()
	}
	ErrDocDelete = func(ctx context.Context, cause error, index, id string) error {
		return kit.NewAppErrBuilder(ErrCodeDocDelete, "es: document delete").C(ctx).Wrap(cause).F(kit.KV{"index": index, "id": id}).Err()
	}
	ErrDocExists = func(ctx context.Context, cause error, index, id string) error {
		return kit.NewAppErrBuilder(ErrCodeDocExists, "es: document exists").C(ctx).Wrap(cause).F(kit.KV{"index": index, "id": id}).Err()
	}
	ErrIndexRefresh = func(ctx context.Context, cause error, index string) error {
		return kit.NewAppErrBuilder(ErrCodeIndexRefresh, "es: index exists").C(ctx).Wrap(cause).F(kit.KV{"index": index}).Err()
	}
	ErrIndexUpdateAliases = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeIndexUpdateAliases, "es: update aliases").C(ctx).Wrap(cause).Err()
	}
	ErrIndexExists = func(ctx context.Context, cause error, index string) error {
		return kit.NewAppErrBuilder(ErrCodeIndexExists, "es: index exists").C(ctx).Wrap(cause).F(kit.KV{"index": index}).Err()
	}
	ErrIndexCreate = func(ctx context.Context, cause error, index string) error {
		return kit.NewAppErrBuilder(ErrCodeIndexCreate, "es: index create").C(ctx).Wrap(cause).F(kit.KV{"index": index}).Err()
	}
	ErrSearch = func(ctx context.Context, cause error, index string) error {
		return kit.NewAppErrBuilder(ErrCodeSearch, "es: search").C(ctx).Wrap(cause).F(kit.KV{"index": index}).Err()
	}
	ErrIndexPutMapping = func(ctx context.Context, cause error, index string) error {
		return kit.NewAppErrBuilder(ErrCodeIndexPutMapping, "es: put mapping").C(ctx).Wrap(cause).F(kit.KV{"index": index}).Err()
	}
	ErrIndexGetMapping = func(ctx context.Context, cause error, index string) error {
		return kit.NewAppErrBuilder(ErrCodeIndexGetMapping, "es: get mapping").C(ctx).Wrap(cause).F(kit.KV{"index": index}).Err()
	}
	ErrIndexGetIndices = func(ctx context.Context, cause error, alias string) error {
		return kit.NewAppErrBuilder(ErrCodeIndexGetIndices, "es: get by alias").C(ctx).Wrap(cause).F(kit.KV{"alias": alias}).Err()
	}
	ErrIndexBuilderNoMappingFound = func(ctx context.Context, index string) error {
		return kit.NewAppErrBuilder(ErrCodeIndexBuilderNoMappingFound, "es index builder: no mapping found").C(ctx).F(kit.KV{"index": index}).Err()
	}
	ErrIndexAliasesNoAliasFound = func(ctx context.Context, alias string) error {
		return kit.NewAppErrBuilder(ErrCodeIndexAliasesNoAliasFound, "es: no alias found").C(ctx).F(kit.KV{"alias": alias}).Err()
	}
	ErrExecuteFuncEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeExecuteFuncEmpty, "es: empty function").C(ctx).Err()
	}
	ErrExecuteResponseFuncEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeExecuteResponseFuncEmpty, "es: empty response function").C(ctx).Err()
	}
	ErrExecuteResponseCode = func(ctx context.Context, code int, body string) error {
		return kit.NewAppErrBuilder(ErrCodeExecuteResponseCode, "es: error status code %d %s", code, body).C(ctx).Err()
	}
	ErrExecuteResponseProcessing = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeExecuteResponseProcessing, "es: response processiong failed").Wrap(cause).C(ctx).Err()
	}
	ErrExecute = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeExecute, "es: execution failed").Wrap(cause).C(ctx).Err()
	}
	ErrIndexDelete = func(ctx context.Context, cause error, indices []string) error {
		return kit.NewAppErrBuilder(ErrCodeIndexDelete, "es: index delete").C(ctx).Wrap(cause).F(kit.KV{"indices": indices}).Err()
	}
	ErrIndexBuilderInvalidModel = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeIndexBuilderInvalidModel, "es index builder: invalid model, check tags").C(ctx).Err()
	}
	ErrIndexBuilderInvalidModelType = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeIndexBuilderInvalidModelType, "es index builder: model must be pointer of struct").C(ctx).Err()
	}
	ErrIndexBuilderMappingSchemaNotExpected = func(ctx context.Context, cause error, index string) error {
		return kit.NewAppErrBuilder(ErrCodeIndexBuilderMappingSchemaNotExpected, "es index builder: mapping schema not expected %s", index).C(ctx).F(kit.KV{"index": index}).Err()
	}
	ErrIndexBuilderMappingExistentFieldsModified = func(ctx context.Context, index string, fields []string) error {
		return kit.NewAppErrBuilder(ErrCodeIndexBuilderMappingExistentFieldsModified, "es index builder: doesn't allowed to change mapping for existent fields.").C(ctx).F(kit.KV{"index": index, "fields": fields}).Err()
	}
	ErrIndexBuilderAliasAndIndexEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeIndexBuilderAliasAndIndexEmpty, "es index builder: neither alias name nor index name specified").C(ctx).Err()
	}
	ErrIndexBuilderModelEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeIndexBuilderModelEmpty, "es index builder: model not specified").C(ctx).Err()
	}
	ErrIndexBuilderNoWriteIndexForAlias = func(ctx context.Context, alias string) error {
		return kit.NewAppErrBuilder(ErrCodeIndexBuilderNoWriteIndexForAlias, "es index builder: no write index").F(kit.KV{"alias": alias}).C(ctx).Err()
	}
)
