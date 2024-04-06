/*
 * Copyright 2024 Hypermode, Inc.
 */

package schemagen

import (
	"testing"

	"hmruntime/plugins"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

func Test_GetGraphQLSchema(t *testing.T) {

	metadata := plugins.PluginMetadata{
		Functions: []plugins.FunctionSignature{
			{
				Name:       "add",
				Parameters: []plugins.NameTypePair{{"a", "i32"}, {"b", "i32"}},
				ReturnType: "i32",
			},
			{
				Name:       "sayHello",
				Parameters: []plugins.NameTypePair{{"name", "string"}},
				ReturnType: "string",
			},
			{
				Name:       "currentTime",
				ReturnType: "Date",
			},
			{
				Name:       "transform",
				Parameters: []plugins.NameTypePair{{"items", "Map<string, string>"}},
				ReturnType: "Map<string, string>",
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
		expectedTypeDefs   []plugins.TypeDefinition
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
		{"Date", "DateTime!", nil, []plugins.TypeDefinition{{Name: "DateTime"}}},
		{"i64", "Int64!", nil, []plugins.TypeDefinition{{Name: "Int64"}}},
		{"u32", "UInt!", nil, []plugins.TypeDefinition{{Name: "UInt"}}},
		{"u64", "UInt64!", nil, []plugins.TypeDefinition{{Name: "UInt64"}}},

		// Custom types
		{"User", "User!",
			[]plugins.TypeDefinition{{
				Name: "User",
				Fields: []plugins.NameTypePair{
					{"firstName", "string"},
					{"lastName", "string"},
					{"age", "u8"},
				},
			}},
			[]plugins.TypeDefinition{{
				Name: "User",
				Fields: []plugins.NameTypePair{
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
			[]plugins.TypeDefinition{{Name: "Foo"}}},

		// Map types
		{"Map<string, string>", "[StringStringPair!]!", nil, []plugins.TypeDefinition{{
			Name: "StringStringPair",
			Fields: []plugins.NameTypePair{
				{"key", "String!"},
				{"value", "String!"},
			},
		}}},
		{"Map<string, string | null>", "[StringNullableStringPair!]!", nil, []plugins.TypeDefinition{{
			Name: "StringNullableStringPair",
			Fields: []plugins.NameTypePair{
				{"key", "String!"},
				{"value", "String"},
			},
		}}},
		{"Map<i32, string>", "[IntStringPair!]!", nil, []plugins.TypeDefinition{{
			Name: "IntStringPair",
			Fields: []plugins.NameTypePair{
				{"key", "Int!"},
				{"value", "String!"},
			},
		}}},
	}

	for _, tc := range testCases {
		t.Run(tc.inputType, func(t *testing.T) {

			typeDefs := make(map[string]plugins.TypeDefinition, len(tc.inputTypeDefs))
			errors := transformTypes(tc.inputTypeDefs, &typeDefs)
			require.Empty(t, errors)

			result, err := convertType(tc.inputType, &typeDefs)

			require.Nil(t, err)
			require.Equal(t, tc.expectedOutputType, result)

			if tc.expectedTypeDefs == nil {
				require.Empty(t, typeDefs)
			} else {
				require.Equal(t, tc.expectedTypeDefs, maps.Values(typeDefs))
			}
		})
	}
}
