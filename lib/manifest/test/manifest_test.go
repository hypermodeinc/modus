/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest_test

import (
	_ "embed"
	"reflect"
	"testing"

	"github.com/hypermodeinc/modus/lib/manifest"
)

//go:embed valid_modus.json
var validManifest []byte

func TestReadManifest(t *testing.T) {
	// This should match the content of valid_modus.json
	expectedManifest := &manifest.Manifest{
		Version: 3,
		Endpoints: map[string]manifest.EndpointInfo{
			"default": manifest.GraphqlEndpointInfo{
				Name: "default",
				Type: manifest.EndpointTypeGraphQL,
				Path: "/graphql",
				Auth: manifest.EndpointAuthBearerToken,
			},
		},
		Models: map[string]manifest.ModelInfo{
			"model-1": {
				Name:        "model-1",
				SourceModel: "example/source-model-1",
				Provider:    "hugging-face",
				Connection:  "hypermode",
			},
			"model-2": {
				Name:        "model-2",
				SourceModel: "source-model-2",
				Connection:  "my-model-connection",
				Path:        "path/to/model-2",
			},
			"model-3": {
				Name:        "model-3",
				SourceModel: "source-model-3",
				Connection:  "my-model-connection",
			},
		},
		Connections: map[string]manifest.ConnectionInfo{
			"my-model-connection": manifest.HTTPConnectionInfo{
				Name:    "my-model-connection",
				Type:    manifest.ConnectionTypeHTTP,
				BaseURL: "https://models.example.com/",
				Headers: map[string]string{
					"X-API-Key": "{{API_KEY}}",
				},
			},
			"another-model-connection": manifest.HTTPConnectionInfo{
				Name:     "another-model-connection",
				Type:     manifest.ConnectionTypeHTTP,
				Endpoint: "https://models.example.com/full/path/to/model-3",
				Headers: map[string]string{
					"X-API-Key": "{{API_KEY}}",
				},
			},
			"my-graphql-api": manifest.HTTPConnectionInfo{
				Name:     "my-graphql-api",
				Type:     manifest.ConnectionTypeHTTP,
				Endpoint: "https://api.example.com/graphql",
				Headers: map[string]string{
					"Authorization": "Bearer {{AUTH_TOKEN}}",
				},
			},
			"my-rest-api": manifest.HTTPConnectionInfo{
				Name:    "my-rest-api",
				Type:    manifest.ConnectionTypeHTTP,
				BaseURL: "https://api.example.com/v1/",
				QueryParameters: map[string]string{
					"api_token": "{{API_TOKEN}}",
				},
			},
			"another-rest-api": manifest.HTTPConnectionInfo{
				Name:    "another-rest-api",
				Type:    manifest.ConnectionTypeHTTP,
				BaseURL: "https://api.example.com/v2/",
				Headers: map[string]string{
					"Authorization": "Basic {{base64(USERNAME:PASSWORD)}}",
				},
			},
			"neon": manifest.PostgresqlConnectionInfo{
				Name:    "neon",
				Type:    manifest.ConnectionTypePostgresql,
				ConnStr: "postgresql://{{POSTGRESQL_USERNAME}}:{{POSTGRESQL_PASSWORD}}@1.2.3.4:5432/data?sslmode=disable",
			},
			"my-mysql": manifest.MysqlConnectionInfo{
				Name:    "my-mysql",
				Type:    manifest.ConnectionTypeMysql,
				ConnStr: "mysql://{{MYSQL_USERNAME}}:{{MYSQL_PASSWORD}}@1.2.3.4:3306/mydb?sslmode=disable",
			},
			"my-dgraph-cloud": manifest.DgraphConnectionInfo{
				Name:       "my-dgraph-cloud",
				Type:       manifest.ConnectionTypeDgraph,
				GrpcTarget: "frozen-mango.grpc.eu-central-1.aws.cloud.dgraph.io:443",
				Key:        "{{DGRAPH_KEY}}",
			},
			"local-dgraph": manifest.DgraphConnectionInfo{
				Name:       "local-dgraph",
				Type:       manifest.ConnectionTypeDgraph,
				GrpcTarget: "localhost:9080",
				Key:        "",
			},
			"dgraph-with-connstr": manifest.DgraphConnectionInfo{
				Name:    "dgraph-with-connstr",
				Type:    manifest.ConnectionTypeDgraph,
				ConnStr: "dgraph://localhost:9080?sslmode=disable",
			},
			"my-neo4j": manifest.Neo4jConnectionInfo{
				Name:     "my-neo4j",
				Type:     manifest.ConnectionTypeNeo4j,
				DbUri:    "bolt://localhost:7687",
				Username: "{{NEO4J_USERNAME}}",
				Password: "{{NEO4J_PASSWORD}}",
			},
		},
		Collections: map[string]manifest.CollectionInfo{
			"collection1": {
				Name: "collection1",
				SearchMethods: map[string]manifest.SearchMethodInfo{
					"searchMethod1": {
						Embedder: "embedder1",
					},
					"searchMethod2": {
						Embedder: "embedder1",
						Index: manifest.IndexInfo{
							Type: "hnsw",
							Options: manifest.OptionsInfo{
								EfConstruction: 100,
								MaxLevels:      3,
							},
						},
					},
				},
			},
		},
	}

	actualManifest, err := manifest.ReadManifest(validManifest)
	if err != nil {
		t.Errorf("Error reading manifest: %v", err)
		return
	}

	if !reflect.DeepEqual(actualManifest, expectedManifest) {
		t.Errorf("Expected manifest: %+v, but got: %+v", expectedManifest, actualManifest)
	}
}

func TestValidateManifest(t *testing.T) {
	if err := manifest.ValidateManifest(validManifest); err != nil {
		t.Error(err)
	}
}

func TestModelInfo_Hash(t *testing.T) {
	model := manifest.ModelInfo{
		Name:        "my-model",
		SourceModel: "my-source-model",
		Provider:    "my-provider",
		Connection:  "my-connection",
		Path:        "path/to/my-model",
	}

	expectedHash := "180de1ebfbaeb038682a121cec743809d22e07c8dcedbfcfe14b9b3329774aea"

	if actualHash := model.Hash(); actualHash != expectedHash {
		t.Errorf("Expected hash: %s, but got: %s", expectedHash, actualHash)
	}
}

func TestHttpConnectionInfo_Hash(t *testing.T) {
	connection := manifest.HTTPConnectionInfo{
		Name:     "my-connection",
		Endpoint: "https://example.com/api",
		BaseURL:  "https://example.com/api",
		Headers: map[string]string{
			"Authorization": "Bearer {{API_TOKEN}}",
		},
		QueryParameters: map[string]string{
			"api_token": "{{API_TOKEN}}",
		},
	}

	expectedHash := "dfb26f6fc9ae7b55d0c60b0b1b78518644c397469624a77b063330ad05a43974"

	actualHash := connection.Hash()
	if actualHash != expectedHash {
		t.Errorf("Expected hash: %s, but got: %s", expectedHash, actualHash)
	}
}

func TestPostgresConnectionInfo_Hash(t *testing.T) {
	connection := manifest.PostgresqlConnectionInfo{
		Name:    "my-database",
		ConnStr: "postgresql://{{USERNAME}}:{{PASSWORD}}@database.example.com:5432/dbname?sslmode=require",
	}

	expectedHash := "bca6ac337c06274506199ceab0f5ce6145b2abd88504284ad0125b8ea1e1d051"

	actualHash := connection.Hash()
	if actualHash != expectedHash {
		t.Errorf("Expected hash: %s, but got: %s", expectedHash, actualHash)
	}
}

func TestMysqlConnectionInfo_Hash(t *testing.T) {
	connection := manifest.MysqlConnectionInfo{
		Name:    "my-database",
		ConnStr: "mysql://{{MYSQL_USERNAME}}:{{MYSQL_PASSWORD}}@1.2.3.4:3306/mydb?sslmode=disable",
	}

	expectedHash := "3b96055cec5bd4195901e1442c856fe5b5493b0af0dde8f64f1d14a4795f5272"

	actualHash := connection.Hash()
	if actualHash != expectedHash {
		t.Errorf("Expected hash: %s, but got: %s", expectedHash, actualHash)
	}
}

func TestDgraphCloudConnectionInfo_Hash(t *testing.T) {
	connection := manifest.DgraphConnectionInfo{
		Name:       "my-dgraph-cloud",
		GrpcTarget: "frozen-mango.grpc.eu-central-1.aws.cloud.dgraph.io:443",
		Key:        "{{DGRAPH_KEY}}",
	}

	expectedHash := "542e5cf68cbff1b2839c2494da557d81c0e0c75a8313bedb06e8040dc4973658"
	actualHash := connection.Hash()
	if actualHash != expectedHash {
		t.Errorf("Expected hash: %s, but got: %s", expectedHash, actualHash)
	}
}

func TestDgraphLocalConnectionInfo_Hash(t *testing.T) {
	connection := manifest.DgraphConnectionInfo{
		Name:       "local-dgraph",
		GrpcTarget: "localhost:9080",
		Key:        "",
	}

	expectedHash := "9e5fd654b1007e1eb2d32480c17e5f977046fdda53f649df166f68bf6545ebdc"
	actualHash := connection.Hash()
	if actualHash != expectedHash {
		t.Errorf("Expected hash: %s, but got: %s", expectedHash, actualHash)
	}
}

func TestNeo4jConnectionInfo_Hash(t *testing.T) {
	connection := manifest.Neo4jConnectionInfo{
		Name:     "my-neo4j",
		DbUri:    "bolt://localhost:7687",
		Username: "{{NEO4J_USERNAME}}",
		Password: "{{NEO4J_PASSWORD}}",
	}

	expectedHash := "51a373d6c2e32442d84fae02a0c87ddb24ec08f05260e49dad14718eca057b29"
	actualHash := connection.Hash()
	if actualHash != expectedHash {
		t.Errorf("Expected hash: %s, but got: %s", expectedHash, actualHash)
	}
}

func TestGetVariablesFromManifest(t *testing.T) {
	// This should match the connection variables that are present in valid_modus.json
	expectedVars := map[string][]string{
		"my-model-connection":      {"API_KEY"},
		"another-model-connection": {"API_KEY"},
		"my-graphql-api":           {"AUTH_TOKEN"},
		"my-rest-api":              {"API_TOKEN"},
		"another-rest-api":         {"USERNAME", "PASSWORD"},
		"neon":                     {"POSTGRESQL_USERNAME", "POSTGRESQL_PASSWORD"},
		"my-mysql":                 {"MYSQL_USERNAME", "MYSQL_PASSWORD"},
		"my-dgraph-cloud":          {"DGRAPH_KEY"},
		"my-neo4j":                 {"NEO4J_USERNAME", "NEO4J_PASSWORD"},
	}

	m, err := manifest.ReadManifest(validManifest)
	if err != nil {
		t.Errorf("Error reading manifest: %v", err)
	}

	vars := m.GetVariables()
	if !reflect.DeepEqual(vars, expectedVars) {
		t.Errorf("Expected vars: %+v, but got: %+v", expectedVars, vars)
	}
}
