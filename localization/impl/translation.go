package impl

import (
	"context"
	"fmt"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/localization"
	"math/rand"
)

type translationSvcImpl struct {
	storage            localization.DataTranslationStorage
	logger             kit.CLoggerFunc
	supportedLanguages []string
}

func NewDataTranslationService(storage localization.DataTranslationStorage, logger kit.CLoggerFunc, supportedLanguages []string) localization.DataTranslationService {
	return &translationSvcImpl{
		storage:            storage,
		logger:             logger,
		supportedLanguages: supportedLanguages,
	}
}

func (t *translationSvcImpl) l() kit.CLogger {
	return t.logger().Cmp("translation-svc")
}

func (t *translationSvcImpl) Merge(ctx context.Context, translations ...*localization.TranslationByKey) ([]string, error) {
	t.l().C(ctx).Mth("merge").Dbg()

	// Validate translation
	for i, tr := range translations {
		modified, err := t.validateAndPopulate(ctx, tr)
		if err != nil {
			return nil, err
		}
		translations[i] = modified
	}

	// merge translations
	err := t.storage.Merge(ctx, translations...)
	if err != nil {
		return nil, err
	}

	return kit.Map(translations, func(tr *localization.TranslationByKey) string { return tr.Key }), nil
}

func (t *translationSvcImpl) MustGet(ctx context.Context, key string, lang string) (*localization.Translation, error) {
	t.l().C(ctx).Mth("must-get").F(kit.KV{"key": key, "lang": lang}).Dbg()

	// validate lang
	err := t.validateLang(ctx, lang)
	if err != nil {
		return nil, err
	}

	translation, err := t.Get(ctx, key, lang)
	if err != nil {
		return nil, err
	}
	if translation == nil {
		return nil, localization.ErrTranslationNotFound(ctx)
	}

	return translation, nil
}

func (t *translationSvcImpl) MustGetByKeys(ctx context.Context, lang string, keys []string) (map[string]*localization.Translation, error) {
	t.l().C(ctx).Mth("must-by-keys").F(kit.KV{"lang": lang}).Dbg()

	// validate lang
	err := t.validateLang(ctx, lang)
	if err != nil {
		return nil, err
	}

	r := make(map[string]*localization.Translation)

	if len(keys) == 0 {
		return r, nil
	}

	// make keys distinct
	keys = kit.Strings(keys).Distinct()

	// Create criteria for GetMany
	criteria := localization.GetTranslationsCriteria{
		Keys: keys,
		Lang: lang,
	}

	// retrieve from storage
	rs, err := t.storage.Search(ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request:       criteria,
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	if err != nil {
		return nil, err
	}

	// Check if all requested keys were found
	if len(rs.Items) != len(keys) {
		return nil, localization.ErrTranslationNotFound(ctx)
	}

	// group by key
	for _, item := range rs.Items {
		r[item.Key] = &localization.Translation{
			Key:      item.Key,
			Value:    item.Translations[lang],
			Category: item.Category,
		}
	}

	return r, nil
}

func (t *translationSvcImpl) Get(ctx context.Context, key string, lang string) (*localization.Translation, error) {
	t.l().C(ctx).Mth("get").F(kit.KV{"key": key, "lang": lang}).Dbg()

	if key == "" {
		return nil, localization.ErrTranslationKeyEmpty(ctx)
	}

	// validate lang
	err := t.validateLang(ctx, lang)
	if err != nil {
		return nil, err
	}

	translation, err := t.storage.Get(ctx, key, lang)
	if err != nil {
		return nil, err
	}

	return translation, nil
}

func (t *translationSvcImpl) GenKey(category string) string {
	return t.BuildKey(category, fmt.Sprintf("%d", rand.Int63()))
}

func (t *translationSvcImpl) BuildKey(category, id string) string {
	return fmt.Sprintf("%s.%s", category, id)
}

func (t *translationSvcImpl) TranslateObjects(ctx context.Context, lang string, objs ...localization.Translatable) (localization.TranslationKV, error) {
	t.l().C(ctx).Mth("trans-obj").F(kit.KV{"lang": lang}).Dbg()

	err := t.validateLang(ctx, lang)
	if err != nil {
		return nil, err
	}

	// gather keys
	var keys []string
	kit.ForAll(objs, func(obj localization.Translatable) {
		keys = append(keys, obj.GetTranslatedKeys()...)
	})

	transByKeys, err := t.MustGetByKeys(ctx, lang, keys)
	if err != nil {
		return nil, err
	}

	r := make(localization.TranslationKV, len(transByKeys))
	for k, v := range transByKeys {
		r[k] = v.Value
	}

	return r, nil

}

func (t *translationSvcImpl) MustTranslateObjects(ctx context.Context, lang string, objs ...localization.Translatable) localization.TranslationKV {
	l := t.l().C(ctx).Mth("must-trans-obj").F(kit.KV{"lang": lang}).Dbg()

	r, err := t.TranslateObjects(ctx, lang, objs...)
	if err != nil {
		// log error
		l.E(err).St().Err()

		// get requested keys
		var keys []string
		kit.ForAll(objs, func(obj localization.Translatable) {
			keys = append(keys, obj.GetTranslatedKeys()...)
		})

		// return KV with keys as values
		r = make(localization.TranslationKV, len(objs))
		kit.ForAll(keys, func(key string) {
			r[key] = key
		})

	}

	return r

}

func (t *translationSvcImpl) Search(ctx context.Context, cr *kit.PagingRequestG[localization.GetTranslationsCriteria]) (*kit.PagingResponseG[localization.TranslationByKey], error) {
	t.l().C(ctx).Mth("search").Dbg()

	if cr.Request.Lang != "" {
		err := t.validateLang(ctx, cr.Request.Lang)
		if err != nil {
			return nil, err
		}
	}

	return t.storage.Search(ctx, cr)
}

func (t *translationSvcImpl) validateAndPopulate(ctx context.Context, translation *localization.TranslationByKey) (*localization.TranslationByKey, error) {

	if translation.Category == "" {
		return nil, localization.ErrTranslationCategoryEmpty(ctx)
	}

	if len(translation.Translations) == 0 {
		return nil, localization.ErrTranslationLangEmpty(ctx)
	}

	for lang, val := range translation.Translations {

		err := t.validateLang(ctx, lang)
		if err != nil {
			return nil, err
		}

		// validate lang
		if val == "" {
			return nil, localization.ErrTranslationInvalid(ctx)
		}
	}

	if translation.Key == "" {
		translation.Key = t.GenKey(translation.Category)
	}

	return translation, nil
}

func (t *translationSvcImpl) validateLang(ctx context.Context, lang string) error {

	if lang == "" {
		return localization.ErrTranslationLangEmpty(ctx)
	}

	if !kit.IsValidISO639_1(lang) {
		return localization.ErrNotValidLang(ctx)
	}

	// check language is supported
	if !kit.Strings(t.supportedLanguages).Contains(lang) {
		return localization.ErrNotSupportedLang(ctx)
	}

	return nil
}
