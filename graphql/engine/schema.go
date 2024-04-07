/*
 * Copyright 2024 Hypermode, Inc.
 */

package engine

import (
	"context"

	"hmruntime/graphql/datasource"

	"github.com/jensneuse/abstractlogger"
	"github.com/wundergraph/graphql-go-tools/execution/engine"
	gql "github.com/wundergraph/graphql-go-tools/execution/graphql"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/plan"
)

// ActivateSchema is a callback function that is called when a schema is loaded.
// It is responsible for preparing the schema and activating it on a new execution engine instance.
func ActivateSchema(ctx context.Context, schemaContent string) error {
	schema, err := gql.NewSchemaFromString(schemaContent)
	if err != nil {
		return err
	}

	queryTypeName := schema.QueryTypeName()
	args := schema.GetAllFieldArguments(gql.NewSkipReservedNamesFunc())
	queryFieldNames := make([]string, 0, len(args))
	for _, arg := range args {
		if arg.TypeName == queryTypeName {
			queryFieldNames = append(queryFieldNames, arg.FieldName)
		}
	}

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

	logger := abstractlogger.NoopLogger // TODO: Use a real logger

	e, err := engine.NewExecutionEngine(ctx, logger, engineConf)
	if err == nil {
		setEngine(e)
	}

	return err
}
