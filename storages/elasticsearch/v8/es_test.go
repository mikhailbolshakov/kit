//go:build example

package v8

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/stretchr/testify/suite"
	"testing"
)

type test struct {
	User     string `json:"user"`
	Message  string `json:"message"`
	Retweets int64  `json:"retweets"`
}

type esTestSuite struct {
	kit.Suite
}

func (s *esTestSuite) SetupSuite() {
	s.Suite.Init(func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) })
}

func (s *esTestSuite) SetupTest() {
}

func (s *esTestSuite) TearDownSuite() {}

func TestEsV8Suite(t *testing.T) {
	suite.Run(t, new(esTestSuite))
}

const (
	url = "http://localhost:9200"
)

func (s *esTestSuite) Test_Simple() {
	client, err := NewEs(&Config{
		Url: url,
	}, s.L)

	s.NoError(err)
	s.True(client.Ping())

	indexName := kit.NewRandString()

	isExist, err := client.Index().Exists(s.Ctx, indexName)
	s.NoError(err)
	s.False(isExist)

	// Create a new index.
	mappings := `{
	"mappings":{
		"properties":{
			"user":{
				"type":"keyword"
			},
			"message":{
				"type":"text"
			},
			"retweets":{
				"type":"long"
			}
		}}}`
	s.NoError(client.Index().Create(s.Ctx, indexName, StringMapping(mappings)))
	defer client.Index().Delete(s.Ctx, []string{indexName})

	// Index a tweet (using JSON serialization)
	tweet1 := test{User: "olivere", Message: "Take Five", Retweets: 0}
	s.NoError(client.Document().Index(s.Ctx, indexName, "1", JsonData(tweet1)))

	// Refresh to make sure the documents are searchable.
	s.NoError(client.Index().Refresh(s.Ctx, indexName))

	// Get tweet with specified ID
	get1, err := client.GetClient().Get(indexName, "1")
	s.NoError(err)
	s.False(get1.IsError())

	// Search with a term query
	//termQuery := elastic.NewTermQuery("user", "olivere")
	//searchResult, err := client.Search().
	//	Search("twitter").        // documentImpl in index "twitter"
	//	QueryBody(termQuery).        // specify the query
	//	Sort("user", true).      // sort by "user" field, ascending
	//	From(0).Size(10).        // take documents 0-9
	//	Pretty(true).            // pretty print request and response JSON
	//	Do(context.Background()) // execute
	//if err != nil {
	//	// Handle error
	//	t.Fatal(err)
	//}

}

func (s *esTestSuite) Test_UpdatePartially() {
	client, err := NewEs(&Config{
		Url: url,
	}, s.L)
	s.NoError(err)

	index := "partially-update-" + kit.NewRandString()

	mappings := `{"mappings":{
		"properties":{
			"ID":{
				"type":"keyword"
			},
			"name":{
				"type":"text"
			},
			"price":{
				"type":"scaled_float",
				"scaling_factor":10000
			}}}}`

	product := []struct {
		Id    string  `json:"ID"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}{
		{Id: kit.NewId(), Name: kit.NewRandString(), Price: 14.5794},
		{Id: kit.NewId(), Name: kit.NewRandString(), Price: 1543},
		{Id: kit.NewId(), Name: kit.NewRandString(), Price: 500.5},
	}

	s.NoError(client.Index().Create(s.Ctx, index, StringMapping(mappings)))

	s.NoError(client.Document().Index(s.Ctx, index, product[0].Id, JsonData(product[0])))
	s.NoError(client.Document().Index(s.Ctx, index, product[1].Id, JsonData(product[1])))
	s.NoError(client.Document().Index(s.Ctx, index, product[2].Id, JsonData(product[2])))

	s.NoError(client.Index().Refresh(s.Ctx, index))

	search, err := client.Search().Search(s.Ctx, index, MatchAllQuery())
	s.NoError(err)
	s.Equal(int64(3), search.Hits.Total.Value)

	// updates all partially - only price
	updates := `{
		"doc": {
			"price": 499.99
	}}`

	// no docs found by identity (404)
	s.Error(client.Document().Update(s.Ctx, index, product[0].Name, StringData(updates)))

	// update 2 products
	s.NoError(client.Document().Update(s.Ctx, index, product[0].Id, StringData(updates)))
	s.NoError(client.Document().Update(s.Ctx, index, product[1].Id, StringData(updates)))

	s.NoError(client.Index().Refresh(s.Ctx, index))

	search, err = client.Search().Search(s.Ctx, index, MatchAllQuery())
	s.NoError(err)
	s.Equal(int64(3), search.Hits.Total.Value)
}

func (s *esTestSuite) Test_Alias() {
	client, err := NewEs(&Config{
		Url: url,
	}, s.L)
	s.NoError(err)
	// create a new index
	aliasName := kit.NewRandString()
	indexNameWritable := aliasName + "-idx-1"
	indexNameNonWritable := aliasName + "-idx-2"
	mappings := `{
	"mappings":{
		"properties":{
			"user":{
				"type":"keyword"
			},
			"message":{
				"type":"text"
			}
		}}}`
	msgData := struct {
		User    string `json:"user"`
		Message string `json:"message"`
	}{
		User:    kit.NewId(),
		Message: "some text",
	}

	s.NoError(client.Index().Create(s.Ctx, indexNameWritable, StringMapping(mappings)))

	s.NoError(client.Index().Create(s.Ctx, indexNameNonWritable, StringMapping(mappings)))

	// add indexes to alias
	s.NoError(client.Index().UpdateAliases(s.Ctx, &AliasAction{
		Actions: []AliasActions{
			{Add: &AliasAddAction{Index: indexNameWritable, Alias: aliasName, IsWriteIndex: true}},
			{Add: &AliasAddAction{Index: indexNameNonWritable, Alias: aliasName, IsWriteIndex: false}},
		},
	}))

	aliases, err := client.Index().GetIndices(s.Ctx, aliasName)
	s.NoError(err)
	s.Len(aliases.Indices, 2)

	// check alias exists
	exists, err := client.Index().Exists(s.Ctx, aliasName)
	s.NoError(err)
	s.True(exists)

	// check index exists
	exists, err = client.Index().Exists(s.Ctx, indexNameWritable)
	s.NoError(err)
	s.True(exists)

	// check index exists
	exists, err = client.Index().Exists(s.Ctx, indexNameNonWritable)
	s.NoError(err)
	s.True(exists)

	// can write through alias to writable index
	s.NoError(client.Document().Index(s.Ctx, aliasName, kit.NewId(), JsonData(msgData)))
	s.NoError(client.Document().Index(s.Ctx, indexNameNonWritable, kit.NewId(), JsonData(msgData)))

	s.NoError(client.Index().Refresh(s.Ctx, indexNameWritable))
	s.NoError(client.Index().Refresh(s.Ctx, indexNameNonWritable))

	// get data from alias
	srchRs, err := client.Search().Search(s.Ctx, aliasName, MatchAllQuery())
	s.NoError(err)
	s.Equal(int64(2), srchRs.Hits.Total.Value)

	// change mapping through alias
	modifiedMappings := `{
		"properties":{
			"user":{
				"type":"keyword"
			},
			"message":{
				"type":"text"
			},
			"field":{
				"type":"text"
			}}}`
	s.NoError(client.Index().PutMapping(s.Ctx, aliasName, StringMapping(modifiedMappings)))

	// get mapping
	curMappings, err := client.Index().GetMapping(s.Ctx, aliasName)
	s.NoError(err)
	s.NotEmpty(curMappings)
}
