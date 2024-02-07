/*
 * Copyright 2023 Hypermode, Inc.
 */

package dgraph

import (
	"context"
	"fmt"
	"hmruntime/utils"
	"log"

	"github.com/dgraph-io/gqlparser/ast"
	"github.com/dgraph-io/gqlparser/parser"
	"github.com/dgraph-io/gqlparser/validator"
)

var DgraphUrl *string

type dgraphRequest struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

func ExecuteDQL[TResponse any](ctx context.Context, stmt string, vars map[string]string, isMutation bool) (TResponse, error) {
	var url string
	if isMutation {
		url = *DgraphUrl + "/mutate?commitNow=true"
	} else {
		url = *DgraphUrl + "/query"
	}

	request := dgraphRequest{
		Query:     stmt,
		Variables: vars,
	}

	response, err := utils.PostHttp[TResponse](url, request, nil)
	if err != nil {
		return response, fmt.Errorf("error posting DQL statement: %w", err)
	}

	return response, nil
}

func ExecuteGQL[TResponse any](ctx context.Context, stmt string, vars map[string]string) (TResponse, error) {
	url := *DgraphUrl + "/graphql"
	request := dgraphRequest{
		Query:     stmt,
		Variables: vars,
	}

	response, err := utils.PostHttp[TResponse](url, request, nil)
	if err != nil {
		return response, fmt.Errorf("error posting GraphQL statement: %w", err)
	}

	return response, nil
}

func GetGQLSchema(ctx context.Context) (string, error) {

	type DqlResponse[T any] struct {
		Data T `json:"data"`
	}

	type SchemaResponse struct {
		Node []struct {
			Schema string `json:"dgraph.graphql.schema"`
		} `json:"node"`
	}

	const query = "{node(func:has(dgraph.graphql.schema)){dgraph.graphql.schema}}"

	response, err := ExecuteDQL[DqlResponse[SchemaResponse]](ctx, query, nil, false)
	if err != nil {
		return "", fmt.Errorf("error getting GraphQL schema from Dgraph: %w", err)
	}

	data := response.Data
	if len(data.Node) == 0 {
		return "", fmt.Errorf("no GraphQL schema found in Dgraph")
	}

	return data.Node[0].Schema, nil
}

type FunctionInfo struct {
	PluginName string
	Schema     functionSchema
}

type functionSchema struct {
	ObjectDef *ast.Definition
	FieldDef  *ast.FieldDefinition
}

func (info FunctionInfo) FunctionName() string {
	return info.Schema.FunctionName()
}

func (schema functionSchema) Resolver() string {
	return schema.ObjectDef.Name + "." + schema.FieldDef.Name
}

func (schema functionSchema) FunctionName() string {
	f := schema.FieldDef

	// If @hm_function(name: "name") is specified, use that.
	d := f.Directives.ForName("hm_function")
	if d != nil {
		a := d.Arguments.ForName("name")
		if a != nil && a.Value != nil {
			return a.Value.Raw
		}
	}

	// No @hm_function directive, or no name argument. Just use the field name.
	return f.Name
}

func (schema functionSchema) FunctionArgs() ast.ArgumentDefinitionList {
	f := schema.FieldDef

	// If @hm_function(args: ["arg1", "arg2"]) is specified, use that.
	// The arguments must correspond to field names on the same parent object.
	// The types will be ascertained from the corresponding fields.
	// This is the case for fields on types other than Query and Mutation.
	d := f.Directives.ForName("hm_function")
	if d != nil {
		a := d.Arguments.ForName("args")
		if a != nil && a.Value != nil {
			v, err := a.Value.Value(nil)
			if err == nil {
				var list ast.ArgumentDefinitionList
				var argName string
				for _, val := range v.([]any) {
					argName = val.(string)
					fld := schema.ObjectDef.Fields.ForName(argName)
					if fld == nil {
						log.Printf("Field %s.%s does not exist", schema.ObjectDef.Name, argName)
						continue
					}
					arg := ast.ArgumentDefinition{
						Name: argName,
						Type: fld.Type,
					}
					list = append(list, &arg)
				}

				return list
			}
		}
	}

	// No @hm_function directive, or no args argument.
	// Just use the arguments on the field.
	// This is the case for Query and Mutation fields.
	return f.Arguments
}

func GetFunctionSchema(schema string) ([]functionSchema, error) {

	// Parse the schema
	doc, parseErr := parser.ParseSchemas(validator.Prelude, &ast.Source{Input: schema})
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse GraphQL schema: %+v", parseErr)
	}

	// Find all fields with the @hm_function directive and add their schema info
	// to the map, using the resolver as a key.
	var results []functionSchema
	for _, def := range doc.Definitions {
		if def.Kind == ast.Object {
			for _, field := range def.Fields {
				if field.Directives.ForName("hm_function") != nil {
					schema := functionSchema{def, field}
					results = append(results, schema)
				}
			}
		}
	}

	return results, nil
}

type ModelSpec struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Endpoint string `json:"endpoint"`
}
type ModelSpecInfo struct {
	Model ModelSpec `json:"model"`
}
type ModelSpecPayload struct {
	Data ModelSpecInfo `json:"data"`
}

func GetModelSpec(modelId string) (ModelSpec, error) {
	serviceURL := fmt.Sprintf("%s/admin", *DgraphUrl)

	query := `
		query GetModelSpec($id: String!) {
			model:getModelSpec(id: $id) {
				id
				type
				endpoint
			}
		}`

	request := map[string]any{
		"query":     query,
		"variables": map[string]any{"id": modelId},
	}

	response, err := utils.PostHttp[ModelSpecPayload](serviceURL, request, nil)
	if err != nil {
		return ModelSpec{}, fmt.Errorf("error getting model spec: %w", err)
	}

	spec := response.Data.Model
	if spec.ID != modelId {
		return ModelSpec{}, fmt.Errorf("error: ID does not match")
	}

	return spec, nil
}
