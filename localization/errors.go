package localization

import (
	"context"
	"github.com/mikhailbolshakov/kit"
)

const (
	ErrCodeLocTranslationKeyEmpty           = "LOC-001"
	ErrCodeLocTranslationLangEmpty          = "LOC-002"
	ErrCodeLocTranslationCategoryEmpty      = "LOC-003"
	ErrCodeLocTranslationInvalid            = "LOC-004"
	ErrCodeTranslationNotFound              = "LOC-005"
	ErrCodeNotSupportedLang                 = "LOC-006"
	ErrCodeNotValidLang                     = "LOC-007"
	ErrCodeDataTranslationDbMerge           = "LOC-008"
	ErrCodeDataTranslationDbGetMany         = "LOC-009"
	ErrCodeDataTranslationEnsureTransFields = "LOC-010"
)

var (
	ErrTranslationKeyEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeLocTranslationKeyEmpty, "translation key is empty").C(ctx).Err()
	}
	ErrTranslationLangEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeLocTranslationLangEmpty, "translation language is empty").C(ctx).Err()
	}
	ErrTranslationCategoryEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeLocTranslationCategoryEmpty, "translation category is empty").C(ctx).Err()
	}
	ErrTranslationInvalid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeLocTranslationInvalid, "translation is invalid").C(ctx).Err()
	}
	ErrTranslationNotFound = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeTranslationNotFound, "translation not found").C(ctx).Err()
	}
	ErrNotSupportedLang = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeNotSupportedLang, "lang isn't supported").C(ctx).Err()
	}
	ErrNotValidLang = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeNotValidLang, "lang isn't valid").C(ctx).Err()
	}
	ErrDataTranslationDbMerge = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeDataTranslationDbMerge, "translation storage: merge").Wrap(cause).C(ctx).Err()
	}
	ErrDataTranslationDbGetMany = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeDataTranslationDbGetMany, "translation storage: get many").Wrap(cause).C(ctx).Err()
	}
	ErrDataTranslationEnsureTransFields = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeDataTranslationEnsureTransFields, "translation storage: ensure language fields").Wrap(cause).C(ctx).Err()
	}
)
