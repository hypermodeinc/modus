package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dgraph-io/dgo/v230"
	"github.com/dgraph-io/dgo/v230/protos/api"
	"github.com/dgraph-io/gqlparser/ast"
	"github.com/dgraph-io/gqlparser/parser"
	"github.com/dgraph-io/gqlparser/validator"
)

func queryDQL(ctx context.Context, q string) ([]byte, error) {

	// TODO: This should use a persistent connection (or a pool of them)
	// TODO: The server endpoint should also be configurable

	// connect to dgraph server
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial("localhost:9080", creds)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Dgraph: %v", err)
	}

	// create a client and transaction
	client := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	txn := client.NewReadOnlyTxn()
	defer txn.Discard(ctx)

	// query dgraph
	response, err := txn.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("error querying Dgraph: %v", err)
	}

	// return the response
	return response.GetJson(), err
}

type schemaResponse struct {
	Node []struct {
		Schema string `json:"dgraph.graphql.schema"`
	} `json:"node"`
}

var schemaQuery = "{node(func:has(dgraph.graphql.schema)){dgraph.graphql.schema}}"

func getGQLSchema(ctx context.Context) (string, error) {

	r, err := queryDQL(ctx, schemaQuery)
	if err != nil {
		return "", fmt.Errorf("error getting GraphQL schema from Dgraph: %v", err)
	}

	var sr schemaResponse
	err = json.Unmarshal(r, &sr)
	if err != nil {
		return "", fmt.Errorf("error deserializing JSON of GraphQL schema: %v", err)
	}

	return sr.Node[0].Schema, nil
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
