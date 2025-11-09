package v8

import (
	"context"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"gitlab.com/algmib/kit"
)

// Config - model of ES configuration
type Config struct {
	Url      string // Url - ES url
	Trace    bool   // Trace enables tracing mode
	Shards   int    // Shards - how many shards to be created for index
	Replicas int    // Replicas - how many replicas to eb created for index
	Username string // Username - ES basic auth (if not set, no auth applied)
	Password string // Password - ES basic auth
	Refresh  bool   // Refresh - enforces refresh after each change. It helpful for tests but MUST NOT BE USED ON PROD
}

// Index provides operations for managing Elasticsearch indices
type Index interface {
	// NewBuilder creates a new index builder instance
	NewBuilder() IndexBuilder
	// Create makes new index with specified mapping
	Create(ctx context.Context, index string, mapping MappingBody) error
	// Exists checks if index exists
	Exists(ctx context.Context, index string) (bool, error)
	// Delete removes indices
	Delete(ctx context.Context, indices []string) error
	// Refresh forces index refresh
	Refresh(ctx context.Context, index string) error
	// DeleteIndices removes indices by alias name
	DeleteIndices(ctx context.Context, alias string) error
	// GetIndices retrieves indices by alias name
	GetIndices(ctx context.Context, alias string) (*GetIndicesResponse, error)
	// GetMapping retrieves mapping for specified index
	GetMapping(ctx context.Context, index string) (*GetMappingResponse, error)
	// PutMapping updates mapping for existing index
	PutMapping(ctx context.Context, index string, mapping MappingBody) error
	// UpdateAliases performs alias operations like add/remove
	UpdateAliases(ctx context.Context, action *AliasAction) error
}

// Document provides operations for managing Elasticsearch documents
type Document interface {
	// Index indexes a document
	Index(ctx context.Context, index string, id string, data DataBody) error
	// Update updates an existing document in the specified index with the provided data.
	Update(ctx context.Context, index string, id string, data DataBody) error
	// Exists checks if a document exists in the index
	Exists(ctx context.Context, index, id string) (bool, error)
	// Delete deletes a document
	Delete(ctx context.Context, index string, id string) error
}

// Search defines an interface for executing search operations against a given index using a specified query.
type Search interface {
	// Search executes a search query on a specified index
	Search(ctx context.Context, index string, query QueryBody) (*SearchResponse, error)
}

// Es provides access to Elasticsearch functionality
type Es interface {
	Instance
	// Ping checks if Elasticsearch server is accessible
	Ping() bool
	// GetClient provides access to ES client
	GetClient() *elastic.Client
	// Close closes a client
	Close(ctx context.Context)
	// Instance returns the current instance
	Instance() Instance
}

type Instance interface {
	// Document returns interface for document operations
	Document() Document
	// Search returns interface for search operations
	Search() Search
	// Index returns interface for index operations
	Index() Index
}

// esImpl implements Es interface
type esImpl struct {
	*indexImpl
	*documentImpl
	*searchImpl
	*processorImpl

	logger kit.CLoggerFunc
	config *Config
	client *elastic.Client
}

// Document returns document implementation
func (s *esImpl) Document() Document {
	return s.documentImpl
}

// Search returns search implementation
func (s *esImpl) Search() Search {
	return s.searchImpl
}

// Index returns index implementation
func (s *esImpl) Index() Index {
	return s.indexImpl
}

// Instance returns the current instance
func (s *esImpl) Instance() Instance {
	return s
}

// l returns logger with a component name
func (s *esImpl) l() kit.CLogger {
	return s.logger().Cmp("es")
}

// NewEs creates a new Elasticsearch client instance
func NewEs(cfg *Config, logger kit.CLoggerFunc) (Es, error) {
	// Initialize implementation struct
	s := &esImpl{
		logger: logger,
		config: cfg,
	}

	// Configure Elasticsearch client
	esCfg := elastic.Config{
		Addresses: []string{cfg.Url},
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	// Create a new client
	cl, err := elastic.NewClient(esCfg)
	if err != nil {
		return nil, ErrNewClient(err)
	}
	s.client = cl

	// Initialize implementations
	s.indexImpl = &indexImpl{esImpl: s}
	s.documentImpl = &documentImpl{esImpl: s}
	s.searchImpl = &searchImpl{esImpl: s}
	s.processorImpl = &processorImpl{}

	s.l().Mth("new").F(kit.KV{"url": cfg.Url, "auth": cfg.Username != ""}).Inf("ok")
	return s, nil
}

// Ping checks if Elasticsearch server is accessible
func (s *esImpl) Ping() bool {
	s.l().Mth("ping").Dbg()
	res, err := s.client.Info()
	if err != nil {
		s.l().Mth("ping").Err(err)
		return false
	}
	defer res.Body.Close()

	return res.StatusCode == 200
}

// GetClient returns underlying Elasticsearch client
func (s *esImpl) GetClient() *elastic.Client {
	return s.client
}

// Close performs cleanup operations
func (s *esImpl) Close(ctx context.Context) {
	// nothing to do
}
