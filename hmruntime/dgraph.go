package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/dgraph-io/gqlparser/ast"
	"github.com/dgraph-io/gqlparser/parser"
	"github.com/dgraph-io/gqlparser/validator"
)

func executeDQL(ctx context.Context, stmt string, isMutation bool) ([]byte, error) {
	reqBody := strings.NewReader(stmt)

	host := "http://localhost:8080" // TODO: make this configurable
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
		return nil, fmt.Errorf("error posting DQL statement: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DQL operation failed with status code %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading DQL response: %v", err)
	}

	return respBody, nil
}

func executeGQL(ctx context.Context, stmt string) ([]byte, error) {
	reqBody := strings.NewReader(stmt)
	resp, err := http.Post("http://localhost:8080/graphql", "application/graphql", reqBody)
	if err != nil {
		return nil, fmt.Errorf("error posting GraphQL statement: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GraphQL operation failed with status code %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading GraphQL response: %v", err)
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
		return "", fmt.Errorf("error getting GraphQL schema from Dgraph: %v", err)
	}

	var response dqlResponse[schemaResponse]
	err = json.Unmarshal(r, &response)
	if err != nil {
		return "", fmt.Errorf("error deserializing JSON of GraphQL schema: %v", err)
	}

	return response.Data.Node[0].Schema, nil
}

type functionSchemaInfo struct {
	ObjectDef *ast.Definition
	FieldDef  *ast.FieldDefinition
}

func (info functionSchemaInfo) Resolver() string {
	return info.ObjectDef.Name + "." + info.FieldDef.Name
}

func (info functionSchemaInfo) FunctionName() string {
	f := info.FieldDef

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

func (info functionSchemaInfo) FunctionArgs() ast.ArgumentDefinitionList {
	f := info.FieldDef

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
					fld := info.ObjectDef.Fields.ForName(argName)
					if fld == nil {
						log.Printf("Field %s.%s does not exist", info.ObjectDef.Name, argName)
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

func getFunctionSchemaInfos(schema string) []functionSchemaInfo {

	// Parse the schema
	doc, err := parser.ParseSchemas(validator.Prelude, &ast.Source{Input: schema})
	if err != nil {
		log.Fatal(err)
	}

	// Find all fields with the @hm_function directive
	var results []functionSchemaInfo
	for _, def := range doc.Definitions {
		if def.Kind == ast.Object {
			for _, field := range def.Fields {
				if field.Directives.ForName("hm_function") != nil {
					results = append(results, functionSchemaInfo{def, field})
				}
			}
		}
	}

	return results
}
