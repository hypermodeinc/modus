/*
 * Copyright 2024 Hypermode, Inc.
 */

package schemagen

import (
	"bytes"
	"cmp"
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"

	"hypruntime/languages"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
)

type GraphQLSchema struct {
	Schema   string
	MapTypes []string
}

func GetGraphQLSchema(ctx context.Context, md *metadata.Metadata) (*GraphQLSchema, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	lti := languages.GetLanguageForSDK(md.SDK).TypeInfo()
	typeDefs, errors := transformTypes(md.Types, lti)
	functions, errs := transformFunctions(md.FnExports, typeDefs, lti)
	errors = append(errors, errs...)

	if len(errors) > 0 {
		return nil, fmt.Errorf("failed to generate schema: %+v", errors)
	}

	functions = filterFunctions(functions)
	types := filterTypes(utils.MapValues(typeDefs), functions)

	buf := bytes.Buffer{}
	writeSchema(&buf, functions, types)

	mapTypes := make([]string, 0, len(typeDefs))
	for _, t := range typeDefs {
		if t.IsMapType {
			mapTypes = append(mapTypes, t.Name)
		}
	}

	return &GraphQLSchema{
		Schema:   buf.String(),
		MapTypes: mapTypes,
	}, nil
}

type TransformError struct {
	Source any
	Error  error
}

func (e *TransformError) String() string {
	return fmt.Sprintf("source: %+v, error: %v", e.Source, e.Error)
}

func transformTypes(types metadata.TypeMap, lti languages.TypeInfo) (map[string]*TypeDefinition, []*TransformError) {
	typeDefs := make(map[string]*TypeDefinition, len(types))
	errors := make([]*TransformError, 0)
	for _, t := range types {
		name := lti.GetNameForType(t.Name)
		if _, ok := typeDefs[name]; ok {
			continue
		}

		fields, err := convertFields(t.Fields, lti, typeDefs)
		if err != nil {
			errors = append(errors, &TransformError{t, err})
			continue
		}

		typeDefs[name] = &TypeDefinition{
			Name:   name,
			Fields: fields,
		}
	}
	return typeDefs, errors
}

type FunctionSignature struct {
	Name       string
	Parameters []*ParameterSignature
	ReturnType string
}

type TypeDefinition struct {
	Name      string
	Fields    []*NameTypePair
	IsMapType bool
}

type NameTypePair struct {
	Name string
	Type string
}

type ParameterSignature struct {
	Name    string
	Type    string
	Default *any
}

func transformFunctions(functions metadata.FunctionMap, typeDefs map[string]*TypeDefinition, lti languages.TypeInfo) ([]*FunctionSignature, []*TransformError) {
	output := make([]*FunctionSignature, len(functions))
	errors := make([]*TransformError, 0)

	i := 0
	fnNames := utils.MapKeys(functions)
	sort.Strings(fnNames)
	for _, name := range fnNames {
		f := functions[name]

		params, err := convertParameters(f.Parameters, lti, typeDefs)
		if err != nil {
			errors = append(errors, &TransformError{f, err})
			continue
		}

		returnType, err := convertResults(f.Results, lti, typeDefs)
		if err != nil {
			errors = append(errors, &TransformError{f, err})
			continue
		}

		output[i] = &FunctionSignature{
			Name:       f.Name,
			Parameters: params,
			ReturnType: returnType,
		}

		i++
	}

	return output, errors
}

func filterFunctions(functions []*FunctionSignature) []*FunctionSignature {
	fnFilter := getFnFilter()
	results := make([]*FunctionSignature, 0, len(functions))
	for _, f := range functions {
		if fnFilter(f) {
			results = append(results, f)
		}
	}

	return results
}

func filterTypes(types []*TypeDefinition, functions []*FunctionSignature) []*TypeDefinition {
	// Filter out types that are not used by any function.
	// Also then recursively filter out types that are not used by any type.

	// Make a map of all types
	typeMap := make(map[string]*TypeDefinition, len(types))
	for _, t := range types {
		name := getBaseType(t.Name)
		typeMap[name] = t
	}

	// Get all types used by functions, including subtypes
	usedTypes := make(map[string]bool)
	for _, f := range functions {
		for _, p := range f.Parameters {
			addUsedTypes(p.Type, typeMap, usedTypes)
		}
		addUsedTypes(f.ReturnType, typeMap, usedTypes)
	}

	// Filter out types that are not used
	results := make([]*TypeDefinition, 0, len(types))
	for _, t := range types {
		name := getBaseType(t.Name)
		if usedTypes[name] {
			results = append(results, t)
		}
	}

	return results
}

func addUsedTypes(name string, types map[string]*TypeDefinition, usedTypes map[string]bool) {
	name = getBaseType(name)
	if usedTypes[name] {
		return
	}
	usedTypes[name] = true
	if t, ok := types[name]; ok {
		for _, f := range t.Fields {
			addUsedTypes(f.Type, types, usedTypes)
		}
	}
}

func getBaseType(name string) string {
	name = strings.TrimSuffix(name, "!")
	if strings.HasPrefix(name, "[") {
		return getBaseType(name[1 : len(name)-2])
	}

	return name
}

func writeSchema(buf *bytes.Buffer, functions []*FunctionSignature, typeDefs []*TypeDefinition) {

	// write header
	buf.WriteString("# Hypermode GraphQL Schema (auto-generated)\n\n")

	// sort functions and type definitions
	slices.SortFunc(functions, func(a, b *FunctionSignature) int {
		return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})
	slices.SortFunc(typeDefs, func(a, b *TypeDefinition) int {
		return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
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
				if p.Default != nil {
					val, err := utils.JsonSerialize(*p.Default)
					if err == nil {
						buf.WriteString(" = ")
						buf.Write(val)
					}
				}
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
		if len(t.Fields) > 0 || strings.HasSuffix(t.Name, "[]") || strings.HasPrefix(t.Name, "Map<") {
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

func convertParameters(parameters []*metadata.Parameter, lti languages.TypeInfo, typeDefs map[string]*TypeDefinition) ([]*ParameterSignature, error) {
	if len(parameters) == 0 {
		return nil, nil
	}

	output := make([]*ParameterSignature, len(parameters))
	for i, p := range parameters {

		t, err := convertType(p.Type, lti, typeDefs, false)
		if err != nil {
			return nil, err
		}

		// maintain compatibility with the deprecated "optional" field
		if p.Optional {
			output[i] = &ParameterSignature{
				Name: p.Name,
				Type: strings.TrimSuffix(t, "!"),
			}
			continue
		}

		output[i] = &ParameterSignature{
			Name:    p.Name,
			Type:    t,
			Default: p.Default,
		}
	}
	return output, nil
}

func convertResults(results []*metadata.Result, lti languages.TypeInfo, typeDefs map[string]*TypeDefinition) (string, error) {
	switch len(results) {
	case 0:
		return newScalar("Void", typeDefs), nil
	case 1:
		// Note: Single result doesn't use the name, even if it's present.
		return convertType(results[0].Type, lti, typeDefs, false)
	}

	fields := make([]*NameTypePair, len(results))
	for i, r := range results {
		name := r.Name
		if name == "" {
			name = fmt.Sprintf("item%d", i+1)
		}

		typ, err := convertType(r.Type, lti, typeDefs, false)
		if err != nil {
			return "", err
		}

		fields[i] = &NameTypePair{
			Name: name,
			Type: typ,
		}
	}

	t := getTypeForFields(fields, typeDefs)
	return t, nil
}

func getTypeForFields(fields []*NameTypePair, typeDefs map[string]*TypeDefinition) string {
	// see if an existing type already matches
	for _, t := range typeDefs {
		if len(t.Fields) != len(fields) {
			continue
		}

		found := true
		for i, f := range fields {
			if t.Fields[i].Name != f.Name || t.Fields[i].Type != f.Type {
				found = false
				break
			}
		}

		if found {
			return t.Name
		}
	}

	// there's no existing type that matches, so create a new one
	var name string
	for i := 1; ; i++ {
		name = fmt.Sprintf("_type%d", i)
		if _, ok := typeDefs[name]; !ok {
			break
		}
	}

	return newType(name, fields, typeDefs)
}

func convertFields(fields []*metadata.Field, lti languages.TypeInfo, typeDefs map[string]*TypeDefinition) ([]*NameTypePair, error) {
	if len(fields) == 0 {
		return nil, nil
	}

	results := make([]*NameTypePair, len(fields))
	for i, f := range fields {
		t, err := convertType(f.Type, lti, typeDefs, true)
		if err != nil {
			return nil, err
		}
		results[i] = &NameTypePair{
			Name: f.Name,
			Type: t,
		}
	}
	return results, nil
}

func convertType(typ string, lti languages.TypeInfo, typeDefs map[string]*TypeDefinition, firstPass bool) (string, error) {

	// Unwrap parentheses if present
	if strings.HasPrefix(typ, "(") && strings.HasSuffix(typ, ")") {
		return convertType(typ[1:len(typ)-1], lti, typeDefs, firstPass)
	}

	// Set the nullable flag.
	// In GraphQL, types are nullable by default,
	// and non-nullable types are indicated by a "!" suffix
	var n string
	if !lti.IsNullable(typ) {
		n = "!"
	}

	// unwrap nullable types (and dereference pointers)
	for lti.IsNullable(typ) {
		t := lti.GetUnderlyingType(typ)
		if t == typ {
			break
		}
		typ = t
	}

	// convert basic types
	// TODO: How do we want to provide GraphQL "ID" scalar types? Maybe they're annotated? or maybe by naming convention?

	if lti.IsStringType(typ) {
		return "String" + n, nil
	}

	if lti.IsByteSequenceType(typ) {
		// Note: If the bytes represent valid UTF-8 strings, Go will serialize them as actual strings.
		// Otherwise, the data will be base64 encoded.
		// TODO: We may want to ensure that the results are _always_ base64 encoded.
		return "String" + n, nil
	}

	if lti.IsBooleanType(typ) {
		return "Boolean" + n, nil
	}

	if lti.IsFloatType(typ) {
		return "Float" + n, nil
	}

	if lti.IsIntegerType(typ) {
		ctx := context.Background() // context is always unused for this purpose
		signed := lti.IsSignedIntegerType(typ)
		size, err := lti.GetSizeOfType(ctx, typ)
		if err != nil {
			return "", err
		}

		switch size {
		case 8:
			if signed {
				return newScalar("Int64", typeDefs) + n, nil
			} else {
				return newScalar("UInt64", typeDefs) + n, nil
			}
		case 4:
			if !signed {
				return newScalar("UInt", typeDefs) + n, nil
			}
		}

		return "Int" + n, nil
	}

	if lti.IsTimestampType(typ) {
		return newScalar("Timestamp", typeDefs) + n, nil
	}

	// check for array types
	if lti.IsListType(typ) {
		elem := lti.GetListSubtype(typ)
		t, err := convertType(elem, lti, typeDefs, firstPass)
		if err != nil {
			return "", err
		}
		return "[" + t + "]" + n, nil
	}

	// check for map types
	if lti.IsMapType(typ) {
		k, v := lti.GetMapSubtypes(typ)
		kt, err := convertType(k, lti, typeDefs, firstPass)
		if err != nil {
			return "", err
		}
		vt, err := convertType(v, lti, typeDefs, firstPass)
		if err != nil {
			return "", err
		}

		// The pair type name will be composed from the key and value types.
		// ex: StringStringPair, IntStringPair, StringNullableStringPair, etc.
		ktn := utils.If(strings.HasSuffix(kt, "!"), kt[:len(kt)-1], "Nullable"+kt)
		vtn := utils.If(strings.HasSuffix(vt, "!"), vt[:len(vt)-1], "Nullable"+vt)
		if ktn[0] == '[' {
			ktn = ktn[1:len(ktn)-2] + "List"
		}
		if vtn[0] == '[' {
			vtn = vtn[1:len(vtn)-2] + "List"
		}
		typeName := ktn + vtn + "Pair"

		newMapType(typeName, []*NameTypePair{{"key", kt}, {"value", vt}}, typeDefs)

		// The map is represented as a list of the pair type.
		// The list might be nullable, but the pair type within the list is always non-nullable.
		// ex: [StringStringPair!] or [StringStringPair!]!
		return "[" + typeName + "!]" + n, nil
	}

	name := lti.GetNameForType(typ)

	// in the first pass, we convert input custom type definitions
	if firstPass {
		return name + n, nil
	}

	// going forward, convert custom types only if they have a type definition
	if _, ok := typeDefs[name]; ok {
		return name + n, nil
	}

	return "", fmt.Errorf("unsupported type or missing type definition: %s", typ)
}

func newScalar(name string, typeDefs map[string]*TypeDefinition) string {
	return newType(name, nil, typeDefs)
}

func newType(name string, fields []*NameTypePair, typeDefs map[string]*TypeDefinition) string {
	if _, ok := typeDefs[name]; !ok {
		typeDefs[name] = &TypeDefinition{
			Name:   name,
			Fields: fields,
		}
	}
	return name
}

func newMapType(name string, fields []*NameTypePair, typeDefs map[string]*TypeDefinition) string {
	if _, ok := typeDefs[name]; !ok {
		typeDefs[name] = &TypeDefinition{
			Name:      name,
			Fields:    fields,
			IsMapType: true,
		}
	}
	return name
}
