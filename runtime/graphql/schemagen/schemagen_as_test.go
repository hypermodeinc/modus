/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package schemagen

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/lib/metadata"
	"github.com/hypermodeinc/modus/runtime/languages"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/stretchr/testify/require"
)

func Test_GetGraphQLSchema_AssemblyScript(t *testing.T) {

	manifest := &manifest.Manifest{
		Models:      map[string]manifest.ModelInfo{},
		Connections: map[string]manifest.ConnectionInfo{},
	}
	manifestdata.SetManifest(manifest)

	md := metadata.NewPluginMetadata()
	md.SDK = "modus-sdk-as" // AssemblyScript

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
		WithResult("assembly/test/Person").
		WithDocs(metadata.Docs{
			Lines: []string{
				"This is a Person object",
			},
		})

	md.FnExports.AddFunction("listPeople").
		WithResult("~lib/array/Array<assembly/test/Person>")

	md.FnExports.AddFunction("addPerson").
		WithParameter("person", "assembly/test/Person")

	md.FnExports.AddFunction("getProductMap").
		WithResult("~lib/map/Map<~lib/string/String,assembly/test/Product>")

	md.FnExports.AddFunction("doNothing")

	// Generated input object from the output object
	md.FnExports.AddFunction("testObj1").
		WithParameter("obj", "assembly/test/Obj1").
		WithResult("assembly/test/Obj1")
	md.Types.AddType("assembly/test/Obj1").
		WithField("id", "i32").
		WithField("name", "~lib/string/String")

	// Separate input and output objects defined
	md.FnExports.AddFunction("testObj2").
		WithParameter("obj", "assembly/test/Obj2Input").
		WithResult("assembly/test/Obj2")
	md.Types.AddType("assembly/test/Obj2").
		WithField("id", "i32").
		WithField("name", "~lib/string/String")
	md.Types.AddType("assembly/test/Obj2Input").
		WithField("name", "~lib/string/String")

	// Generated input object without output object
	md.FnExports.AddFunction("testObj3").
		WithParameter("obj", "assembly/test/Obj3")
	md.Types.AddType("assembly/test/Obj3").
		WithField("name", "~lib/string/String")

	// Single input object defined without output object
	md.FnExports.AddFunction("testObj4").
		WithParameter("obj", "assembly/test/Obj4Input")
	md.Types.AddType("assembly/test/Obj4Input").
		WithField("name", "~lib/string/String")

	md.Types.AddType("~lib/array/Array<i32>")
	md.Types.AddType("~lib/array/Array<f64>")
	md.Types.AddType("~lib/array/Array<assembly/test/Person>")
	md.Types.AddType("~lib/map/Map<~lib/string/String,~lib/string/String>")
	md.Types.AddType("~lib/map/Map<~lib/string/String,assembly/test/Product>")

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

	md.Types.AddType("assembly/test/Address").
		WithField("street", "~lib/string/String", &metadata.Docs{
			Lines: []string{
				"Street that the user lives on",
			},
		}).
		WithField("city", "~lib/string/String").
		WithField("state", "~lib/string/String").
		WithField("country", "~lib/string/String", &metadata.Docs{
			Lines: []string{
				"Country that the user is from",
			},
		}).
		WithField("postalCode", "~lib/string/String").
		WithField("location", "assembly/test/Coordinates").
		WithDocs(metadata.Docs{
			Lines: []string{
				"Address represents a physical address.",
				"Each field corresponds to a specific part of the address.",
				"The location field stores geospatial coordinates.",
			},
		})
	md.Types.AddType("assembly/test/Coordinates").
		WithField("lat", "f64").
		WithField("lon", "f64")

	// This should be excluded from the final schema
	md.Types.AddType("assembly/test/Header").
		WithField("name", "~lib/string/String").
		WithField("values", "~lib/array/Array<~lib/string/String>")

	result, err := GetGraphQLSchema(t.Context(), md)

	t.Log(result.Schema)

	expectedSchema := `
# Modus GraphQL Schema (auto-generated)

type Query {
  currentTime: Timestamp!
  doNothing: Void
  people: [Person!]!
  """
  This is a Person object
  """
  person: Person!
  productMap: [StringProductPair!]!
  sayHello(name: String!): String!
  testDefaultArrayParams(a: [Int!]!, b: [Int!]! = [], c: [Int!]! = [1,2,3], d: [Int!], e: [Int!] = null, f: [Int!] = [], g: [Int!] = [1,2,3]): Void
  testDefaultIntParams(a: Int!, b: Int! = 0, c: Int! = 1): Void
  testDefaultStringParams(a: String!, b: String! = "", c: String! = "a\"b", d: String, e: String = null, f: String = "", g: String = "test"): Void
  testObj1(obj: Obj1Input!): Obj1!
  testObj2(obj: Obj2Input!): Obj2!
  testObj3(obj: Obj3Input!): Void
  testObj4(obj: Obj4Input!): Void
  transform(items: [StringStringPairInput!]!): [StringStringPair!]!
}

type Mutation {
  add(a: Int!, b: Int!): Int!
  addPerson(person: PersonInput!): Void
}

scalar Timestamp
scalar Void

"""
Address represents a physical address.
Each field corresponds to a specific part of the address.
The location field stores geospatial coordinates.
"""
input AddressInput {
  """
  Street that the user lives on
  """
  street: String!
  city: String!
  state: String!
  """
  Country that the user is from
  """
  country: String!
  postalCode: String!
  location: CoordinatesInput!
}

input CoordinatesInput {
  lat: Float!
  lon: Float!
}

input Obj1Input {
  id: Int!
  name: String!
}

input Obj2Input {
  name: String!
}

input Obj3Input {
  name: String!
}

input Obj4Input {
  name: String!
}

input PersonInput {
  name: String!
  age: Int!
  addresses: [AddressInput!]!
}

input StringStringPairInput {
  key: String!
  value: String!
}

"""
Address represents a physical address.
Each field corresponds to a specific part of the address.
The location field stores geospatial coordinates.
"""
type Address {
  """
  Street that the user lives on
  """
  street: String!
  city: String!
  state: String!
  """
  Country that the user is from
  """
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

type Obj1 {
  id: Int!
  name: String!
}

type Obj2 {
  id: Int!
  name: String!
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
		sourceType          string
		forInput            bool
		expectedGraphQLType string
		sourceTypeDefs      []*metadata.TypeDefinition
		expectedTypeDefs    []*TypeDefinition
	}{
		// Plain non-nullable types
		{"~lib/string/String", false, "String!", nil, nil},
		{"~lib/string/String", true, "String!", nil, nil},
		{"bool", false, "Boolean!", nil, nil},
		{"bool", true, "Boolean!", nil, nil},
		{"i8", false, "Int!", nil, nil},
		{"i8", true, "Int!", nil, nil},
		{"i16", false, "Int!", nil, nil},
		{"i16", true, "Int!", nil, nil},
		{"i32", false, "Int!", nil, nil},
		{"i32", true, "Int!", nil, nil},
		{"u8", false, "Int!", nil, nil},
		{"u8", true, "Int!", nil, nil},
		{"u16", false, "Int!", nil, nil},
		{"u16", true, "Int!", nil, nil},
		{"f32", false, "Float!", nil, nil},
		{"f32", true, "Float!", nil, nil},
		{"f64", false, "Float!", nil, nil},
		{"f64", true, "Float!", nil, nil},

		// Array types
		{"~lib/array/Array<~lib/string/String>", false, "[String!]!", nil, nil},
		{"~lib/array/Array<~lib/string/String>", true, "[String!]!", nil, nil},
		{"~lib/array/Array<~lib/array/Array<~lib/string/String>>", false, "[[String!]!]!", nil, nil},
		{"~lib/array/Array<~lib/array/Array<~lib/string/String>>", true, "[[String!]!]!", nil, nil},
		{"~lib/array/Array<~lib/string/String|null>", false, "[String]!", nil, nil},
		{"~lib/array/Array<~lib/string/String|null>", true, "[String]!", nil, nil},

		// Custom scalar types
		{"~lib/date/Date", false, "Timestamp!", nil, []*TypeDefinition{{Name: "Timestamp"}}},
		{"~lib/date/Date", true, "Timestamp!", nil, []*TypeDefinition{{Name: "Timestamp"}}},
		{"i64", false, "Int64!", nil, []*TypeDefinition{{Name: "Int64"}}},
		{"i64", true, "Int64!", nil, []*TypeDefinition{{Name: "Int64"}}},
		{"u32", false, "UInt!", nil, []*TypeDefinition{{Name: "UInt"}}},
		{"u32", true, "UInt!", nil, []*TypeDefinition{{Name: "UInt"}}},
		{"u64", false, "UInt64!", nil, []*TypeDefinition{{Name: "UInt64"}}},
		{"u64", true, "UInt64!", nil, []*TypeDefinition{{Name: "UInt64"}}},

		// Custom types
		{"assembly/test/User", false, "User!",
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
				Fields: []*FieldDefinition{
					{Name: "firstName", Type: "String!"},
					{Name: "lastName", Type: "String!"},
					{Name: "age", Type: "Int!"},
				},
			}}},
		{"assembly/test/User", true, "UserInput!",
			[]*metadata.TypeDefinition{{
				Name: "User",
				Fields: []*metadata.Field{
					{Name: "firstName", Type: "~lib/string/String"},
					{Name: "lastName", Type: "~lib/string/String"},
					{Name: "age", Type: "u8"},
				},
			}},
			[]*TypeDefinition{{
				Name: "UserInput",
				Fields: []*FieldDefinition{
					{Name: "firstName", Type: "String!"},
					{Name: "lastName", Type: "String!"},
					{Name: "age", Type: "Int!"},
				},
			}}},

		// bool and numeric types can't be nullable in AssemblyScript
		// but string and custom types can
		{"~lib/string/String | null", false, "String", nil, nil},
		{"~lib/string/String | null", true, "String", nil, nil},
		{"assembly/test/Foo | null", false, "Foo", // scalar
			[]*metadata.TypeDefinition{{Name: "assembly/test/Foo"}},
			[]*TypeDefinition{{Name: "Foo"}}},
		{"assembly/test/Foo | null", true, "Foo", // scalar
			[]*metadata.TypeDefinition{{Name: "assembly/test/Foo"}},
			[]*TypeDefinition{{Name: "Foo"}}},

		// Map types
		{"~lib/map/Map<~lib/string/String,~lib/string/String>", false, "[StringStringPair!]!", nil, []*TypeDefinition{{
			Name: "StringStringPair",
			Fields: []*FieldDefinition{
				{Name: "key", Type: "String!"},
				{Name: "value", Type: "String!"},
			},
			IsMapType: true,
		}}},
		{"~lib/map/Map<~lib/string/String,~lib/string/String>", true, "[StringStringPairInput!]!", nil, []*TypeDefinition{{
			Name: "StringStringPairInput",
			Fields: []*FieldDefinition{
				{Name: "key", Type: "String!"},
				{Name: "value", Type: "String!"},
			},
			IsMapType: true,
		}}},
		{"~lib/map/Map<~lib/string/String,~lib/string/String|null>", false, "[StringNullableStringPair!]!", nil, []*TypeDefinition{{
			Name: "StringNullableStringPair",
			Fields: []*FieldDefinition{
				{Name: "key", Type: "String!"},
				{Name: "value", Type: "String"},
			},
			IsMapType: true,
		}}},
		{"~lib/map/Map<~lib/string/String,~lib/string/String|null>", true, "[StringNullableStringPairInput!]!", nil, []*TypeDefinition{{
			Name: "StringNullableStringPairInput",
			Fields: []*FieldDefinition{
				{Name: "key", Type: "String!"},
				{Name: "value", Type: "String"},
			},
			IsMapType: true,
		}}},
		{"~lib/map/Map<i32,~lib/string/String>", false, "[IntStringPair!]!", nil, []*TypeDefinition{{
			Name: "IntStringPair",
			Fields: []*FieldDefinition{
				{Name: "key", Type: "Int!"},
				{Name: "value", Type: "String!"},
			},
			IsMapType: true,
		}}},
		{"~lib/map/Map<i32,~lib/string/String>", true, "[IntStringPairInput!]!", nil, []*TypeDefinition{{
			Name: "IntStringPairInput",
			Fields: []*FieldDefinition{
				{Name: "key", Type: "Int!"},
				{Name: "value", Type: "String!"},
			},
			IsMapType: true,
		}}},
		{"~lib/map/Map<~lib/string/String,~lib/map/Map<~lib/string/String,f32>>", false, "[StringStringFloatPairListPair!]!", nil, []*TypeDefinition{
			{
				Name: "StringStringFloatPairListPair",
				Fields: []*FieldDefinition{
					{Name: "key", Type: "String!"},
					{Name: "value", Type: "[StringFloatPair!]!"},
				},
				IsMapType: true,
			},
			{
				Name: "StringFloatPair",
				Fields: []*FieldDefinition{
					{Name: "key", Type: "String!"},
					{Name: "value", Type: "Float!"},
				},
				IsMapType: true,
			},
		}},
		{"~lib/map/Map<~lib/string/String,~lib/map/Map<~lib/string/String,f32>>", true, "[StringStringFloatPairListPairInput!]!", nil, []*TypeDefinition{
			{
				Name: "StringStringFloatPairListPairInput",
				Fields: []*FieldDefinition{
					{Name: "key", Type: "String!"},
					{Name: "value", Type: "[StringFloatPairInput!]!"},
				},
				IsMapType: true,
			},
			{
				Name: "StringFloatPairInput",
				Fields: []*FieldDefinition{
					{Name: "key", Type: "String!"},
					{Name: "value", Type: "Float!"},
				},
				IsMapType: true,
			},
		}},
	}

	nameRegex := regexp.MustCompile(`(?:[~\w]+/)+?`)

	for _, tc := range testCases {
		testName := strings.ReplaceAll(nameRegex.ReplaceAllString(tc.sourceType, ""), " ", "")
		if tc.forInput {
			testName += "_input"
		}
		t.Run(testName, func(t *testing.T) {

			types := make(metadata.TypeMap, len(tc.sourceTypeDefs))
			for _, td := range tc.sourceTypeDefs {
				types[td.Name] = td
			}

			typeDefs, errors := transformTypes(types, lti, tc.forInput)
			require.Empty(t, errors)

			result, err := convertType(tc.sourceType, lti, typeDefs, false, tc.forInput)

			require.Nil(t, err)
			require.Equal(t, tc.expectedGraphQLType, result)

			if tc.expectedTypeDefs == nil {
				require.Empty(t, typeDefs)
			} else {
				require.ElementsMatch(t, tc.expectedTypeDefs, utils.MapValues(typeDefs))
			}
		})
	}
}
