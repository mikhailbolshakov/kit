package impl

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.com/algmib/kit"
	"gitlab.com/algmib/kit/localization"
	kitMock "gitlab.com/algmib/kit/mocks"
	"strings"
	"testing"
)

type translationTestSuite struct {
	kit.Suite
	storage            *kitMock.LocalizationDataTranslationStorage
	logger             kit.CLoggerFunc
	translationSvc     localization.DataTranslationService
	supportedLanguages []string
}

func (s *translationTestSuite) SetupSuite() {
	s.logger = func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) }
	s.Suite.Init(s.logger)
	s.supportedLanguages = []string{"en", "fr"}
}

func (s *translationTestSuite) SetupTest() {
	s.storage = &kitMock.LocalizationDataTranslationStorage{}
	s.translationSvc = NewDataTranslationService(s.storage, s.logger, s.supportedLanguages)
}

func TestTranslationSuite(t *testing.T) {
	suite.Run(t, new(translationTestSuite))
}

func (s *translationTestSuite) Test_Merge_EmptyCategory() {
	translation := &localization.TranslationByKey{
		Key:          "test.key",
		Category:     "",
		Translations: map[string]string{"en": "Test value"},
	}
	result, err := s.translationSvc.Merge(s.Ctx, translation)
	s.Nil(result)
	s.Error(err)
	s.AssertAppErr(err, localization.ErrCodeLocTranslationCategoryEmpty)
}

func (s *translationTestSuite) Test_Merge_EmptyLang() {
	translation := &localization.TranslationByKey{
		Key:          "test.key",
		Category:     "test",
		Translations: map[string]string{},
	}
	result, err := s.translationSvc.Merge(s.Ctx, translation)
	s.Nil(result)
	s.Error(err)
	s.AssertAppErr(err, localization.ErrCodeLocTranslationLangEmpty)
}

func (s *translationTestSuite) Test_Merge_EmptyValue() {
	translation := &localization.TranslationByKey{
		Key:          "test.key",
		Category:     "test",
		Translations: map[string]string{"en": ""},
	}
	result, err := s.translationSvc.Merge(s.Ctx, translation)
	s.Nil(result)
	s.Error(err)
	s.AssertAppErr(err, localization.ErrCodeLocTranslationInvalid)
}

func (s *translationTestSuite) Test_Merge_UnsupportedLang() {
	translation := &localization.TranslationByKey{
		Key:          "test.key",
		Category:     "test",
		Translations: map[string]string{"es": "Test value"},
	}
	result, err := s.translationSvc.Merge(s.Ctx, translation)
	s.Nil(result)
	s.Error(err)
	s.AssertAppErr(err, localization.ErrCodeNotSupportedLang)
}

func (s *translationTestSuite) Test_Merge_GeneratesKey() {
	translation := &localization.TranslationByKey{
		Key:          "",
		Category:     "test",
		Translations: map[string]string{"en": "Test value"},
	}
	s.storage.On("Merge", s.Ctx, mock.MatchedBy(func(tt []*localization.TranslationByKey) bool {
		t := tt[0]
		return t.Translations["en"] == "Test value" && t.Category == "test" && t.Key != ""
	})).Return(nil)
	result, err := s.translationSvc.Merge(s.Ctx, translation)
	s.NoError(err)
	s.Len(result, 1)
	s.NotEmpty(result[0])
	s.True(strings.Contains(result[0], "test."))
}

func (s *translationTestSuite) Test_Merge_MultipleTrans() {
	trans1 := &localization.TranslationByKey{
		Key:          "test.key1",
		Category:     "test",
		Translations: map[string]string{"en": "Test value 1"},
	}

	trans2 := &localization.TranslationByKey{
		Key:          "test.key2",
		Category:     "test",
		Translations: map[string]string{"fr": "Test value 2"},
	}

	s.storage.On("Merge", s.Ctx, []*localization.TranslationByKey{trans1, trans2}).Return(nil)

	result, err := s.translationSvc.Merge(s.Ctx, trans1, trans2)

	s.NoError(err)
	s.Len(result, 2)
	s.Equal("test.key1", result[0])
	s.Equal("test.key2", result[1])
}

func (s *translationTestSuite) Test_Get_EmptyKey() {
	result, err := s.translationSvc.Get(s.Ctx, "", "en")
	s.Nil(result)
	s.Error(err)
	s.AssertAppErr(err, localization.ErrCodeLocTranslationKeyEmpty)
}

func (s *translationTestSuite) Test_Get_EmptyLang() {
	result, err := s.translationSvc.Get(s.Ctx, "test.key", "")
	s.Nil(result)
	s.Error(err)
	s.AssertAppErr(err, localization.ErrCodeLocTranslationLangEmpty)
}

func (s *translationTestSuite) Test_Get_UnsupportedLang() {
	result, err := s.translationSvc.Get(s.Ctx, "test.key", "es")
	s.Nil(result)
	s.Error(err)
	s.AssertAppErr(err, localization.ErrCodeNotSupportedLang)
}

func (s *translationTestSuite) Test_Get_NotFound() {
	s.storage.On("Get", s.Ctx, "test.key", "en").Return(nil, nil)
	result, err := s.translationSvc.Get(s.Ctx, "test.key", "en")
	s.Nil(result)
	s.NoError(err)
}

func (s *translationTestSuite) Test_Get_Success() {
	expected := &localization.Translation{
		Key:      "test.key",
		Value:    "Test value",
		Category: "test",
	}
	s.storage.On("Get", s.Ctx, "test.key", "en").Return(expected, nil)
	result, err := s.translationSvc.Get(s.Ctx, "test.key", "en")
	s.NoError(err)
	s.Equal(expected, result)
}

func (s *translationTestSuite) Test_MustGet_NotFound() {
	s.storage.On("Get", s.Ctx, "test.key", "en").Return(nil, nil)
	result, err := s.translationSvc.MustGet(s.Ctx, "test.key", "en")
	s.Nil(result)
	s.Error(err)
	s.AssertAppErr(err, localization.ErrCodeTranslationNotFound)
}

func (s *translationTestSuite) Test_MustGet_Success() {
	expected := &localization.Translation{
		Key:      "test.key",
		Value:    "Test value",
		Category: "test",
	}
	s.storage.On("Get", s.Ctx, "test.key", "en").Return(expected, nil)
	result, err := s.translationSvc.MustGet(s.Ctx, "test.key", "en")
	s.NoError(err)
	s.Equal(expected, result)
}

func (s *translationTestSuite) Test_MustGetByKeys_EmptyLang() {
	result, err := s.translationSvc.MustGetByKeys(s.Ctx, "", []string{"key1", "key2"})
	s.Nil(result)
	s.Error(err)
	s.AssertAppErr(err, localization.ErrCodeLocTranslationLangEmpty)
}

func (s *translationTestSuite) Test_MustGetByKeys_UnsupportedLang() {
	result, err := s.translationSvc.MustGetByKeys(s.Ctx, "es", []string{"key1", "key2"})
	s.Nil(result)
	s.Error(err)
	s.AssertAppErr(err, localization.ErrCodeNotSupportedLang)
}

func (s *translationTestSuite) Test_MustGetByKeys_EmptyKeys() {
	result, err := s.translationSvc.MustGetByKeys(s.Ctx, "en", []string{})
	s.NoError(err)
	s.NotNil(result)
	s.Empty(result)
}

func (s *translationTestSuite) Test_MustGetByKeys_NotAllFound() {
	keys := []string{"key1", "key2"}

	criteria := localization.GetTranslationsCriteria{
		Keys: []string{"key1", "key2"},
		Lang: "en",
	}

	response := &kit.PagingResponseG[localization.TranslationByKey]{
		Items: []*localization.TranslationByKey{
			{
				Key:          "key1",
				Category:     "test",
				Translations: map[string]string{"en": "Value 1"},
			},
		},
	}

	s.storage.On("Search", s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request:       criteria,
		PagingRequest: kit.PagingRequest{Skip: true},
	}).Return(response, nil)

	result, err := s.translationSvc.MustGetByKeys(s.Ctx, "en", keys)

	s.Nil(result)
	s.Error(err)
	s.AssertAppErr(err, localization.ErrCodeTranslationNotFound)
}

func (s *translationTestSuite) Test_MustGetByKeys_Success() {
	keys := []string{"key1", "key2"}

	criteria := localization.GetTranslationsCriteria{
		Keys: []string{"key1", "key2"},
		Lang: "en",
	}

	response := &kit.PagingResponseG[localization.TranslationByKey]{
		Items: []*localization.TranslationByKey{
			{
				Key:          "key1",
				Category:     "test",
				Translations: map[string]string{"en": "Value 1"},
			},
			{
				Key:          "key2",
				Category:     "test",
				Translations: map[string]string{"en": "Value 2"},
			},
		},
	}

	s.storage.On("Search", s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request:       criteria,
		PagingRequest: kit.PagingRequest{Skip: true},
	}).Return(response, nil)

	result, err := s.translationSvc.MustGetByKeys(s.Ctx, "en", keys)

	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 2)
	s.Equal("Value 1", result["key1"].Value)
	s.Equal("Value 2", result["key2"].Value)
}

func (s *translationTestSuite) Test_MustGetByKeys_DuplicateKeys() {
	keys := []string{"key1", "key1", "key2"}

	criteria := localization.GetTranslationsCriteria{
		Keys: []string{"key1", "key2"},
		Lang: "en",
	}

	response := &kit.PagingResponseG[localization.TranslationByKey]{
		Items: []*localization.TranslationByKey{
			{
				Key:          "key1",
				Category:     "test",
				Translations: map[string]string{"en": "Value 1"},
			},
			{
				Key:          "key2",
				Category:     "test",
				Translations: map[string]string{"en": "Value 2"},
			},
		},
	}

	s.storage.On("Search", s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request:       criteria,
		PagingRequest: kit.PagingRequest{Skip: true},
	}).Return(response, nil)

	result, err := s.translationSvc.MustGetByKeys(s.Ctx, "en", keys)

	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 2)
	s.Equal("Value 1", result["key1"].Value)
	s.Equal("Value 2", result["key2"].Value)
}

type mockTranslatable struct {
	keys []string
}

// Implement the Translatable interface
func (m *mockTranslatable) GetTranslatedKeys() []string {
	return m.keys
}

func (s *translationTestSuite) Test_TranslateObjects() {

	obj1 := &mockTranslatable{keys: []string{"key1", "key2"}}
	obj2 := &mockTranslatable{keys: []string{"key3"}}

	criteria := localization.GetTranslationsCriteria{
		Keys: []string{"key1", "key2", "key3"},
		Lang: "en",
	}

	response := &kit.PagingResponseG[localization.TranslationByKey]{
		Items: []*localization.TranslationByKey{
			{
				Key:          "key1",
				Category:     "test",
				Translations: map[string]string{"en": "Value 1"},
			},
			{
				Key:          "key2",
				Category:     "test",
				Translations: map[string]string{"en": "Value 2"},
			},
			{
				Key:          "key3",
				Category:     "test",
				Translations: map[string]string{"en": "Value 3"},
			},
		},
	}

	s.storage.On("Search", s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request:       criteria,
		PagingRequest: kit.PagingRequest{Skip: true},
	}).Return(response, nil)

	result, err := s.translationSvc.TranslateObjects(s.Ctx, "en", obj1, obj2)

	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 3)
	s.Equal("Value 1", result["key1"])
	s.Equal("Value 2", result["key2"])
	s.Equal("Value 3", result["key3"])
}

func (s *translationTestSuite) Test_MustTranslateObjects_Success() {

	obj1 := &mockTranslatable{keys: []string{"key1", "key2"}}

	criteria := localization.GetTranslationsCriteria{
		Keys: []string{"key1", "key2"},
		Lang: "en",
	}

	response := &kit.PagingResponseG[localization.TranslationByKey]{
		Items: []*localization.TranslationByKey{
			{
				Key:          "key1",
				Category:     "test",
				Translations: map[string]string{"en": "Value 1"},
			},
			{
				Key:          "key2",
				Category:     "test",
				Translations: map[string]string{"en": "Value 2"},
			},
		},
	}

	s.storage.On("Search", s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request:       criteria,
		PagingRequest: kit.PagingRequest{Skip: true},
	}).Return(response, nil)

	result := s.translationSvc.MustTranslateObjects(s.Ctx, "en", obj1)

	s.NotNil(result)
	s.Len(result, 2)
	s.Equal("Value 1", result["key1"])
	s.Equal("Value 2", result["key2"])
}

func (s *translationTestSuite) Test_MustTranslateObjects_Error() {

	obj1 := &mockTranslatable{keys: []string{"key1", "key2"}}

	criteria := localization.GetTranslationsCriteria{
		Keys: []string{"key1", "key2"},
		Lang: "en",
	}

	s.storage.On("Search", s.Ctx, &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request:       criteria,
		PagingRequest: kit.PagingRequest{Skip: true},
	}).Return(nil, localization.ErrTranslationNotFound(s.Ctx))

	result := s.translationSvc.MustTranslateObjects(s.Ctx, "en", obj1)

	s.NotNil(result)
	s.Len(result, 2)
	s.Equal("key1", result["key1"])
	s.Equal("key2", result["key2"])
}

func (s *translationTestSuite) Test_Search() {
	criteria := localization.GetTranslationsCriteria{
		Lang: "en",
	}

	request := &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request:       criteria,
		PagingRequest: kit.PagingRequest{Skip: true},
	}

	response := &kit.PagingResponseG[localization.TranslationByKey]{
		Items: []*localization.TranslationByKey{
			{
				Key:          "key1",
				Category:     "test",
				Translations: map[string]string{"en": "Value 1"},
			},
		},
	}

	s.storage.On("Search", s.Ctx, request).Return(response, nil)

	result, err := s.translationSvc.Search(s.Ctx, request)

	s.NoError(err)
	s.Equal(response, result)
}

func (s *translationTestSuite) Test_Search_UnsupportedLang() {
	criteria := localization.GetTranslationsCriteria{
		Lang: "es",
	}

	request := &kit.PagingRequestG[localization.GetTranslationsCriteria]{
		Request:       criteria,
		PagingRequest: kit.PagingRequest{Skip: true},
	}

	result, err := s.translationSvc.Search(s.Ctx, request)

	s.Nil(result)
	s.Error(err)
	s.AssertAppErr(err, localization.ErrCodeNotSupportedLang)
}

func (s *translationTestSuite) Test_GenKey() {
	key1 := s.translationSvc.GenKey("test.category")
	key2 := s.translationSvc.GenKey("test.category")

	s.NotEmpty(key1)
	s.NotEmpty(key2)
	s.NotEqual(key1, key2)
	s.True(strings.Contains(key1, "test.category."))
	s.True(strings.Contains(key2, "test.category."))
}

func (s *translationTestSuite) Test_BuildKey() {
	key := s.translationSvc.BuildKey("test.category", "123")
	s.Equal("test.category.123", key)
}
