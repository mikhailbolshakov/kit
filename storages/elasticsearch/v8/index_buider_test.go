//go:build example

package v8

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/stretchr/testify/suite"
	"testing"
)

type esIndexTestSuite struct {
	kit.Suite
}

func (s *esIndexTestSuite) SetupSuite() {
	s.Suite.Init(func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) })
}

func (s *esIndexTestSuite) SetupTest() {
}

func (s *esIndexTestSuite) TearDownSuite() {}

func TestEsV8IndexSuite(t *testing.T) {
	suite.Run(t, new(esIndexTestSuite))
}

func (s *esIndexTestSuite) Test_Index_ChangingModelMappingWithCfgSettings() {
	es, err := NewEs(&Config{
		Url:      url,
		Trace:    true,
		Shards:   2,
		Replicas: 2,
	}, s.L)
	s.NoError(err)

	index := kit.NewRandString()
	type model struct {
		Field string `json:"field" es:"type:keyword"`
	}
	err = es.Index().NewBuilder().WithIndex(index).WithMappingModel(&model{}).Build(s.Ctx)
	s.NoError(err)

	type modelNew struct {
		Field  string `json:"field" es:"type:keyword"`
		Field2 string `json:"field2" es:"type:keyword"`
	}
	err = es.Index().NewBuilder().WithIndex(index).WithMappingModel(&modelNew{}).Build(s.Ctx)
	s.NoError(err)
}

func (s *esIndexTestSuite) Test_IndexNested_SubProperties() {
	es, err := NewEs(&Config{
		Url:   url,
		Trace: false,
	}, s.L)
	s.NoError(err)

	type nested struct {
		Name  string  `json:"n" es:"type:keyword"`
		Value float64 `json:"v" es:"type:double"`
	}

	type model struct {
		Field  string   `json:"field" es:"type:keyword"`
		Array  []nested `json:"array" es:"type:nested"`
		Scaled float64  `json:"scaled" es:"type:scaled_float;scaling_factor:1000"`
	}

	index := kit.NewRandString()
	err = es.Index().NewBuilder().WithIndex(index).WithMappingModel(&model{}).Build(s.Ctx)
	s.NoError(err)
}

func (s *esIndexTestSuite) Test_Index_ChangingModelMapping() {
	es, err := NewEs(&Config{
		Url:   url,
		Trace: false,
	}, s.L)
	s.NoError(err)

	type model struct {
		Field string `json:"field" es:"type:keyword"`
	}
	index := kit.NewRandString()
	err = es.Index().NewBuilder().WithIndex(index).WithMappingModel(&model{}).Build(s.Ctx)
	s.NoError(err)

	type modelNew struct {
		Field  string `json:"field" es:"type:keyword"`
		Field2 string `json:"field2" es:"type:keyword"`
	}

	err = es.Index().NewBuilder().WithIndex(index).WithMappingModel(&modelNew{}).Build(s.Ctx)
	s.NoError(err)
}

func (s *esIndexTestSuite) Test_Index_ChangingMapping_ExistentFields() {
	es, err := NewEs(&Config{
		Url:   url,
		Trace: false,
	}, s.L)
	s.NoError(err)

	type model struct {
		Field  string  `json:"field" es:"type:keyword"`
		Scaled float64 `json:"scaled" es:"type:scaled_float;scaling_factor:1000"`
	}

	index := kit.NewRandString()
	err = es.Index().NewBuilder().WithIndex(index).WithMappingModel(&model{}).Build(s.Ctx)
	s.NoError(err)
	type modelNew struct {
		Field string `json:"field" es:"type:text"`
	}
	err = es.Index().NewBuilder().WithIndex(index).WithMappingModel(&modelNew{}).Build(s.Ctx)
	s.NotNil(err)
}

func (s *esIndexTestSuite) Test_Index_ChangingMapping_ExplicitMapping() {
	es, err := NewEs(&Config{
		Url:      url,
		Trace:    true,
		Shards:   2,
		Replicas: 2,
	}, s.L)
	s.NoError(err)

	mapping := `
{
	"mappings": {
		"properties": {
			"field1": {
				"type": "keyword"
			}
		}
	}
}
`
	index := kit.NewRandString()
	err = es.Index().NewBuilder().WithIndex(index).WithExplicitMapping(mapping).Build(s.Ctx)
	s.NoError(err)
	newMapping := `
{
	"mappings": {
		"properties": {
			"field1": {
				"type": "keyword"
			},
			"field2": {
				"type": "keyword"
			}
		}
	}
}
`
	err = es.Index().NewBuilder().WithIndex(index).WithExplicitMapping(newMapping).Build(s.Ctx)
	s.NoError(err)
}

func (s *esIndexTestSuite) Test_Index_Mapping_WhenNotIndexField() {
	es, err := NewEs(&Config{
		Url:   url,
		Trace: false,
	}, s.L)
	s.NoError(err)

	type model struct {
		Field1 string `json:"field1" es:"type:keyword"`
		Field2 string `json:"field2" es:"type:keyword;-"`
	}

	index := kit.NewRandString()
	err = es.Index().NewBuilder().WithIndex(index).WithMappingModel(&model{}).Build(s.Ctx)
	s.NoError(err)
}

func (s *esIndexTestSuite) Test_Alias_ChangingModelMapping() {
	es, err := NewEs(&Config{
		Url:      url,
		Trace:    true,
		Shards:   2,
		Replicas: 2,
	}, s.L)
	s.NoError(err)

	// create alias and index
	alias := kit.NewRandString()
	type model struct {
		Field string `json:"field" es:"type:keyword"`
	}
	err = es.Index().NewBuilder().WithAlias(alias).WithMappingModel(&model{}).Build(s.Ctx)
	s.NoError(err)

	// index document
	err = es.Document().Index(s.Ctx, alias, kit.NewId(), JsonData(&model{Field: "value"}))
	s.NoError(err)

	// search
	err = es.Index().Refresh(s.Ctx, alias)
	s.NoError(err)

	// get data from alias
	srchRs, err := es.Search().Search(s.Ctx, alias, MatchAllQuery())
	s.NoError(err)
	s.Equal(int64(1), srchRs.Hits.Total.Value)

	// change mapping
	type modelNew struct {
		Field  string `json:"field" es:"type:keyword"`
		Field2 string `json:"field2" es:"type:keyword"`
	}
	err = es.Index().NewBuilder().WithAlias(alias).WithMappingModel(&modelNew{}).Build(s.Ctx)
	s.NoError(err)

	// index document
	err = es.Document().Index(s.Ctx, alias, kit.NewId(), JsonData(&modelNew{Field: "value", Field2: "value2"}))
	s.NoError(err)

	// search
	err = es.Index().Refresh(s.Ctx, alias)
	s.NoError(err)

	// get data from alias
	srchRs, err = es.Search().Search(s.Ctx, alias, MatchAllQuery())
	s.NoError(err)
	s.Equal(int64(2), srchRs.Hits.Total.Value)
}
