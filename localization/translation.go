package localization

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	kitPg "github.com/mikhailbolshakov/kit/storages/pg"
	kitRedis "github.com/mikhailbolshakov/kit/storages/redis"
)

// Config defines language configuration settings
type Config struct {
	Languages struct {
		Default   string   // Default language code to use
		Supported []string // List of supported language codes
	}
}

type TranslationKV map[string]string

type TranslationByKey struct {
	Key          string        // Key uniquely identifies the translation
	Translations TranslationKV // Translations by language
	Category     string        // Category groups related translations
}

type Translation struct {
	Key      string // Key uniquely identifies the translation
	Value    string // Value translation to the requested language
	Category string // Category groups related translations
}

type GetManyByLangRq struct {
	Keys       []string // Keys list of translation keys to retrieve
	Categories []string // Categories list of categories to filter by
}

type GetTranslationsCriteria struct {
	Lang       string   // Lang specifies the language code to filter by
	Keys       []string // Keys list of translation keys to retrieve
	Categories []string // Categories list of categories to filter by
	Fts        string   // Fts full text search
}

// Translatable this interface must be implemented by any struct supporting translations
type Translatable interface {
	// GetTranslatedKeys retrieves translation keys
	GetTranslatedKeys() []string
}

type DataTranslationService interface {
	// Merge validates and stores multiple translations. Returns merged keys
	Merge(ctx context.Context, t ...*TranslationByKey) ([]string, error)
	// Get retrieves a translation by key and language
	Get(ctx context.Context, key string, lang string) (*Translation, error)
	// MustGet retrieves a translation by key and language, errors if not found
	MustGet(ctx context.Context, key string, lang string) (*Translation, error)
	// MustGetByKeys retrieves multiple translations by keys for a language, errors if any not found
	MustGetByKeys(ctx context.Context, lang string, keys []string) (map[string]*Translation, error)
	// Search retrieves translations filtered by criteria
	Search(ctx context.Context, cr *kit.PagingRequestG[GetTranslationsCriteria]) (*kit.PagingResponseG[TranslationByKey], error)
	// GenKey generates a new translation key for a category
	GenKey(category string) string
	// BuildKey creates a translation key from category and ID
	BuildKey(category, id string) string
	// TranslateObjects translates the provided Translatable objects and returns a map of translation keys to their values.
	TranslateObjects(ctx context.Context, lang string, objs ...Translatable) (TranslationKV, error)
	// MustTranslateObjects translates the provided Translatable objects and returns a map of translation keys to their values.
	// if an error occurs, translation keys are returned instead of translation values
	MustTranslateObjects(ctx context.Context, lang string, objs ...Translatable) TranslationKV
}

type DataStorageConfig struct {
	DB                 *kitPg.Storage
	Redis              *kitRedis.Redis
	SupportedLanguages []string
	DefaultLang        string
	Logger             kit.CLoggerFunc
}

type DataTranslationStorage interface {
	// Merge stores multiple translations
	Merge(ctx context.Context, t ...*TranslationByKey) error
	// Get retrieves a translation by key and language
	Get(ctx context.Context, key string, lang string) (*Translation, error)
	// Search retrieves translations matching the given criteria
	Search(ctx context.Context, cr *kit.PagingRequestG[GetTranslationsCriteria]) (*kit.PagingResponseG[TranslationByKey], error)
	// EnsureTranslationFields checks if the data_translations table has the necessary
	// language fields and adds any missing ones
	EnsureTranslationFields(ctx context.Context) error
}
