/*
 * Copyright 2024 Hypermode, Inc.
 */

package engine

import (
	"fmt"
	"sync"

	"context"
	"strings"

	"hypruntime/config"
	"hypruntime/graphql/datasource"
	"hypruntime/graphql/schemagen"
	"hypruntime/logger"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
	"hypruntime/wasmhost"

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

func Activate(ctx context.Context, md *metadata.Metadata) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	schema, cfg, err := generateSchema(ctx, md)
	if err != nil {
		return err
	}

	datasourceConfig, err := getDatasourceConfig(ctx, schema, cfg)
	if err != nil {
		return err
	}

	engine, err := makeEngine(ctx, schema, datasourceConfig)
	if err != nil {
		return err
	}

	setEngine(engine)
	return nil
}

func generateSchema(ctx context.Context, md *metadata.Metadata) (*gql.Schema, *datasource.HypDSConfig, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	generated, err := schemagen.GetGraphQLSchema(ctx, md)
	if err != nil {
		return nil, nil, err
	}

	if utils.HypermodeDebugEnabled() {
		if config.UseJsonLogging {
			logger.Debug(ctx).Str("schema", generated.Schema).Msg("Generated schema")
		} else {
			fmt.Printf("\n%s\n", generated.Schema)
		}
	}

	schema, err := gql.NewSchemaFromString(generated.Schema)
	if err != nil {
		return nil, nil, err
	}

	cfg := &datasource.HypDSConfig{
		WasmHost: wasmhost.GetWasmHost(ctx),
		MapTypes: generated.MapTypes,
	}

	return schema, cfg, nil
}

func getDatasourceConfig(ctx context.Context, schema *gql.Schema, cfg *datasource.HypDSConfig) (plan.DataSourceConfiguration[datasource.HypDSConfig], error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	queryTypeName := schema.QueryTypeName()
	queryFieldNames := getAllQueryFields(ctx, schema)
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

	return plan.NewDataSourceConfiguration(
		datasource.DataSourceName,
		datasource.NewHypDSFactory(ctx),
		&plan.DataSourceMetadata{RootNodes: rootNodes, ChildNodes: childNodes},
		*cfg,
	)
}

func makeEngine(ctx context.Context, schema *gql.Schema, datasourceConfig plan.DataSourceConfiguration[datasource.HypDSConfig]) (*engine.ExecutionEngine, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	engineConfig := engine.NewConfiguration(schema)
	engineConfig.SetDataSources([]plan.DataSource{datasourceConfig})

	resolverOptions := resolve.ResolverOptions{
		MaxConcurrency:               1024,
		PropagateSubgraphErrors:      true,
		SubgraphErrorPropagationMode: resolve.SubgraphErrorPropagationModePassThrough,
	}

	adapter := newLoggerAdapter(ctx)
	return engine.NewExecutionEngine(ctx, adapter, engineConfig, resolverOptions)
}

func getAllQueryFields(ctx context.Context, s *gql.Schema) []string {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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
