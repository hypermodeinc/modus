/*
 * Copyright 2023 Hypermode, Inc.
 */
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgraph-io/gqlparser/ast"
	"github.com/dgraph-io/gqlparser/parser"
	"github.com/dgraph-io/gqlparser/validator"
)

func executeDQL(ctx context.Context, stmt string, isMutation bool) ([]byte, error) {
	reqBody := strings.NewReader(stmt)

	host := *dgraphUrl
	var endpoint, contentType string
	if isMutation {
		endpoint = "/mutate?commitNow=true"
		contentType = "application/rdf"
	} else {
		endpoint = "/query"
		contentType = "application/dql"
	}

	resp, err := http.Post(host+endpoint, contentType, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error posting DQL statement: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DQL operation failed with status code %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading DQL response: %w", err)
	}

	return respBody, nil
}

func executeGQL(ctx context.Context, stmt string) ([]byte, error) {
	reqBody := strings.NewReader(stmt)
	resp, err := http.Post(fmt.Sprintf("%s/graphql", *dgraphUrl), "application/graphql", reqBody)
	if err != nil {
		return nil, fmt.Errorf("error posting GraphQL statement: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GraphQL operation failed with status code %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading GraphQL response: %w", err)
	}

	return respBody, nil
}

type dqlResponse[T any] struct {
	Data T `json:"data"`
}

type schemaResponse struct {
	Node []struct {
		Schema string `json:"dgraph.graphql.schema"`
	} `json:"node"`
}

var schemaQuery = "{node(func:has(dgraph.graphql.schema)){dgraph.graphql.schema}}"

func getGQLSchema(ctx context.Context) (string, error) {

	r, err := executeDQL(ctx, schemaQuery, false)
	if err != nil {
		return "", fmt.Errorf("error getting GraphQL schema from Dgraph: %w", err)
	}

	var response dqlResponse[schemaResponse]
	err = json.Unmarshal(r, &response)
	if err != nil {
		return "", fmt.Errorf("error deserializing JSON of GraphQL schema: %w", err)
	}

	data := response.Data
	if len(data.Node) == 0 {
		return "", fmt.Errorf("no GraphQL schema found in Dgraph")
	}

	return data.Node[0].Schema, nil
}

type functionInfo struct {
	PluginName string
	Schema     functionSchema
}

type functionSchema struct {
	ObjectDef *ast.Definition
	FieldDef  *ast.FieldDefinition
}

func (info functionInfo) FunctionName() string {
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
	// The types will be assertained from the corresponding fields.
	// This is the case for fields on types other than Query and Mutation.
	d := f.Directives.ForName("hm_function")
	if d != nil {
		a := d.Arguments.ForName("args")
		if a != nil && a.Value != nil {
			v, err := a.Value.Value(nil)
			if err == nil {
				var list ast.ArgumentDefinitionList
				var argName string
				for _, val := range v.([]interface{}) {
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

func getFunctionSchema(schema string) ([]functionSchema, error) {

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

const (
	alphaService    string = "%v-alpha-service"
	classifierModel string = "classifier"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func getModelEndpoint(mid string) (string, error) {
	serviceURL := fmt.Sprintf("%s/admin", *dgraphUrl)

	query := `
		query GetModelSpec($id: String!) {
			model:getModelSpec(id: $id) {
				id
				type
				endpoint
			}
		}`

	payload := map[string]interface{}{
		"query":     query,
		"variables": map[string]string{"id": mid},
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling payload: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", serviceURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Create an instance of the ModelSpec struct
	var spec ModelSpecPayload

	// Unmarshal the JSON data into the ModelSpec struct
	err = json.Unmarshal(body, &spec)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response body: %w", err)
	}

	if spec.Data.Model.ID != mid {
		return "", fmt.Errorf("error: ID does not match")
	}

	if spec.Data.Model.Type != classifierModel {
		return "", fmt.Errorf("error: model type is not classifier")
	}

	return spec.Data.Model.Endpoint, nil
}
