/*
 * Copyright 2024 Hypermode, Inc.
 */

package engine

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"context"
	"strings"

	"hmruntime/config"
	"hmruntime/graphql/datasource"
	"hmruntime/graphql/schemagen"
	"hmruntime/logger"
	"hmruntime/plugins"

	"github.com/wundergraph/graphql-go-tools/execution/engine"
	gql "github.com/wundergraph/graphql-go-tools/execution/graphql"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/plan"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

var instance *engine.ExecutionEngine
var mutex sync.RWMutex

// GetEngine provides thread-safe access to the current GraphQL execution engine.
func GetEngine() *engine.ExecutionEngine {
	mutex.RLock()
	defer mutex.RUnlock()
	return instance
}

func setEngine(engine *engine.ExecutionEngine) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = engine
}

func Activate(ctx context.Context, metadata plugins.PluginMetadata) error {

	schemaContent, err := schemagen.GetGraphQLSchema(metadata, true)
	if err != nil {
		return err
	}

	if b, err := strconv.ParseBool(os.Getenv("HYPERMODE_DEBUG")); err == nil && b {
		if config.UseJsonLogging {
			logger.Debug(ctx).Str("schema", schemaContent).Msg("Generated schema")
		} else {
			fmt.Printf("\n%s\n", schemaContent)
		}
	}

	schema, err := gql.NewSchemaFromString(schemaContent)
	if err != nil {
		return err
	}

	queryTypeName := schema.QueryTypeName()
	queryFieldNames := getAllQueryFields(schema)
	rootNodes := []plan.TypeField{
		{
			TypeName:   queryTypeName,
			FieldNames: queryFieldNames,
		},
	}

	var childNodes []plan.TypeField
	for _, f := range queryFieldNames {
		fields := schema.GetAllNestedFieldChildrenFromTypeField(queryTypeName, f, gql.NewSkipReservedNamesFunc())
		for _, field := range fields {
			childNodes = append(childNodes, plan.TypeField{
				TypeName:   field.TypeName,
				FieldNames: field.FieldNames,
			})
		}
	}

	dsCfg, err := plan.NewDataSourceConfiguration(
		datasource.DataSourceName,
		&datasource.Factory[datasource.Configuration]{Ctx: ctx},
		&plan.DataSourceMetadata{RootNodes: rootNodes, ChildNodes: childNodes},
		datasource.Configuration{},
	)
	if err != nil {
		return err
	}

	engineConf := engine.NewConfiguration(schema)
	engineConf.SetDataSources([]plan.DataSource{
		dsCfg,
	})

	resolverOptions := resolve.ResolverOptions{
		MaxConcurrency:          1024,
		PropagateSubgraphErrors: true,
	}

	adapter := NewLoggerAdapter(ctx)
	e, err := newExecutionEngine(ctx, adapter, engineConf, resolverOptions)
	if err == nil {
		setEngine(e)
	}

	return err
}

func getAllQueryFields(s *gql.Schema) []string {
	doc := s.Document()
	queryTypeName := s.QueryTypeName()

	fields := make([]string, 0)
	for _, objectType := range doc.ObjectTypeDefinitions {
		typeName := doc.Input.ByteSliceString(objectType.Name)
		if typeName == queryTypeName {
			for _, fieldRef := range objectType.FieldsDefinition.Refs {
				field := doc.FieldDefinitions[fieldRef]
				fieldName := doc.Input.ByteSliceString(field.Name)
				if !strings.HasPrefix(fieldName, "__") {
					fields = append(fields, fieldName)
				}
			}
			break
		}
	}

	return fields
}
