/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package engine

import (
	"fmt"
	"os"
	"sync"

	"context"
	"strings"

	"github.com/fatih/color"
	"github.com/hypermodeinc/modus/lib/metadata"
	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/graphql/datasource"
	"github.com/hypermodeinc/modus/runtime/graphql/schemagen"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"

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

	if utils.DebugModeEnabled() {
		if app.Config().UseJsonLogging() {
			logger.Debug(ctx).Str("schema", generated.Schema).Msg("Generated schema")
		} else {
			fmt.Fprintf(os.Stderr, "\n%s\n", color.BlueString(generated.Schema))
		}
	}

	schema, err := gql.NewSchemaFromString(generated.Schema)
	if err != nil {
		return nil, nil, err
	}

	cfg := &datasource.HypDSConfig{
		WasmHost:          wasmhost.GetWasmHost(ctx),
		FieldsToFunctions: generated.FieldsToFunctions,
		MapTypes:          generated.MapTypes,
	}

	return schema, cfg, nil
}

func getDatasourceConfig(ctx context.Context, schema *gql.Schema, cfg *datasource.HypDSConfig) (plan.DataSourceConfiguration[datasource.HypDSConfig], error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	queryTypeName := schema.QueryTypeName()
	queryFieldNames := getTypeFields(ctx, schema, queryTypeName)

	mutationTypeName := schema.MutationTypeName()
	mutationFieldNames := getTypeFields(ctx, schema, mutationTypeName)

	rootNodes := []plan.TypeField{
		{
			TypeName:   queryTypeName,
			FieldNames: queryFieldNames,
		},
		{
			TypeName:   mutationTypeName,
			FieldNames: mutationFieldNames,
		},
	}

	childNodes := []plan.TypeField{}
	childNodes = append(childNodes, getChildNodes(queryFieldNames, schema, queryTypeName)...)
	childNodes = append(childNodes, getChildNodes(mutationFieldNames, schema, mutationTypeName)...)

	return plan.NewDataSourceConfiguration(
		datasource.DataSourceName,
		datasource.NewHypDSFactory(ctx),
		&plan.DataSourceMetadata{RootNodes: rootNodes, ChildNodes: childNodes},
		*cfg,
	)
}

func getChildNodes(fieldNames []string, schema *gql.Schema, typeName string) []plan.TypeField {
	var foundFields = make(map[string]bool)
	var childNodes []plan.TypeField
	for _, fieldName := range fieldNames {
		fields := schema.GetAllNestedFieldChildrenFromTypeField(typeName, fieldName, gql.NewSkipReservedNamesFunc())
		for _, field := range fields {
			if !foundFields[field.TypeName] {
				foundFields[field.TypeName] = true
				childNodes = append(childNodes, plan.TypeField{
					TypeName:   field.TypeName,
					FieldNames: field.FieldNames,
				})
			}
		}
	}
	return childNodes
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

func getTypeFields(ctx context.Context, s *gql.Schema, typeName string) []string {
	span, _ := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	doc := s.Document()
	fields := make([]string, 0)
	for _, objectType := range doc.ObjectTypeDefinitions {
		if doc.Input.ByteSliceString(objectType.Name) == typeName {
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
