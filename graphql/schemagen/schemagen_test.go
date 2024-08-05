/*
 * Copyright 2024 Hypermode, Inc.
 */

package schemagen

import (
	"context"
	"testing"

	"hmruntime/plugins/metadata"
	"hmruntime/utils"

	"github.com/hypermodeAI/manifest"
	"github.com/stretchr/testify/require"
)

func Test_GetGraphQLSchema(t *testing.T) {

	manifest := &manifest.HypermodeManifest{
		Models: map[string]manifest.ModelInfo{},
		Hosts:  map[string]manifest.HostInfo{},
		Collections: map[string]manifest.CollectionInfo{
			"collection1": {
				SearchMethods: map[string]manifest.SearchMethodInfo{
					"search1": {
						Embedder: "myEmbedder",
					},
				},
			},
		},
	}

	makeDefault := func(val any) *any {
		return &val
	}

	md := metadata.Metadata{
		Functions: []metadata.Function{
			{
				Name: "add",
				Parameters: []metadata.Parameter{
					{Name: "a", Type: metadata.TypeInfo{Name: "i32"}},
					{Name: "b", Type: metadata.TypeInfo{Name: "i32"}},
				},
				ReturnType: metadata.TypeInfo{Name: "i32"},
			},
			{
				Name: "sayHello",
				Parameters: []metadata.Parameter{
					{Name: "name", Type: metadata.TypeInfo{Name: "string"}},
				},
				ReturnType: metadata.TypeInfo{Name: "string"},
			},
			{
				Name:       "currentTime",
				ReturnType: metadata.TypeInfo{Name: "Date"},
			},
			{
				Name: "transform",
				Parameters: []metadata.Parameter{
					{Name: "items", Type: metadata.TypeInfo{Name: "Map<string,string>"}},
				},
				ReturnType: metadata.TypeInfo{Name: "Map<string, string>"},
			},
			{
				Name: "testDefaultIntParams",
				Parameters: []metadata.Parameter{
					{Name: "a", Type: metadata.TypeInfo{Name: "i32"}, Default: nil},
					{Name: "b", Type: metadata.TypeInfo{Name: "i32"}, Default: makeDefault(0)},
					{Name: "c", Type: metadata.TypeInfo{Name: "i32"}, Default: makeDefault(1)},
				},
				ReturnType: metadata.TypeInfo{Name: "void"},
			},
			{
				Name: "testDefaultStringParams",
				Parameters: []metadata.Parameter{
					{Name: "a", Type: metadata.TypeInfo{Name: "string"}, Default: nil},
					{Name: "b", Type: metadata.TypeInfo{Name: "string"}, Default: makeDefault("")},
					{Name: "c", Type: metadata.TypeInfo{Name: "string"}, Default: makeDefault(`a"b`)},
					{Name: "d", Type: metadata.TypeInfo{Name: "string | null"}, Default: nil},
					{Name: "e", Type: metadata.TypeInfo{Name: "string | null"}, Default: makeDefault(nil)},
					{Name: "f", Type: metadata.TypeInfo{Name: "string | null"}, Default: makeDefault("")},
					{Name: "g", Type: metadata.TypeInfo{Name: "string | null"}, Default: makeDefault("test")},
				},
				ReturnType: metadata.TypeInfo{Name: "void"},
			},
			{
				Name: "testDefaultArrayParams",
				Parameters: []metadata.Parameter{
					{Name: "a", Type: metadata.TypeInfo{Name: "i32[]"}, Default: nil},
					{Name: "b", Type: metadata.TypeInfo{Name: "i32[]"}, Default: makeDefault([]int32{})},
					{Name: "c", Type: metadata.TypeInfo{Name: "i32[]"}, Default: makeDefault([]int32{1, 2, 3})},
					{Name: "d", Type: metadata.TypeInfo{Name: "i32[] | null"}, Default: nil},
					{Name: "e", Type: metadata.TypeInfo{Name: "i32[] | null"}, Default: makeDefault(nil)},
					{Name: "f", Type: metadata.TypeInfo{Name: "i32[] | null"}, Default: makeDefault([]int32{})},
					{Name: "g", Type: metadata.TypeInfo{Name: "i32[] | null"}, Default: makeDefault([]int32{1, 2, 3})},
				},
				ReturnType: metadata.TypeInfo{Name: "void"},
			},
			{
				Name:       "getPerson",
				ReturnType: metadata.TypeInfo{Name: "Person"},
			},
			{
				Name:       "getPeople",
				ReturnType: metadata.TypeInfo{Name: "Person[]"},
			},
			{
				Name:       "getProductMap",
				ReturnType: metadata.TypeInfo{Name: "Map<string, Product>"},
			},
			{
				Name:       "doNothing",
				ReturnType: metadata.TypeInfo{Name: "void"},
			},
			// This should be excluded from the final schema
			{
				Name: "myEmbedder",
				Parameters: []metadata.Parameter{
					{Name: "text", Type: metadata.TypeInfo{Name: "string"}},
				},
				ReturnType: metadata.TypeInfo{Name: "f64[]"},
			},
		},
		Types: []metadata.TypeDefinition{
			{
				Name: "Company",
				Fields: []metadata.Field{
					{Name: "name", Type: metadata.TypeInfo{Name: "string"}},
				},
			},
			{
				Name: "Product",
				Fields: []metadata.Field{
					{Name: "name", Type: metadata.TypeInfo{Name: "string"}},
					{Name: "price", Type: metadata.TypeInfo{Name: "f64"}},
					{Name: "manufacturer", Type: metadata.TypeInfo{Name: "Company"}},
					{Name: "components", Type: metadata.TypeInfo{Name: "Product[]"}},
				},
			},
			{
				Name: "Person",
				Fields: []metadata.Field{
					{Name: "name", Type: metadata.TypeInfo{Name: "string"}},
					{Name: "age", Type: metadata.TypeInfo{Name: "i32"}},
					{Name: "addresses", Type: metadata.TypeInfo{Name: "Address[]"}},
				},
			},
			{
				Name: "Person[]",
			},
			{
				Name: "Address",
				Fields: []metadata.Field{
					{Name: "street", Type: metadata.TypeInfo{Name: "string"}},
					{Name: "city", Type: metadata.TypeInfo{Name: "string"}},
					{Name: "state", Type: metadata.TypeInfo{Name: "string"}},
					{Name: "country", Type: metadata.TypeInfo{Name: "string"}},
					{Name: "postalCode", Type: metadata.TypeInfo{Name: "string"}},
					{Name: "location", Type: metadata.TypeInfo{Name: "Coordinates"}},
				},
			},
			{
				Name: "Coordinates",
				Fields: []metadata.Field{
					{Name: "lat", Type: metadata.TypeInfo{Name: "f64"}},
					{Name: "lon", Type: metadata.TypeInfo{Name: "f64"}},
				},
			},
			// This should be excluded from the final schema
			{
				Name: "Header",
				Fields: []metadata.Field{
					{Name: "name", Type: metadata.TypeInfo{Name: "string"}},
					{Name: "values", Type: metadata.TypeInfo{Name: "string[]"}},
				},
			},
		},
	}

	result, err := GetGraphQLSchema(context.Background(), md, manifest, false)

	expectedSchema := `type Query {
  add(a: Int!, b: Int!): Int!
  currentTime: Timestamp!
  doNothing: Void
  getPeople: [Person!]!
  getPerson: Person!
  getProductMap: [StringProductPair!]!
  sayHello(name: String!): String!
  testDefaultArrayParams(a: [Int!]!, b: [Int!]! = [], c: [Int!]! = [1,2,3], d: [Int!], e: [Int!] = null, f: [Int!] = [], g: [Int!] = [1,2,3]): Void
  testDefaultIntParams(a: Int!, b: Int! = 0, c: Int! = 1): Void
  testDefaultStringParams(a: String!, b: String! = "", c: String! = "a\"b", d: String, e: String = null, f: String = "", g: String = "test"): Void
  transform(items: [StringStringPair!]!): [StringStringPair!]!
}

scalar Timestamp
scalar Void

type Address {
  street: String!
  city: String!
  state: String!
  country: String!
  postalCode: String!
  location: Coordinates!
}

type Company {
  name: String!
}

type Coordinates {
  lat: Float!
  lon: Float!
}

type Person {
  name: String!
  age: Int!
  addresses: [Address!]!
}

type Product {
  name: String!
  price: Float!
  manufacturer: Company!
  components: [Product!]!
}

type StringProductPair {
  key: String!
  value: Product!
}

type StringStringPair {
  key: String!
  value: String!
}
`
	require.Nil(t, err)
	require.Equal(t, expectedSchema, result.Schema)
}

func Test_ConvertType(t *testing.T) {
	testCases := []struct {
		inputType          string
		expectedOutputType string
		inputTypeDefs      []metadata.TypeDefinition
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
		{"Date", "Timestamp!", nil, []TypeDefinition{{Name: "Timestamp"}}},
		{"i64", "Int64!", nil, []TypeDefinition{{Name: "Int64"}}},
		{"u32", "UInt!", nil, []TypeDefinition{{Name: "UInt"}}},
		{"u64", "UInt64!", nil, []TypeDefinition{{Name: "UInt64"}}},

		// Custom types
		{"User", "User!",
			[]metadata.TypeDefinition{{
				Name: "User",
				Fields: []metadata.Field{
					{Name: "firstName", Type: metadata.TypeInfo{Name: "string"}},
					{Name: "lastName", Type: metadata.TypeInfo{Name: "string"}},
					{Name: "age", Type: metadata.TypeInfo{Name: "u8"}},
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
			[]metadata.TypeDefinition{{Name: "Foo"}},
			[]TypeDefinition{{Name: "Foo"}}},

		// Map types
		{"Map<string, string>", "[StringStringPair!]!", nil, []TypeDefinition{{
			Name: "StringStringPair",
			Fields: []NameTypePair{
				{"key", "String!"},
				{"value", "String!"},
			},
			IsMapType: true,
		}}},
		{"Map<string, string | null>", "[StringNullableStringPair!]!", nil, []TypeDefinition{{
			Name: "StringNullableStringPair",
			Fields: []NameTypePair{
				{"key", "String!"},
				{"value", "String"},
			},
			IsMapType: true,
		}}},
		{"Map<i32, string>", "[IntStringPair!]!", nil, []TypeDefinition{{
			Name: "IntStringPair",
			Fields: []NameTypePair{
				{"key", "Int!"},
				{"value", "String!"},
			},
			IsMapType: true,
		}}},
		{"Map<string, Map<string, f32>>", "[StringStringFloatPairListPair!]!", nil, []TypeDefinition{
			{
				Name: "StringStringFloatPairListPair",
				Fields: []NameTypePair{
					{"key", "String!"},
					{"value", "[StringFloatPair!]!"},
				},
				IsMapType: true,
			},
			{
				Name: "StringFloatPair",
				Fields: []NameTypePair{
					{"key", "String!"},
					{"value", "Float!"},
				},
				IsMapType: true,
			},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.inputType, func(t *testing.T) {

			typeDefs := make(map[string]TypeDefinition, len(tc.inputTypeDefs))
			errors := transformTypes(tc.inputTypeDefs, typeDefs)
			require.Empty(t, errors)

			result, err := convertType(tc.inputType, typeDefs, false)

			require.Nil(t, err)
			require.Equal(t, tc.expectedOutputType, result)

			if tc.expectedTypeDefs == nil {
				require.Empty(t, typeDefs)
			} else {
				require.ElementsMatch(t, tc.expectedTypeDefs, utils.MapValues(typeDefs))
			}
		})
	}
}
