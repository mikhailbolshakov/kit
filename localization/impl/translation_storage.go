package impl

import (
	"context"
	"fmt"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/localization"
	kitPg "github.com/mikhailbolshakov/kit/storages/pg"
	kitRedis "github.com/mikhailbolshakov/kit/storages/redis"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type dataTranslationStorageImpl struct {
	supportedLanguages []string
	defaultLang        string
	db                 *kitPg.Storage
	redis              *kitRedis.Redis
	logger             kit.CLoggerFunc
}

func NewDataTranslationStorage(cfg localization.DataStorageConfig) localization.DataTranslationStorage {
	return &dataTranslationStorageImpl{
		supportedLanguages: cfg.SupportedLanguages,
		db:                 cfg.DB,
		redis:              cfg.Redis,
		logger:             cfg.Logger,
		defaultLang:        cfg.DefaultLang,
	}
}

func (s *dataTranslationStorageImpl) l() kit.CLogger {
	return s.logger().Cmp("translation-storage")
}

func (s *dataTranslationStorageImpl) EnsureTranslationFields(ctx context.Context) error {
	l := s.l().Mth("ensure-fields").Dbg()

	// Get existing columns
	var columns []string
	err := s.db.Instance.Raw(`SELECT column_name FROM information_schema.columns WHERE table_name = 'data_translations'`).
		Scan(&columns).Error
	if err != nil {
		return localization.ErrDataTranslationEnsureTransFields(ctx, err)
	}

	// Convert to map
	columnsSet := kit.Strings(columns).ToMap()

	// Check if any supported languages are missing
	var missingLang []string
	for _, lang := range s.supportedLanguages {
		if _, ok := columnsSet[lang]; !ok {
			missingLang = append(missingLang, lang)
		}
	}

	// Add missing language columns if any
	for _, lang := range missingLang {
		err = s.db.Instance.Exec(fmt.Sprintf("ALTER TABLE data_translations ADD COLUMN %s TEXT", lang)).Error
		if err != nil {
			return localization.ErrDataTranslationEnsureTransFields(ctx, err)
		}
	}

	l.Dbg("added: %v", missingLang)

	return nil
}

func (s *dataTranslationStorageImpl) Get(ctx context.Context, key, lang string) (*localization.Translation, error) {
	s.l().Mth("get").C(ctx).F(kit.KV{"key": key, "lang": lang}).Dbg()

	if key == "" {
		return nil, localization.ErrTranslationKeyEmpty(ctx)
	}

	if lang == "" {
		return nil, localization.ErrTranslationLangEmpty(ctx)
	}

	type Dto struct {
		Key      string
		Category string
		Value    string
	}

	dto := &Dto{}

	res := s.db.Instance.Scopes(kitPg.Single()).
		Table("data_translations").
		Select(fmt.Sprintf("key, category, COALESCE(NULLIF(%s, ''), %s) as value", lang, s.defaultLang)).
		Where("key = ?", key).
		Find(dto)

	if res.Error != nil {
		return nil, localization.ErrDataTranslationDbMerge(ctx, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &localization.Translation{
		Key:      dto.Key,
		Category: dto.Category,
		Value:    dto.Value,
	}, nil
}

func (s *dataTranslationStorageImpl) Merge(ctx context.Context, translations ...*localization.TranslationByKey) error {
	s.l().Mth("merge").C(ctx).Dbg()

	if len(translations) == 0 {
		return nil
	}

	return s.db.Instance.Transaction(func(tx *gorm.DB) error {

		// Process each translation
		for _, translation := range translations {

			// Create a map for the values to insert/update
			values := map[string]any{
				"key":      translation.Key,
				"category": translation.Category,
			}

			// Add language values from translations
			for _, lang := range s.supportedLanguages {
				if value, ok := translation.Translations[lang]; ok {
					values[lang] = value
				}
			}

			// Create a slice of update columns (excluding the primary key)
			updateColumns := make([]clause.Column, 0, len(values)-1)
			updateValues := make([]string, 0, len(values)-1)
			for col := range values {
				if col != "key" { // Skip primary key for updates
					updateColumns = append(updateColumns, clause.Column{Name: col})
					updateValues = append(updateValues, col)
				}
			}

			// Use GORM's Clauses for upsert operation
			result := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "key"}},
				DoUpdates: clause.AssignmentColumns(updateValues),
			}).Table("data_translations").Create(values)

			if result.Error != nil {
				return localization.ErrDataTranslationDbMerge(ctx, result.Error)
			}

			// Then update the FTS vector with a separate query
			// This ensures we have access to all language values (including ones not being updated)
			var ftsExpressions []string
			for _, lang := range s.supportedLanguages {
				ftsExpressions = append(ftsExpressions, "COALESCE("+lang+", '')")
			}

			// Use GORM's Expression to safely build the SQL function call
			ftsExpression := gorm.Expr("to_tsvector('simple', ?)",
				gorm.Expr(strings.Join(ftsExpressions, " || ' ' || ")))

			// Update the FTS vector using GORM's query builder
			ftsRes := tx.Table("data_translations").
				Where("key = ?", translation.Key).
				Update("fts", ftsExpression)

			if ftsRes.Error != nil {
				return localization.ErrDataTranslationDbMerge(ctx, result.Error)
			}

		}

		return nil
	})
}

func (s *dataTranslationStorageImpl) Search(ctx context.Context, cr *kit.PagingRequestG[localization.GetTranslationsCriteria]) (*kit.PagingResponseG[localization.TranslationByKey], error) {
	s.l().Mth("search").C(ctx).Dbg()

	rs := &kit.PagingResponseG[localization.TranslationByKey]{
		PagingResponse: kit.PagingResponse{
			Limit: kitPg.PagingLimit(cr.Size),
		},
	}

	// Start building the query
	q := s.db.Instance.Table("data_translations")

	if len(cr.Request.Keys) > 0 {
		q = q.Where("key in (?)", cr.Request.Keys)
	}
	if len(cr.Request.Categories) > 0 {
		q = q.Where("category in (?)", cr.Request.Categories)
	}
	if cr.Request.Lang != "" {
		q = q.Select(fmt.Sprintf("key, category, COALESCE(NULLIF(%s, ''), %s) as %s", cr.Request.Lang, s.defaultLang, cr.Request.Lang))
	}
	if cr.Request.Fts != "" {
		ftsQ := kitPg.FtsQuery(cr.Request.Fts)
		if ftsQ != "" {
			q = q.Where("fts @@ to_tsquery('simple', ?)", ftsQ)
		}
	}

	// apply paging if requested
	if !cr.PagingRequest.Skip {
		if len(cr.SortBy) == 0 {
			cr.SortBy = []*kit.SortRequest{{
				Field: "key",
			}}
		}
		rs.PagingResponse.Limit = kitPg.PagingLimit(cr.Size)
		q = q.Scopes(kitPg.Paging(cr.PagingRequest))
	}

	// Execute the query with GORM's map scanning
	var results []map[string]any
	r := q.Find(&results)
	if r.Error != nil {
		return nil, localization.ErrDataTranslationDbGetMany(ctx, r.Error)
	}
	if r.RowsAffected == 0 {
		return nil, nil
	}

	rs.Items = s.toTransByKeysDomain(results)

	return rs, nil
}

func (s *dataTranslationStorageImpl) toTransByKeysDomain(dtos []map[string]any) []*localization.TranslationByKey {

	translations := make([]*localization.TranslationByKey, 0, len(dtos))
	for _, dto := range dtos {

		// Create a translation with the data from the map
		tr := &localization.TranslationByKey{
			Key:          dto["key"].(string),
			Translations: make(map[string]string),
			Category:     dto["category"].(string),
		}

		// Add the translation
		for _, lang := range s.supportedLanguages {
			if val, ok := dto[lang]; ok && val != nil && val.(string) != "" {
				tr.Translations[lang] = val.(string)
			}
		}

		translations = append(translations, tr)
	}

	return translations
}
