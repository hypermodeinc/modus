/*
 * Copyright 2024 Hypermode, Inc.
 */

package schemagen

import (
	"context"
	"testing"

	"hmruntime/languages"
	"hmruntime/manifestdata"
	"hmruntime/plugins/metadata"
	"hmruntime/utils"

	"github.com/hypermodeAI/manifest"
	"github.com/stretchr/testify/require"
)

func Test_GetGraphQLSchema_AssemblyScript(t *testing.T) {

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
	manifestdata.SetManifest(manifest)

	md := metadata.NewPluginMetadata()
	md.SDK = "functions-as" // AssemblyScript

	md.FnExports.AddFunction("add").
		WithParameter("a", "i32").
		WithParameter("b", "i32").
		WithResult("i32")

	md.FnExports.AddFunction("sayHello").
		WithParameter("name", "string").
		WithResult("string")

	md.FnExports.AddFunction("currentTime").
		WithResult("Date")

	md.FnExports.AddFunction("transform").
		WithParameter("items", "Map<string,string>").
		WithResult("Map<string,string>")

	md.FnExports.AddFunction("testDefaultIntParams").
		WithParameter("a", "i32").
		WithParameter("b", "i32", 0).
		WithParameter("c", "i32", 1)

	md.FnExports.AddFunction("testDefaultStringParams").
		WithParameter("a", "string").
		WithParameter("b", "string", "").
		WithParameter("c", "string", `a"b`).
		WithParameter("d", "string|null").
		WithParameter("e", "string|null", nil).
		WithParameter("f", "string|null", "").
		WithParameter("g", "string|null", "test")

	md.FnExports.AddFunction("testDefaultArrayParams").
		WithParameter("a", "i32[]").
		WithParameter("b", "i32[]", []int32{}).
		WithParameter("c", "i32[]", []int32{1, 2, 3}).
		WithParameter("d", "i32[]|null").
		WithParameter("e", "i32[]|null", nil).
		WithParameter("f", "i32[]|null", []int32{}).
		WithParameter("g", "i32[]|null", []int32{1, 2, 3})

	md.FnExports.AddFunction("getPerson").
		WithResult("Person")

	md.FnExports.AddFunction("getPeople").
		WithResult("Person[]")

	md.FnExports.AddFunction("getProductMap").
		WithResult("Map<string,Product>")

	md.FnExports.AddFunction("doNothing")

	// This should be excluded from the final schema
	md.FnExports.AddFunction("myEmbedder").
		WithParameter("text", "string").
		WithResult("f64[]")

	md.Types.AddType("Company").
		WithField("name", "string")

	md.Types.AddType("Product").
		WithField("name", "string").
		WithField("price", "f64").
		WithField("manufacturer", "Company").
		WithField("components", "Product[]")

	md.Types.AddType("Person").
		WithField("name", "string").
		WithField("age", "i32").
		WithField("addresses", "Address[]")

	md.Types.AddType("Person[]")

	md.Types.AddType("Address").
		WithField("street", "string").
		WithField("city", "string").
		WithField("state", "string").
		WithField("country", "string").
		WithField("postalCode", "string").
		WithField("location", "Coordinates")

	md.Types.AddType("Coordinates").
		WithField("lat", "f64").
		WithField("lon", "f64")

	// This should be excluded from the final schema
	md.Types.AddType("Header").
		WithField("name", "string").
		WithField("values", "string[]")

	result, err := GetGraphQLSchema(context.Background(), md)

	expectedSchema := `
# Hypermode GraphQL Schema (auto-generated)

type Query {
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
`[1:]

	require.Nil(t, err)
	require.Equal(t, expectedSchema, result.Schema)
}

func Test_ConvertType_AssemblyScript(t *testing.T) {

	lti := languages.AssemblyScript().TypeInfo()

	testCases := []struct {
		inputType          string
		expectedOutputType string
		inputTypeDefs      []*metadata.TypeDefinition
		expectedTypeDefs   []*TypeDefinition
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
		{"(string|null)[]", "[String]!", nil, nil},

		// Custom scalar types
		{"Date", "Timestamp!", nil, []*TypeDefinition{{Name: "Timestamp"}}},
		{"i64", "Int64!", nil, []*TypeDefinition{{Name: "Int64"}}},
		{"u32", "UInt!", nil, []*TypeDefinition{{Name: "UInt"}}},
		{"u64", "UInt64!", nil, []*TypeDefinition{{Name: "UInt64"}}},

		// Custom types
		{"User", "User!",
			[]*metadata.TypeDefinition{{
				Name: "User",
				Fields: []*metadata.Field{
					{Name: "firstName", Type: "string"},
					{Name: "lastName", Type: "string"},
					{Name: "age", Type: "u8"},
				},
			}},
			[]*TypeDefinition{{
				Name: "User",
				Fields: []*NameTypePair{
					{"firstName", "String!"},
					{"lastName", "String!"},
					{"age", "Int!"},
				},
			}}},

		// bool and numeric types can't be nullable in AssemblyScript
		// but string and custom types can
		{"string|null", "String", nil, nil},
		{"Foo|null", "Foo",
			[]*metadata.TypeDefinition{{Name: "Foo"}},
			[]*TypeDefinition{{Name: "Foo"}}},

		// Map types
		{"Map<string,string>", "[StringStringPair!]!", nil, []*TypeDefinition{{
			Name: "StringStringPair",
			Fields: []*NameTypePair{
				{"key", "String!"},
				{"value", "String!"},
			},
			IsMapType: true,
		}}},
		{"Map<string,string|null>", "[StringNullableStringPair!]!", nil, []*TypeDefinition{{
			Name: "StringNullableStringPair",
			Fields: []*NameTypePair{
				{"key", "String!"},
				{"value", "String"},
			},
			IsMapType: true,
		}}},
		{"Map<i32,string>", "[IntStringPair!]!", nil, []*TypeDefinition{{
			Name: "IntStringPair",
			Fields: []*NameTypePair{
				{"key", "Int!"},
				{"value", "String!"},
			},
			IsMapType: true,
		}}},
		{"Map<string,Map<string,f32>>", "[StringStringFloatPairListPair!]!", nil, []*TypeDefinition{
			{
				Name: "StringStringFloatPairListPair",
				Fields: []*NameTypePair{
					{"key", "String!"},
					{"value", "[StringFloatPair!]!"},
				},
				IsMapType: true,
			},
			{
				Name: "StringFloatPair",
				Fields: []*NameTypePair{
					{"key", "String!"},
					{"value", "Float!"},
				},
				IsMapType: true,
			},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.inputType, func(t *testing.T) {

			inputTypes := make(metadata.TypeMap, len(tc.inputTypeDefs))
			for _, td := range tc.inputTypeDefs {
				inputTypes[td.Name] = td
			}

			typeDefs := make(map[string]*TypeDefinition, len(tc.inputTypeDefs))
			errors := transformTypes(inputTypes, typeDefs, lti)
			require.Empty(t, errors)

			result, err := convertType(tc.inputType, lti, typeDefs, false)

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
