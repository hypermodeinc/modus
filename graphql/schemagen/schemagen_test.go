/*
 * Copyright 2024 Hypermode, Inc.
 */

package schemagen

import (
	"context"
	"testing"

	"hypruntime/languages"
	"hypruntime/manifestdata"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"

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
		WithParameter("name", "~lib/string/String").
		WithResult("~lib/string/String")

	md.FnExports.AddFunction("currentTime").
		WithResult("~lib/date/Date")

	md.FnExports.AddFunction("transform").
		WithParameter("items", "~lib/map/Map<~lib/string/String,~lib/string/String>").
		WithResult("~lib/map/Map<~lib/string/String,~lib/string/String>")

	md.FnExports.AddFunction("testDefaultIntParams").
		WithParameter("a", "i32").
		WithParameter("b", "i32", 0).
		WithParameter("c", "i32", 1)

	md.FnExports.AddFunction("testDefaultStringParams").
		WithParameter("a", "~lib/string/String").
		WithParameter("b", "~lib/string/String", "").
		WithParameter("c", "~lib/string/String", `a"b`).
		WithParameter("d", "~lib/string/String | null").
		WithParameter("e", "~lib/string/String | null", nil).
		WithParameter("f", "~lib/string/String | null", "").
		WithParameter("g", "~lib/string/String | null", "test")

	md.FnExports.AddFunction("testDefaultArrayParams").
		WithParameter("a", "~lib/array/Array<i32>").
		WithParameter("b", "~lib/array/Array<i32>", []int32{}).
		WithParameter("c", "~lib/array/Array<i32>", []int32{1, 2, 3}).
		WithParameter("d", "~lib/array/Array<i32> | null").
		WithParameter("e", "~lib/array/Array<i32> | null", nil).
		WithParameter("f", "~lib/array/Array<i32> | null", []int32{}).
		WithParameter("g", "~lib/array/Array<i32> | null", []int32{1, 2, 3})

	md.FnExports.AddFunction("getPerson").
		WithResult("assembly/test/Person")

	md.FnExports.AddFunction("getPeople").
		WithResult("~lib/array/Array<assembly/test/Person>")

	md.FnExports.AddFunction("getProductMap").
		WithResult("~lib/map/Map<~lib/string/String,assembly/test/Product>")

	md.FnExports.AddFunction("doNothing")

	// This should be excluded from the final schema
	md.FnExports.AddFunction("myEmbedder").
		WithParameter("text", "~lib/string/String").
		WithResult("~lib/array/Array<f64>")

	md.Types.AddType("assembly/test/Company").
		WithField("name", "~lib/string/String")

	md.Types.AddType("assembly/test/Product").
		WithField("name", "~lib/string/String").
		WithField("price", "f64").
		WithField("manufacturer", "assembly/test/Company").
		WithField("components", "~lib/array/Array<assembly/test/Product>")

	md.Types.AddType("assembly/test/Person").
		WithField("name", "~lib/string/String").
		WithField("age", "i32").
		WithField("addresses", "~lib/array/Array<assembly/test/Address>")

	md.Types.AddType("~lib/array/Array<assembly/test/Person>")

	md.Types.AddType("assembly/test/Address").
		WithField("street", "~lib/string/String").
		WithField("city", "~lib/string/String").
		WithField("state", "~lib/string/String").
		WithField("country", "~lib/string/String").
		WithField("postalCode", "~lib/string/String").
		WithField("location", "assembly/test/Coordinates")

	md.Types.AddType("assembly/test/Coordinates").
		WithField("lat", "f64").
		WithField("lon", "f64")

	// This should be excluded from the final schema
	md.Types.AddType("assembly/test/Header").
		WithField("name", "~lib/string/String").
		WithField("values", "~lib/array/Array<~lib/string/String>")

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
		{"~lib/string/String", "String!", nil, nil},
		{"bool", "Boolean!", nil, nil},
		{"i32", "Int!", nil, nil},
		{"i16", "Int!", nil, nil},
		{"i8", "Int!", nil, nil},
		{"u16", "Int!", nil, nil},
		{"u8", "Int!", nil, nil},
		{"f64", "Float!", nil, nil},
		{"f32", "Float!", nil, nil},

		// Array types
		{"~lib/array/Array<~lib/string/String>", "[String!]!", nil, nil},
		{"~lib/array/Array<~lib/array/Array<~lib/string/String>>", "[[String!]!]!", nil, nil},
		{"~lib/array/Array<~lib/string/String|null>", "[String]!", nil, nil},

		// Custom scalar types
		{"~lib/date/Date", "Timestamp!", nil, []*TypeDefinition{{Name: "Timestamp"}}},
		{"i64", "Int64!", nil, []*TypeDefinition{{Name: "Int64"}}},
		{"u32", "UInt!", nil, []*TypeDefinition{{Name: "UInt"}}},
		{"u64", "UInt64!", nil, []*TypeDefinition{{Name: "UInt64"}}},

		// Custom types
		{"assembly/test/User", "User!",
			[]*metadata.TypeDefinition{{
				Name: "User",
				Fields: []*metadata.Field{
					{Name: "firstName", Type: "~lib/string/String"},
					{Name: "lastName", Type: "~lib/string/String"},
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
		{"~lib/string/String | null", "String", nil, nil},
		{"assembly/test/Foo | null", "Foo",
			[]*metadata.TypeDefinition{{Name: "assembly/test/Foo"}},
			[]*TypeDefinition{{Name: "Foo"}}},

		// Map types
		{"~lib/map/Map<~lib/string/String,~lib/string/String>", "[StringStringPair!]!", nil, []*TypeDefinition{{
			Name: "StringStringPair",
			Fields: []*NameTypePair{
				{"key", "String!"},
				{"value", "String!"},
			},
			IsMapType: true,
		}}},
		{"~lib/map/Map<~lib/string/String,~lib/string/String|null>", "[StringNullableStringPair!]!", nil, []*TypeDefinition{{
			Name: "StringNullableStringPair",
			Fields: []*NameTypePair{
				{"key", "String!"},
				{"value", "String"},
			},
			IsMapType: true,
		}}},
		{"~lib/map/Map<i32,~lib/string/String>", "[IntStringPair!]!", nil, []*TypeDefinition{{
			Name: "IntStringPair",
			Fields: []*NameTypePair{
				{"key", "Int!"},
				{"value", "String!"},
			},
			IsMapType: true,
		}}},
		{"~lib/map/Map<~lib/string/String,~lib/map/Map<~lib/string/String,f32>>", "[StringStringFloatPairListPair!]!", nil, []*TypeDefinition{
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

			typeDefs, errors := transformTypes(inputTypes, lti)
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
