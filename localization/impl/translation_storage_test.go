//go:build integration

package impl

import (
	"fmt"
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/localization"
	kitPg "github.com/mikhailbolshakov/kit/storages/pg"
	kitRedis "github.com/mikhailbolshakov/kit/storages/redis"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	connectionString = "postgres://product:product@localhost:55432/test?connect_timeout=10&sslmode=disable"

	redisConfig = &kitRedis.Config{
		Host: "localhost",
		Port: "6379",
	}
)

type dataTranslationStorageTestSuite struct {
	kit.Suite
	logger  kit.CLoggerFunc
	storage localization.DataTranslationStorage
	db      *kitPg.Storage
	redis   *kitRedis.Redis
}

func (s *dataTranslationStorageTestSuite) SetupSuite() {
	s.logger = func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) }
	s.Suite.Init(s.logger)

	var err error
	s.db, err = kitPg.Open(&kitPg.DbConfig{ConnectionString: connectionString}, s.logger)
	s.NoError(err)

	s.redis, err = kitRedis.Open(s.Ctx, redisConfig, s.logger)
	s.NoError(err)
}

func (s *dataTranslationStorageTestSuite) TearDownSuite() {
	s.db.Close()
}

func (s *dataTranslationStorageTestSuite) SetupTest() {
	s.storage = NewDataTranslationStorage(localization.DataStorageConfig{
		DB:                 s.db,
		Redis:              s.redis,
		SupportedLanguages: []string{"en", "sr", "ru"},
		DefaultLang:        "en",
		Logger:             s.logger,
	})
}

func TestDataTranslationSuite(t *testing.T) {
	suite.Run(t, new(dataTranslationStorageTestSuite))
}

func (s *dataTranslationStorageTestSuite) Test_NewLangField() {
	translation := s.getTranslation()

	// Save translation
	err := s.storage.Merge(s.Ctx, translation)
	s.NoError(err)

	s.storage = NewDataTranslationStorage(localization.DataStorageConfig{
		DB:                 s.db,
		Redis:              s.redis,
		SupportedLanguages: []string{"en", "sr", "ru", "de"},
		DefaultLang:        "en",
		Logger:             s.logger,
	})
	s.NoError(s.storage.EnsureTranslationFields(s.Ctx))

	translation.Translations = map[string]string{"de": "yaya"}

	// Save translation
	err = s.storage.Merge(s.Ctx, translation)
	s.NoError(err)

	// Get translation by key and lang
	retrieved, err := s.storage.Get(s.Ctx, translation.Key, "de")
	s.NoError(err)
	s.NotEmpty(retrieved.Value)

	retrieved, err = s.storage.Get(s.Ctx, translation.Key, "en")
	s.NoError(err)
	s.NotEmpty(retrieved.Value)
}

func (s *dataTranslationStorageTestSuite) TestMerge_Update() {
	// Create a test translation
	translation := s.getTranslation()

	// Save translation
	err := s.storage.Merge(s.Ctx, translation)
	s.NoError(err)

	// Get translation by key and lang
	retrieved, err := s.storage.Get(s.Ctx, translation.Key, "en")
	s.NoError(err)
	s.NotNil(retrieved)
	s.Equal(translation.Key, retrieved.Key)
	s.Equal(translation.Category, retrieved.Category)
	s.Equal(translation.Translations["en"], retrieved.Value)

	// Update the translation
	updatedTranslation := &localization.TranslationByKey{
		Key:          translation.Key,
		Category:     translation.Category,
		Translations: map[string]string{"en": "Updated value"},
	}

	// Save updated translation
	err = s.storage.Merge(s.Ctx, updatedTranslation)
	s.NoError(err)

	// Get updated translation
	updated, err := s.storage.Get(s.Ctx, translation.Key, "en")
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updatedTranslation.Translations["en"], updated.Value)
}

func (s *dataTranslationStorageTestSuite) TestMerge_Multiple() {
	// Create test translations
	translation1 := s.getTranslation()
	translation2 := s.getTranslation()

	// Save translations
	err := s.storage.Merge(s.Ctx, translation1, translation2)
	s.NoError(err)

	// Get and check first translation
	actual1, err := s.storage.Get(s.Ctx, translation1.Key, "en")
	s.NoError(err)
	s.NotNil(actual1)
	s.Equal(translation1.Key, actual1.Key)
	s.Equal(translation1.Category, actual1.Category)
	s.Equal(translation1.Translations["en"], actual1.Value)

	// Get and check second translation
	actual2, err := s.storage.Get(s.Ctx, translation2.Key, "en")
	s.NoError(err)
	s.NotNil(actual2)
	s.Equal(translation2.Key, actual2.Key)
	s.Equal(translation2.Category, actual2.Category)
	s.Equal(translation2.Translations["en"], actual2.Value)

	// Get both translations using Search
	searchResults, err := s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Lang: "en",
			Keys: []string{translation1.Key, translation2.Key},
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 2)
}

func (s *dataTranslationStorageTestSuite) Test_DefaultLanguage() {
	// Create test translations
	translation1 := s.getTranslation()
	translation2 := s.getTranslation()
	translation2.Translations["ru"] = "что-то"

	// Save translations
	err := s.storage.Merge(s.Ctx, translation1, translation2)
	s.NoError(err)

	// Get and check first translation
	actual1, err := s.storage.Get(s.Ctx, translation1.Key, "ru")
	s.NoError(err)
	s.NotNil(actual1)
	s.Equal(translation1.Key, actual1.Key)
	s.Equal(translation1.Category, actual1.Category)
	s.Equal(translation1.Translations["en"], actual1.Value)

	// Get and check second translation
	actual2, err := s.storage.Get(s.Ctx, translation2.Key, "ru")
	s.NoError(err)
	s.NotNil(actual2)
	s.Equal(translation2.Key, actual2.Key)
	s.Equal(translation2.Category, actual2.Category)
	s.Equal(translation2.Translations["ru"], actual2.Value)

	// Get both translations using Search
	searchResults, err := s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Lang: "ru",
			Keys: []string{translation1.Key, translation2.Key},
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 2)
	s.NotEmpty(searchResults.Items[0].Translations["ru"])
	s.NotEmpty(searchResults.Items[1].Translations["ru"])
}

func (s *dataTranslationStorageTestSuite) TestMerge_MultipleLanguages() {
	// Create a test translation with multiple languages
	translation := &localization.TranslationByKey{
		Key:      "test.multilang." + kit.NewId(),
		Category: fmt.Sprintf("test-multilang-%s", kit.NewRandString()),
		Translations: map[string]string{
			"en": "English value",
			"sr": "Serbian value",
			"ru": "Russian value",
		},
	}

	// Save translation
	err := s.storage.Merge(s.Ctx, translation)
	s.NoError(err)

	// Get and check English translation
	enTranslation, err := s.storage.Get(s.Ctx, translation.Key, "en")
	s.NoError(err)
	s.NotNil(enTranslation)
	s.Equal(translation.Key, enTranslation.Key)
	s.Equal(translation.Translations["en"], enTranslation.Value)

	// Get and check Serbian translation
	srTranslation, err := s.storage.Get(s.Ctx, translation.Key, "sr")
	s.NoError(err)
	s.NotNil(srTranslation)
	s.Equal(translation.Key, srTranslation.Key)
	s.Equal(translation.Translations["sr"], srTranslation.Value)

	// Get and check Russian translation
	ruTranslation, err := s.storage.Get(s.Ctx, translation.Key, "ru")
	s.NoError(err)
	s.NotNil(ruTranslation)
	s.Equal(translation.Key, ruTranslation.Key)
	s.Equal(translation.Translations["ru"], ruTranslation.Value)

	// Get the translation using Search
	searchResults, err := s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Keys: []string{translation.Key},
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 1)
	s.Equal(translation.Key, searchResults.Items[0].Key)
	s.Equal(translation.Category, searchResults.Items[0].Category)
	s.Equal(translation.Translations["en"], searchResults.Items[0].Translations["en"])
	s.Equal(translation.Translations["sr"], searchResults.Items[0].Translations["sr"])
	s.Equal(translation.Translations["ru"], searchResults.Items[0].Translations["ru"])
}

func (s *dataTranslationStorageTestSuite) TestGet_NonExistent() {
	// Try to get a non-existent translation
	translation, err := s.storage.Get(s.Ctx, "non.existent.key."+kit.NewId(), "en")
	s.NoError(err)
	s.Nil(translation)
}

func (s *dataTranslationStorageTestSuite) TestSearch_ByCategory() {
	// Create test translations with specific categories
	category1 := "test-category-" + kit.NewRandString()
	category2 := "test-category-" + kit.NewRandString()

	translation1 := &localization.TranslationByKey{
		Key:          "test.cat1.1." + kit.NewId(),
		Category:     category1,
		Translations: map[string]string{"en": "Category 1 Item 1"},
	}

	translation2 := &localization.TranslationByKey{
		Key:          "test.cat1.2." + kit.NewId(),
		Category:     category1,
		Translations: map[string]string{"en": "Category 1 Item 2"},
	}

	translation3 := &localization.TranslationByKey{
		Key:          "test.cat2.1." + kit.NewId(),
		Category:     category2,
		Translations: map[string]string{"en": "Category 2 Item 1"},
	}

	translation4 := &localization.TranslationByKey{
		Key:          "test.cat2.2." + kit.NewId(),
		Category:     category2,
		Translations: map[string]string{"en": "Category 2 Item 2"},
	}

	// Save translations
	err := s.storage.Merge(s.Ctx, translation1, translation2, translation3, translation4)
	s.NoError(err)

	// Search by category1
	searchResults, err := s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Categories: []string{category1},
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 2)

	// Search by category2
	searchResults, err = s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Categories: []string{category2},
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 2)

	// Search by both categories
	searchResults, err = s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Categories: []string{category1, category2},
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 4)
}

func (s *dataTranslationStorageTestSuite) TestSearch_ByFullTextSearch() {
	// Create test translations with specific content for FTS
	category := kit.NewRandString()
	translation1 := &localization.TranslationByKey{
		Key:          "test.fts.apple." + kit.NewId(),
		Category:     category,
		Translations: map[string]string{"en": "Red apple healthy fruit"},
	}

	translation2 := &localization.TranslationByKey{
		Key:          "test.fts.banana." + kit.NewId(),
		Category:     category,
		Translations: map[string]string{"en": "Yellow banana fruit tropical"},
	}

	translation3 := &localization.TranslationByKey{
		Key:          "test.fts.carrot." + kit.NewId(),
		Category:     category,
		Translations: map[string]string{"en": "Orange carrot healthy vegetable"},
	}

	// Save translations
	err := s.storage.Merge(s.Ctx, translation1, translation2, translation3)
	s.NoError(err)

	// Search by "fruit" - should find apple and banana
	searchResults, err := s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Categories: []string{category},
			Fts:        "fruit",
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 2)

	// Search by "healthy" - should find apple and carrot
	searchResults, err = s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Categories: []string{category},
			Fts:        "healthy",
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 2)

	// Search by "orange" - should find only carrot
	searchResults, err = s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Categories: []string{category},
			Fts:        "orange",
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 1)
	s.Equal(translation3.Key, searchResults.Items[0].Key)
}

func (s *dataTranslationStorageTestSuite) TestSearch_WithPaging() {
	// Create multiple test translations
	category := "paging" + kit.NewRandString()
	var translations []*localization.TranslationByKey
	for i := 0; i < 10; i++ {
		translation := &localization.TranslationByKey{
			Key:          fmt.Sprintf("test.paging.%d.%s", i, kit.NewId()),
			Category:     category,
			Translations: map[string]string{"en": fmt.Sprintf("Paging test item %d", i)},
		}
		translations = append(translations, translation)
	}

	// Save translations
	err := s.storage.Merge(s.Ctx, translations...)
	s.NoError(err)

	// Search with paging - first page (3 items)
	searchResults, err := s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		PagingRequest: kit.PagingRequest{
			Index: 1,
			Size:  3,
		},
		Request: localization.GetTranslationsCriteria{
			Categories: []string{category},
		},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 3)

	// Search with paging - second page (3 more items)
	searchResults, err = s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		PagingRequest: kit.PagingRequest{
			Index: 3,
			Size:  3,
		},
		Request: localization.GetTranslationsCriteria{
			Categories: []string{category},
		},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 3)

	// Verify we get all 10 items when using RetrieveAll
	searchResults, err = s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Categories: []string{category},
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 10)
}

func (s *dataTranslationStorageTestSuite) TestMultiLanguageFTS() {
	// Create a translation with multiple languages for FTS testing
	translation := &localization.TranslationByKey{
		Key:      "test.multilang.fts." + kit.NewId(),
		Category: "test-fts-multilang",
		Translations: map[string]string{
			"en": "English apple fruit",
			"sr": "Serbian jabuka voće",
			"ru": "Russian яблоко фрукты",
		},
	}

	// Save translation
	err := s.storage.Merge(s.Ctx, translation)
	s.NoError(err)

	// Search for English term
	searchResults, err := s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Keys: []string{translation.Key},
			Fts:  "apple",
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 1)
	s.Equal(translation.Key, searchResults.Items[0].Key)

	// Search for Serbian term
	searchResults, err = s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Keys: []string{translation.Key},
			Fts:  "jabuka",
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 1)
	s.Equal(translation.Key, searchResults.Items[0].Key)

	// Search for Russian term
	searchResults, err = s.storage.Search(s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request: localization.GetTranslationsCriteria{
			Keys: []string{translation.Key},
			Fts:  "яблоко",
		},
		PagingRequest: kit.PagingRequest{Skip: true},
	})
	s.NoError(err)
	s.NotNil(searchResults)
	s.Len(searchResults.Items, 1)
	s.Equal(translation.Key, searchResults.Items[0].Key)
}

func (s *dataTranslationStorageTestSuite) getTranslation() *localization.TranslationByKey {
	return &localization.TranslationByKey{
		Key:          "test.translation." + kit.NewId(),
		Category:     fmt.Sprintf("test-%s", kit.NewRandString()),
		Translations: map[string]string{"en": "some translation"},
	}
}
