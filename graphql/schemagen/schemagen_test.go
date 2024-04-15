/*
 * Copyright 2024 Hypermode, Inc.
 */

package schemagen

import (
	"testing"

	"hmruntime/plugins"
	"hmruntime/utils"

	"github.com/stretchr/testify/require"
)

func Test_GetGraphQLSchema(t *testing.T) {

	metadata := plugins.PluginMetadata{
		Functions: []plugins.FunctionSignature{
			{
				Name: "add",
				Parameters: []plugins.Parameter{
					{Name: "a", Type: plugins.TypeInfo{Name: "i32"}},
					{Name: "b", Type: plugins.TypeInfo{Name: "i32"}},
				},
				ReturnType: plugins.TypeInfo{Name: "i32"},
			},
			{
				Name: "sayHello",
				Parameters: []plugins.Parameter{
					{Name: "name", Type: plugins.TypeInfo{Name: "string"}},
				},
				ReturnType: plugins.TypeInfo{Name: "string"},
			},
			{
				Name:       "currentTime",
				ReturnType: plugins.TypeInfo{Name: "Date"},
			},
			{
				Name: "transform",
				Parameters: []plugins.Parameter{
					{Name: "items", Type: plugins.TypeInfo{Name: "Map<string,string>"}},
				},
				ReturnType: plugins.TypeInfo{Name: "Map<string, string>"},
			},
		},
	}

	result, err := GetGraphQLSchema(metadata, false)

	expectedSchema := `type Query {
  add(a: Int!, b: Int!): Int!
  currentTime: DateTime!
  sayHello(name: String!): String!
  transform(items: [StringStringPair!]!): [StringStringPair!]!
}

scalar DateTime

type StringStringPair {
  key: String!
  value: String!
}
`
	require.Nil(t, err)
	require.Equal(t, expectedSchema, result)
}

func Test_ConvertType(t *testing.T) {
	testCases := []struct {
		inputType          string
		expectedOutputType string
		inputTypeDefs      []plugins.TypeDefinition
		expectedTypeDefs   []TypeDefinition
	}{
		// Plain non-nullable types
		{"string", "String!", nil, nil},
		{"bool", "Boolean!", nil, nil},
		{"i32", "Int!", nil, nil},
		{"i16", "Int!", nil, nil},
		{"i8", "Int!", nil, nil},
		{"u16", "Int!", nil, nil},
		{"u8", "Int!", nil, nil},
		{"f64", "Float!", nil, nil},
		{"f32", "Float!", nil, nil},

		// Array types
		{"string[]", "[String!]!", nil, nil},
		{"string[][]", "[[String!]!]!", nil, nil},
		{"(string | null)[]", "[String]!", nil, nil},

		// Custom scalar types
		{"Date", "DateTime!", nil, []TypeDefinition{{Name: "DateTime"}}},
		{"i64", "Int64!", nil, []TypeDefinition{{Name: "Int64"}}},
		{"u32", "UInt!", nil, []TypeDefinition{{Name: "UInt"}}},
		{"u64", "UInt64!", nil, []TypeDefinition{{Name: "UInt64"}}},

		// Custom types
		{"User", "User!",
			[]plugins.TypeDefinition{{
				Name: "User",
				Fields: []plugins.Field{
					{Name: "firstName", Type: plugins.TypeInfo{Name: "string"}},
					{Name: "lastName", Type: plugins.TypeInfo{Name: "string"}},
					{Name: "age", Type: plugins.TypeInfo{Name: "u8"}},
				},
			}},
			[]TypeDefinition{{
				Name: "User",
				Fields: []NameTypePair{
					{"firstName", "String!"},
					{"lastName", "String!"},
					{"age", "Int!"},
				},
			}}},

		// bool and numeric types can't be nullable in AssemblyScript
		// but string and custom types can
		{"string | null", "String", nil, nil},
		{"Foo | null", "Foo",
			[]plugins.TypeDefinition{{Name: "Foo"}},
			[]TypeDefinition{{Name: "Foo"}}},

		// Map types
		{"Map<string, string>", "[StringStringPair!]!", nil, []TypeDefinition{{
			Name: "StringStringPair",
			Fields: []NameTypePair{
				{"key", "String!"},
				{"value", "String!"},
			},
		}}},
		{"Map<string, string | null>", "[StringNullableStringPair!]!", nil, []TypeDefinition{{
			Name: "StringNullableStringPair",
			Fields: []NameTypePair{
				{"key", "String!"},
				{"value", "String"},
			},
		}}},
		{"Map<i32, string>", "[IntStringPair!]!", nil, []TypeDefinition{{
			Name: "IntStringPair",
			Fields: []NameTypePair{
				{"key", "Int!"},
				{"value", "String!"},
			},
		}}},
	}

	for _, tc := range testCases {
		t.Run(tc.inputType, func(t *testing.T) {

			typeDefs := make(map[string]TypeDefinition, len(tc.inputTypeDefs))
			errors := transformTypes(tc.inputTypeDefs, &typeDefs)
			require.Empty(t, errors)

			result, err := convertType(tc.inputType, &typeDefs)

			require.Nil(t, err)
			require.Equal(t, tc.expectedOutputType, result)

			if tc.expectedTypeDefs == nil {
				require.Empty(t, typeDefs)
			} else {
				require.Equal(t, tc.expectedTypeDefs, utils.MapValues(typeDefs))
			}
		})
	}
}
