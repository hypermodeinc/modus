/*
 * Copyright 2023 Hypermode, Inc.
 */

package schema

import (
	"fmt"

	"github.com/dgraph-io/gqlparser/ast"
	"github.com/dgraph-io/gqlparser/parser"
	"github.com/dgraph-io/gqlparser/validator"
	"github.com/rs/zerolog/log"
)

type FunctionInfo struct {
	PluginName string
	Schema     FunctionSchema
}

type FunctionSchema struct {
	ObjectDef *ast.Definition
	FieldDef  *ast.FieldDefinition
}

func (info FunctionInfo) FunctionName() string {
	return info.Schema.FunctionName()
}

func (schema FunctionSchema) Resolver() string {
	return schema.ObjectDef.Name + "." + schema.FieldDef.Name
}

func (schema FunctionSchema) FunctionName() string {
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

func (schema FunctionSchema) FunctionArgs() ast.ArgumentDefinitionList {
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
				for _, val := range v.([]any) {
					argName = val.(string)
					fld := schema.ObjectDef.Fields.ForName(argName)
					if fld == nil {
						log.Warn().
							Str("field", argName).
							Str("object", schema.ObjectDef.Name).
							Msg("Field does not exist.")
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

func GetFunctionSchema(scma string) ([]FunctionSchema, error) {

	// Parse the schema
	doc, parseErr := parser.ParseSchemas(validator.Prelude, &ast.Source{Input: scma})
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse GraphQL schema: %+v", parseErr)
	}

	// Find all fields with the @hm_function directive and add their schema info
	// to the map, using the resolver as a key.
	var results []FunctionSchema
	for _, def := range doc.Definitions {
		if def.Kind == ast.Object {
			for _, field := range def.Fields {
				if field.Directives.ForName("hm_function") != nil {
					schema := FunctionSchema{ObjectDef: def, FieldDef: field}
					results = append(results, schema)
				}
			}
		}
	}

	return results, nil
}
