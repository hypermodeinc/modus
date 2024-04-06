/*
 * Copyright 2024 Hypermode, Inc.
 */

package schemagen

import (
	"bytes"
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"hmruntime/plugins"
	"hmruntime/utils"

	"golang.org/x/exp/maps"
)

var firstPassComplete = false

func GetGraphQLSchema(metadata plugins.PluginMetadata, includeHeader bool) (string, error) {
	typeDefs := make(map[string]plugins.TypeDefinition, len(metadata.Types))
	errors := transformTypes(metadata.Types, &typeDefs)

	firstPassComplete = true

	functions, errs := transformFunctions(metadata.Functions, &typeDefs)
	errors = append(errors, errs...)

	if len(errors) > 0 {
		return "", fmt.Errorf("failed to generate schema: %+v", errors)
	}

	buf := bytes.Buffer{}
	if includeHeader {
		writeSchemaHeader(&buf, metadata)
	}
	writeSchema(&buf, functions, maps.Values(typeDefs))
	return buf.String(), nil
}

type TransformError struct {
	Source any
	Error  error
}

func transformTypes(types []plugins.TypeDefinition, typeDefs *map[string]plugins.TypeDefinition) []TransformError {
	errors := make([]TransformError, 0)
	for _, t := range types {
		if _, ok := (*typeDefs)[t.Name]; ok {
			errors = append(errors, TransformError{t, fmt.Errorf("type already exists: %s", t.Name)})
			continue
		}

		fields, err := convertNameTypePairs(t.Fields, typeDefs)
		if err != nil {
			errors = append(errors, TransformError{t, err})
			continue
		}

		(*typeDefs)[t.Name] = plugins.TypeDefinition{
			Name:   t.Name,
			Fields: fields,
		}
	}
	return errors
}

func transformFunctions(functions []plugins.FunctionSignature, typeDefs *map[string]plugins.TypeDefinition) ([]plugins.FunctionSignature, []TransformError) {
	results := make([]plugins.FunctionSignature, len(functions))
	errors := make([]TransformError, 0)
	for i, f := range functions {
		params, err := convertNameTypePairs(f.Parameters, typeDefs)
		if err != nil {
			errors = append(errors, TransformError{f, err})
			continue
		}

		returnType, err := convertType(f.ReturnType, typeDefs)
		if err != nil {
			errors = append(errors, TransformError{f, err})
			continue
		}

		results[i] = plugins.FunctionSignature{
			Name:       f.Name,
			Parameters: params,
			ReturnType: returnType,
		}
	}

	return results, errors
}

const timeFormat = "2006-01-02T15:04:05.999Z"

func writeSchemaHeader(buf *bytes.Buffer, metadata plugins.PluginMetadata) {
	buf.WriteString("# Hypermode Functions GraphQL Schema (auto-generated)\n")
	buf.WriteString("# \n")
	buf.WriteString("# Plugin: ")
	buf.WriteString(metadata.Plugin)
	buf.WriteByte('\n')
	buf.WriteString("# Library: ")
	buf.WriteString(metadata.Library)
	buf.WriteByte('\n')
	buf.WriteString("# Build ID: ")
	buf.WriteString(metadata.BuildId)
	buf.WriteByte('\n')
	buf.WriteString("# Build Time: ")
	buf.WriteString(metadata.BuildTime.Format(timeFormat))
	buf.WriteByte('\n')
	if metadata.GitRepo != "" {
		buf.WriteString("# Git Repo: ")
		buf.WriteString(metadata.GitRepo)
		buf.WriteByte('\n')
	}
	if metadata.GitCommit != "" {
		buf.WriteString("# Git Commit: ")
		buf.WriteString(metadata.GitCommit)
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
}

func writeSchema(buf *bytes.Buffer, functions []plugins.FunctionSignature, typeDefs []plugins.TypeDefinition) {

	// sort functions and type definitions
	slices.SortFunc(functions, func(a, b plugins.FunctionSignature) int {
		return cmp.Compare(a.Name, b.Name)
	})
	slices.SortFunc(typeDefs, func(a, b plugins.TypeDefinition) int {
		return cmp.Compare(a.Name, b.Name)
	})

	// write query functions
	buf.WriteString("type Query {\n")
	for _, f := range functions {
		buf.WriteString("  ")
		buf.WriteString(f.Name)
		if len(f.Parameters) > 0 {
			buf.WriteByte('(')
			for i, p := range f.Parameters {
				if i > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(p.Name)
				buf.WriteString(": ")
				buf.WriteString(p.Type)
			}
			buf.WriteByte(')')
		}
		buf.WriteString(": ")
		buf.WriteString(f.ReturnType)
		buf.WriteByte('\n')
	}
	buf.WriteByte('}')

	// write scalars
	wroteScalar := false
	for _, t := range typeDefs {
		if (len(t.Fields)) > 0 {
			continue
		}
		if !wroteScalar {
			wroteScalar = true
			buf.WriteByte('\n')
		}

		buf.WriteByte('\n')
		buf.WriteString("scalar ")
		buf.WriteString(t.Name)
	}

	// write types
	for _, t := range typeDefs {
		if (len(t.Fields)) == 0 {
			continue
		}

		buf.WriteString("\n\n")
		buf.WriteString("type ")
		buf.WriteString(t.Name)
		buf.WriteString(" {\n")
		for _, f := range t.Fields {
			buf.WriteString("  ")
			buf.WriteString(f.Name)
			buf.WriteString(": ")
			buf.WriteString(f.Type)
			buf.WriteByte('\n')
		}
		buf.WriteByte('}')
	}

	buf.WriteByte('\n')
}

func convertNameTypePairs(items []plugins.NameTypePair, typeDefs *map[string]plugins.TypeDefinition) ([]plugins.NameTypePair, error) {
	if len(items) == 0 {
		return nil, nil
	}

	results := make([]plugins.NameTypePair, len(items))
	for i, p := range items {
		t, err := convertType(p.Type, typeDefs)
		if err != nil {
			return nil, err
		}
		results[i] = plugins.NameTypePair{
			Name: p.Name,
			Type: t,
		}
	}
	return results, nil
}

var mapRegex = regexp.MustCompile(`^Map<(\w+<.+>|.+),\s*(\w+<.+>|.+)>$`)

func convertType(asType string, typeDefs *map[string]plugins.TypeDefinition) (string, error) {

	// Unwrap parentheses if present
	if strings.HasPrefix(asType, "(") && strings.HasSuffix(asType, ")") {
		return convertType(asType[1:len(asType)-1], typeDefs)
	}

	// Set the nullable flag.
	// In GraphQL, types are nullable by default,
	// and non-nullable types are indicated by a "!" suffix
	var n string
	if strings.HasSuffix(asType, " | null") {
		n = ""
		asType = asType[:len(asType)-7]
	} else {
		n = "!"
	}

	// check for array types
	if strings.HasSuffix(asType, "[]") {
		t, err := convertType(asType[:len(asType)-2], typeDefs)
		if err != nil {
			return "", err
		}
		return "[" + t + "]" + n, nil
	}

	// check for map types
	matches := mapRegex.FindStringSubmatch(asType)
	if len(matches) == 3 {
		kt, err := convertType(matches[1], typeDefs)
		if err != nil {
			return "", err
		}
		vt, err := convertType(matches[2], typeDefs)
		if err != nil {
			return "", err
		}

		// The pair type name will be composed from the key and value types.
		// ex: StringStringPair, IntStringPair, StringNullableStringPair, etc.
		ktn := utils.If(strings.HasSuffix(kt, "!"), kt[:len(kt)-1], "Nullable"+kt)
		vtn := utils.If(strings.HasSuffix(vt, "!"), vt[:len(vt)-1], "Nullable"+vt)
		typeName := ktn + vtn + "Pair"

		newType(typeName, []plugins.NameTypePair{{"key", kt}, {"value", vt}}, typeDefs)

		// The map is represented as a list of the pair type.
		// The list might be nullable, but the pair type within the list is always non-nullable.
		// ex: [StringStringPair!] or [StringStringPair!]!
		return "[" + typeName + "!]" + n, nil
	}

	// convert scalar types
	// TODO: How do we want to provide GraphQL ID scalar types? Maybe they're annotated? or maybe by naming convention?
	switch asType {
	case "string":
		return "String" + n, nil
	case "bool":
		return "Boolean" + n, nil
	case "i32", "i16", "i8", "u16", "u8":
		return "Int" + n, nil
	case "f64", "f32":
		return "Float" + n, nil
	case "i64":
		return newScalar("Int64", typeDefs) + n, nil
	case "u32":
		return newScalar("UInt", typeDefs) + n, nil
	case "u64":
		return newScalar("UInt64", typeDefs) + n, nil
	case "Date":
		return newScalar("DateTime", typeDefs) + n, nil
	}

	// in the first pass, we convert input custom type definitions
	if !firstPassComplete {
		return asType + n, nil
	}

	// going forward, convert custom types only if they have a type definition
	if _, ok := (*typeDefs)[asType]; ok {
		return asType + n, nil
	}

	return "", fmt.Errorf("unsupported type or missing type definition: %s", asType)
}

func newScalar(name string, typeDefs *map[string]plugins.TypeDefinition) string {
	return newType(name, nil, typeDefs)
}

func newType(name string, fields []plugins.NameTypePair, typeDefs *map[string]plugins.TypeDefinition) string {
	if _, ok := (*typeDefs)[name]; !ok {
		(*typeDefs)[name] = plugins.TypeDefinition{
			Name:   name,
			Fields: fields,
		}
	}
	return name
}
