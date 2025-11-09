package v8

const (
	EsTypeKeyword         = "keyword"
	EsTypeText            = "text"
	EsTypeDate            = "date"
	EsTypeBool            = "boolean"
	EsTypeLong            = "long"
	EsTypeInteger         = "integer"
	EsTypeSearchAsYouType = "search_as_you_type"
	EsTypeFlattened       = "flattened"
	EsTypeNested          = "nested"
	EsTypeFloat           = "float"
	EsTypeDouble          = "double"
	EsTypeScaledFloat     = "scaled_float"
)

const (
	EsSortRequestMissingFirst = "_first"
	EsSortRequestMissingLast  = "_last"
)

var typesMap = map[string]struct{}{
	EsTypeKeyword:         {},
	EsTypeText:            {},
	EsTypeDate:            {},
	EsTypeBool:            {},
	EsTypeSearchAsYouType: {},
	EsTypeLong:            {},
	EsTypeInteger:         {},
	EsTypeFlattened:       {},
	EsTypeNested:          {},
	EsTypeFloat:           {},
	EsTypeDouble:          {},
	EsTypeScaledFloat:     {},
}

// Index block

type GetIndicesResponse struct {
	Indices map[string]EsIndex
}

type EsIndex struct {
	Aliases map[string]EsAlias `json:"aliases"`
}

type EsAlias struct {
	IsWriteIndex bool `json:"is_write_index"`
}

type EsProperties map[string]*EsProperty

type EsProperty struct {
	Type          string        `json:"type,omitempty"`           // Type specifies a datatype
	Index         *bool         `json:"index,omitempty"`          // Index - if false, field isn't indexed
	ScalingFactor *float64      `json:"scaling_factor,omitempty"` // ScalingFactor defines the multiplier for scaling numeric types in Elasticsearch.
	Properties    *EsProperties `json:"properties,omitempty"`     // Properties defines nested properties structure for an Elasticsearch property, used for handling complex or object types.
}

type GetMappingResponse struct {
	Mappings map[string]EsMapping `json:"mappings"`
}

type EsMapping struct {
	Settings struct {
		NumberOfShards   int `json:"number_of_shards"`
		NumberOfReplicas int `json:"number_of_replicas"`
	} `json:"settings"`
	Mappings struct {
		Properties EsProperties `json:"properties"`
	} `json:"mappings"`
}

type AliasAction struct {
	Actions []AliasActions `json:"actions"`
}

type AliasActions struct {
	Add *AliasAddAction `json:"add,omitempty"`
}

type AliasAddAction struct {
	Index        string `json:"index"`
	Alias        string `json:"alias"`
	IsWriteIndex bool   `json:"is_write_index,omitempty"`
}

// Search block

type SearchResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
	} `json:"hits"`
}
