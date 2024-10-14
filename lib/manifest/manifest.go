/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/tailscale/hujson"
	"github.com/tidwall/gjson"
)

// This version should only be incremented if there are major breaking changes
// to the manifest schema.  In general, we don't consider the schema to be versioned,
// from the user's perspective, so this should be rare.
// NOTE: We intentionally do not expose the *current* version number outside this package.
const currentVersion = 1

//go:embed modus_schema.json
var schemaContent string

type Manifest struct {
	Version     int                       `json:"-"`
	Models      map[string]ModelInfo      `json:"models"`
	Connections map[string]ConnectionInfo `json:"connections"`
	Collections map[string]CollectionInfo `json:"collections"`
}

func (m *Manifest) IsCurrentVersion() bool {
	return m.Version == currentVersion
}

func (m *Manifest) GetVariables() map[string][]string {
	results := make(map[string][]string, len(m.Connections))

	for _, c := range m.Connections {
		vars := c.Variables()
		if len(vars) > 0 {
			results[c.ConnectionName()] = vars
		}
	}

	return results
}

func IsCurrentVersion(version int) bool {
	return version == currentVersion
}

func ValidateManifest(content []byte) error {

	sch, err := jsonschema.CompileString("modus.json", schemaContent)
	if err != nil {
		return err
	}

	if content, err := standardizeJSON(content); err != nil {
		return fmt.Errorf("failed to standardize manifest: %w", err)
	} else {
		var v interface{}
		if err := json.Unmarshal(content, &v); err != nil {
			return fmt.Errorf("failed to deserialize manifest: %w", err)
		}
		if err := sch.Validate(v); err != nil {
			return fmt.Errorf("failed to validate manifest: %w", err)
		}
	}

	return nil
}

func ReadManifest(content []byte) (*Manifest, error) {
	data, err := standardizeJSON(content)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := parseManifestJson(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}
	return &manifest, nil
}

func parseManifestJson(data []byte, manifest *Manifest) error {
	var m struct {
		Models      map[string]ModelInfo       `json:"models"`
		Connections map[string]json.RawMessage `json:"connections"`
		Collections map[string]CollectionInfo  `json:"collections"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	manifest.Version = currentVersion
	manifest.Models = m.Models
	manifest.Collections = m.Collections

	// Copy map keys to Name fields
	for key, model := range manifest.Models {
		model.Name = key
		manifest.Models[key] = model
	}
	for key, collection := range manifest.Collections {
		collection.Name = key
		manifest.Collections[key] = collection
	}

	// Parse the connections by type
	manifest.Connections = make(map[string]ConnectionInfo, len(m.Connections))
	for name, rawCon := range m.Connections {
		t := gjson.GetBytes(rawCon, "type")
		if !t.Exists() {
			return fmt.Errorf("missing type for connection [%s]", name)
		}
		conType := t.String()

		switch conType {
		case ConnectionTypeHTTP:
			var info HTTPConnectionInfo
			if err := json.Unmarshal(rawCon, &info); err != nil {
				return fmt.Errorf("failed to parse http connection [%s]: %w", name, err)
			}
			info.Name = name
			manifest.Connections[name] = info
		case ConnectionTypePostgresql:
			var info PostgresqlConnectionInfo
			if err := json.Unmarshal(rawCon, &info); err != nil {
				return fmt.Errorf("failed to parse postgresql connection [%s]: %w", name, err)
			}
			info.Name = name
			manifest.Connections[name] = info
		case ConnectionTypeDgraph:
			var info DgraphConnectionInfo
			if err := json.Unmarshal(rawCon, &info); err != nil {
				return fmt.Errorf("failed to parse dgraph connection [%s]: %w", name, err)
			}
			info.Name = name
			manifest.Connections[name] = info
		default:
			return fmt.Errorf("unknown type [%s] for connection [%s]", conType, name)
		}
	}

	return nil
}

// standardizeJSON removes comments and trailing commas to make the JSON valid
func standardizeJSON(b []byte) ([]byte, error) {
	ast, err := hujson.Parse(b)
	if err != nil {
		return b, err
	}
	ast.Standardize()
	return ast.Pack(), nil
}
