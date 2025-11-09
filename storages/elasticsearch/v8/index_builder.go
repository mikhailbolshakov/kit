package v8

import (
	"context"
	"fmt"
	"github.com/mikhailbolshakov/kit"
	"reflect"
	"strconv"
	"strings"
)

// IndexBuilder allows creating / modification a ES index
type IndexBuilder interface {
	// WithAlias specifies an index with alias
	WithAlias(name string) IndexBuilder
	// WithIndex specifies an index name
	// if alias specified with WithAlias call, you don't need specify index name explicitly
	WithIndex(name string) IndexBuilder
	// WithMappingModel specifies index mapping based on model provided
	// if index doesn't exist, a new index is created
	// if index exists it checks whether existent mapping is modified and if it is, it fails. If only new fields added, it handles them as PUT
	// Note, "json" tag must be specified together with "es" tag
	// example:
	// type IndexModel struct {
	//   Field1 string `json:"field1" es:"type:text"`   	// field is mapped with text type
	//   Field2 string `json:"field2" es:"type:keyword"` 	// field is mapped with keyword type
	//   Field3 time.Time `json:"field3" es:"type:date"` 	// field is mapped with date type
	//   Field4 time.Time `json:"field4" es:"-"` 			// field is mapped with "Index=false"
	// }
	//
	// model must be pointer type
	WithMappingModel(model interface{}) IndexBuilder
	// WithExplicitMapping specifies index mapping explicitly as a serialized json mapping object
	// see ES doc about https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping.html
	// it checks if index exists, if not creates it
	WithExplicitMapping(mapping string) IndexBuilder
	// Build builds a new alias/index or modifies mapping of an existent index
	Build(ctx context.Context) error
}

type indexBuilder struct {
	Index
	alias           string
	index           string
	mappingModel    interface{}
	mappingExplicit string
	cfg             *Config
	logger          kit.CLoggerFunc
}

func newIndexBuilder(index Index, cfg *Config, logger kit.CLoggerFunc) IndexBuilder {
	return &indexBuilder{
		Index:  index,
		cfg:    cfg,
		logger: logger,
	}
}

func (e *indexBuilder) l() kit.CLogger {
	return e.logger().Cmp("es-v8-idx-builder")
}

func (e *indexBuilder) WithAlias(name string) IndexBuilder {
	e.alias = name
	return e
}

func (e *indexBuilder) WithIndex(name string) IndexBuilder {
	e.index = name
	return e
}

func (e *indexBuilder) WithMappingModel(model interface{}) IndexBuilder {
	e.mappingModel = model
	return e
}

func (e *indexBuilder) WithExplicitMapping(mapping string) IndexBuilder {
	e.mappingExplicit = mapping
	return e
}

func (e *indexBuilder) getMapping(ctx context.Context, index string) (*EsMapping, error) {

	// get current mapping
	curMappings, err := e.GetMapping(ctx, index)
	if err != nil {
		return nil, err
	}
	curMapping, ok := curMappings.Mappings[index]
	if !ok {
		return nil, ErrIndexBuilderNoMappingFound(ctx, index)
	}

	mappingJson, _ := kit.JsonEncode(curMapping)
	currentMapping := &EsMapping{}

	currentMapping, err = kit.JsonDecode[EsMapping](mappingJson)
	if err != nil {
		return nil, ErrIndexBuilderMappingSchemaNotExpected(ctx, err, index)
	}

	return currentMapping, nil
}

// modelToMapping creates ES mapping based on model tag
// check model_mapping_test for usage details
func (e *indexBuilder) modelToMapping(ctx context.Context, modelObj interface{}) (*EsMapping, error) {
	e.l().C(ctx).Mth("model-to-mapping").Dbg()

	if modelObj == nil {
		return nil, nil
	}

	type params map[string]string

	if reflect.ValueOf(modelObj).Kind() != reflect.Ptr || reflect.TypeOf(modelObj).Elem().Kind() != reflect.Struct {
		return nil, ErrIndexBuilderInvalidModelType(ctx)
	}

	// takes type description
	r := reflect.TypeOf(modelObj).Elem()
	mappingProperties := make(EsProperties)

	// build mapping fields map
	// go through fields
	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)

		// check json tag
		// we use index field name from json mapping
		// if json tag missing, field is skipped
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {

			jsonParams := strings.Split(jsonTag, ",")
			// if there is no field name in json tag
			if len(jsonParams) == 0 {
				return nil, ErrIndexBuilderInvalidModel(ctx)
			}
			indexFieldName := jsonParams[0]

			// take es tag
			esTag := field.Tag.Get("es")
			// if es tag missing, skip the field
			if esTag != "" {
				esTagParams := make(params)

				// take params separated by ;
				params := strings.Split(esTag, ";")
				for _, p := range params {
					kv := strings.Split(p, ":")
					if len(kv) == 2 {
						esTagParams[kv[0]] = kv[1]
					} else {
						esTagParams[kv[0]] = ""
					}
				}

				// populate mapping params
				mappingProperties[indexFieldName] = &EsProperty{}

				if _, ok := esTagParams["-"]; ok {
					// if empty sign exists
					f := false
					mappingProperties[indexFieldName].Type = EsTypeText
					mappingProperties[indexFieldName].Index = &f
				} else {
					if v, ok := esTagParams["type"]; ok {
						if _, ok := typesMap[v]; !ok {
							return nil, ErrIndexBuilderInvalidModel(ctx)
						} else {
							if v == EsTypeNested {
								// slice of struct (The nested type is a specialised version of the object data type that allows arrays of objects to be indexed in a way that they can be queried independently of each other.)
								if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
									subStruct := reflect.New(field.Type.Elem()).Interface()
									subProps, err := e.modelToMapping(ctx, subStruct)
									if err != nil {
										return nil, err
									}
									if subProps != nil {
										mappingProperties[indexFieldName].Properties = &subProps.Mappings.Properties
									}
								}
							}
							mappingProperties[indexFieldName].Type = v
						}
					}
					if v, ok := esTagParams["scaling_factor"]; ok {
						sf, err := strconv.ParseFloat(v, 64)
						if err != nil {
							return nil, ErrIndexBuilderInvalidModel(ctx)
						}
						mappingProperties[indexFieldName].ScalingFactor = &sf
					}
				}

			}
		}
	}

	// return ES mapping if specified
	if len(mappingProperties) == 0 {
		return nil, ErrIndexBuilderInvalidModel(ctx)
	} else {
		r := &EsMapping{}
		r.Mappings.Properties = mappingProperties
		return r, nil
	}

}

func (e *indexBuilder) getNewMapping(ctx context.Context, index string) (*EsMapping, error) {
	// get new mapping
	newMapping := &EsMapping{}
	var err error
	if e.mappingModel != nil {
		// if mapping specified as a model
		newMapping, err = e.modelToMapping(ctx, e.mappingModel)
		if err != nil {
			return nil, err
		}
	} else {
		// if mapping specified explicitly
		err = kit.Unmarshal([]byte(e.mappingExplicit), newMapping)
		if err != nil {
			return nil, ErrIndexBuilderMappingSchemaNotExpected(ctx, err, index)
		}
	}
	return newMapping, nil
}

func (e *indexBuilder) modifyMapping(ctx context.Context, index string, curMapping, newMapping *EsMapping) error {
	l := e.l().C(ctx).Mth("modify-mapping").Dbg()

	// check if there are changes in existent fields
	if v := e.checkExistentFieldsMappingModified(curMapping, newMapping); len(v) > 0 {
		return ErrIndexBuilderMappingExistentFieldsModified(ctx, index, v)
	}

	// extract added fields
	if addedFieldsMapping := e.addedFieldsMapping(curMapping, newMapping); len(addedFieldsMapping.Mappings.Properties) > 0 {
		if err := e.PutMapping(ctx, index, JsonMapping(addedFieldsMapping.Mappings)); err != nil {
			return err
		}
		l.DbgF("fields added: %+v", addedFieldsMapping.Mappings.Properties)
	}
	return nil
}

func (e *indexBuilder) createIndex(ctx context.Context, index string, mapping *EsMapping) error {
	l := e.l().C(ctx).Mth("create-index").F(kit.KV{"index": index}).Dbg()
	// create
	if err := e.Create(ctx, index, JsonMapping(e.setSettings(mapping))); err != nil {
		return err
	}
	l.Dbg("created")
	return nil
}

func (e *indexBuilder) buildAlias(ctx context.Context, alias string) error {
	l := e.l().C(ctx).Mth("build-alias").F(kit.KV{"alias": alias}).Dbg()

	// check alias exists
	exists, err := e.Exists(ctx, alias)
	if err != nil {
		return err
	}

	if exists {
		// we allow adding new fields to mapping
		// but don't allow changing existent ones
		l.DbgF("alias %s exists", alias)

		// get indexes by alias
		aliasesRs, err := e.GetIndices(ctx, alias)
		if err != nil {
			return err
		}
		if len(aliasesRs.Indices) == 0 {
			return ErrIndexAliasesNoAliasFound(ctx, alias)
		}
		// get writable index
		var writeIndexName string
	loop:
		for idxName, idx := range aliasesRs.Indices {
			for _, ia := range idx.Aliases {
				if ia.IsWriteIndex {
					writeIndexName = idxName
					break loop
				}
			}
		}
		if writeIndexName == "" {
			return ErrIndexBuilderNoWriteIndexForAlias(ctx, alias)
		}
		l.F(kit.KV{"writeIndex": writeIndexName})

		// get current mapping
		currentMapping, err := e.getMapping(ctx, writeIndexName)
		if err != nil {
			return err
		}

		// get new mapping
		newMapping, err := e.getNewMapping(ctx, writeIndexName)
		if err != nil {
			return err
		}

		// modify mapping for alias (it modifies mapping for all the indexes)
		if err = e.modifyMapping(ctx, alias, currentMapping, newMapping); err != nil {
			return err
		}

		l.Dbg("modified")
	} else {
		// new alias

		// get new mapping
		newMapping, err := e.getNewMapping(ctx, alias)
		if err != nil {
			return err
		}

		// create write index
		idxName := fmt.Sprintf("%s-idx-%s", alias, kit.Now().Format("20060102150405"))
		if err = e.createIndex(ctx, idxName, newMapping); err != nil {
			return err
		}

		if err = e.UpdateAliases(ctx, &AliasAction{
			Actions: []AliasActions{
				{
					Add: &AliasAddAction{
						Index:        idxName,
						Alias:        alias,
						IsWriteIndex: true,
					},
				},
			},
		}); err != nil {
			return err
		}

		l.Dbg("created")
	}
	return nil
}

func (e *indexBuilder) buildIndex(ctx context.Context, index string) error {
	l := e.l().C(ctx).Mth("build-index").F(kit.KV{"index": index}).Dbg()

	// check index exists
	exists, err := e.Exists(ctx, index)
	if err != nil {
		return err
	}

	if exists {

		// get current mapping
		currentMapping, err := e.getMapping(ctx, index)
		if err != nil {
			return err
		}

		// get new mapping
		newMapping, err := e.getNewMapping(ctx, index)
		if err != nil {
			return err
		}

		// modify mapping for index
		err = e.modifyMapping(ctx, index, currentMapping, newMapping)
		if err != nil {
			return err
		}

		l.Dbg("modified")
	} else {

		// new index
		// get new mapping
		newMapping, err := e.getNewMapping(ctx, index)
		if err != nil {
			return err
		}

		// create write index
		err = e.createIndex(ctx, index, newMapping)
		if err != nil {
			return err
		}

		l.Dbg("created")
	}
	return nil
}

func (e *indexBuilder) Build(ctx context.Context) error {
	e.l().Mth("build").Dbg()

	// check passed params
	if e.alias == "" && e.index == "" {
		return ErrIndexBuilderAliasAndIndexEmpty(ctx)
	}
	if e.mappingExplicit == "" && e.mappingModel == nil {
		return ErrIndexBuilderModelEmpty(ctx)
	}

	// alias-based
	if e.alias != "" {
		return e.buildAlias(ctx, e.alias)
	} else {
		return e.buildIndex(ctx, e.index)
	}

}

func (e *indexBuilder) setSettings(mapping *EsMapping) *EsMapping {
	if mapping.Settings.NumberOfReplicas == 0 {
		mapping.Settings.NumberOfReplicas = e.cfg.Replicas
		if mapping.Settings.NumberOfReplicas == 0 {
			mapping.Settings.NumberOfReplicas = 1
		}
	}
	if mapping.Settings.NumberOfShards == 0 {
		mapping.Settings.NumberOfShards = e.cfg.Shards
		if mapping.Settings.NumberOfShards == 0 {
			mapping.Settings.NumberOfShards = 1
		}
	}
	return mapping
}

func (e *indexBuilder) addedFieldsMapping(currentMapping, newMapping *EsMapping) *EsMapping {
	addedFieldsMapping := &EsMapping{}
	addedFieldsMapping.Mappings.Properties = make(EsProperties)
	for f, v := range newMapping.Mappings.Properties {
		if _, found := currentMapping.Mappings.Properties[f]; !found {
			addedFieldsMapping.Mappings.Properties[f] = v
		}
	}
	return addedFieldsMapping
}

// checkExistentFieldsMappingModified compares current and provided mapping and returns true if there are changes in existent fields
func (e *indexBuilder) checkExistentFieldsMappingModified(currentMapping, newMapping *EsMapping) []string {
	var modifiedFields []string
	for curFieldName, curField := range currentMapping.Mappings.Properties {
		for newFieldName, newField := range newMapping.Mappings.Properties {
			if curFieldName == newFieldName && curField.Type != newField.Type {
				modifiedFields = append(modifiedFields, curFieldName)
			}
		}
	}
	return modifiedFields
}
