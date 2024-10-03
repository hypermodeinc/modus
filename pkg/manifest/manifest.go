/*
 * Copyright 2024 Hypermode, Inc.
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
const currentVersion = 2

//go:embed hypermode.json
var schemaContent string

type HostInfo interface {
	HostName() string
	HostType() string
	GetVariables() []string
	Hash() string
}

type HypermodeManifest struct {
	Version     int                       `json:"-"`
	Models      map[string]ModelInfo      `json:"models"`
	Hosts       map[string]HostInfo       `json:"hosts"`
	Collections map[string]CollectionInfo `json:"collections"`
}

func (m *HypermodeManifest) IsCurrentVersion() bool {
	return m.Version == currentVersion
}

func (m *HypermodeManifest) GetHostVariables() map[string][]string {
	results := make(map[string][]string, len(m.Hosts))

	for _, host := range m.Hosts {
		vars := host.GetVariables()
		if len(vars) > 0 {
			results[host.HostName()] = vars
		}
	}

	return results
}

func IsCurrentVersion(version int) bool {
	return version == currentVersion
}

func ValidateManifest(content []byte) error {
	sch, err := jsonschema.CompileString("hypermode.json", schemaContent)
	if err != nil {
		return err
	}

	content, err = standardizeJSON(content)
	if err != nil {
		return fmt.Errorf("failed to standardize manifest: %w", err)
	}

	var v interface{}
	if err := json.Unmarshal(content, &v); err != nil {
		return fmt.Errorf("failed to deserialize manifest: %w", err)
	}

	if err := sch.Validate(v); err != nil {
		return fmt.Errorf("failed to validate manifest: %w", err)
	}

	return nil
}

func ReadManifest(content []byte) (HypermodeManifest, error) {
	// Create standard JSON before attempting to parse
	var manifest HypermodeManifest
	data, err := standardizeJSON(content)
	if err != nil {
		return manifest, err
	}

	// Try to parse using the current format first
	errParse := parseManifestJson(data, &manifest)
	if errParse == nil {
		return manifest, nil
	}

	// Try the older format if that failed
	if err := parseManifestJsonV1(data, &manifest); err == nil {
		return manifest, nil
	}

	// We should return the error from parsing using the current format
	return manifest, fmt.Errorf("failed to parse manifest: %w", errParse)
}

func parseManifestJson(data []byte, manifest *HypermodeManifest) error {
	var m struct {
		Models      map[string]ModelInfo       `json:"models"`
		Hosts       map[string]json.RawMessage `json:"hosts"`
		Collections map[string]CollectionInfo  `json:"collections"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	manifest.Version = currentVersion

	// Copy map keys to Name fields
	manifest.Models = m.Models
	for key, model := range manifest.Models {
		model.Name = key
		manifest.Models[key] = model
	}

	// parse the hosts
	manifest.Hosts = make(map[string]HostInfo, len(m.Hosts))
	for name, rawHost := range m.Hosts {
		hostType := gjson.GetBytes(rawHost, "type")

		switch hostType.String() {
		case HostTypeHTTP, "":
			var h HTTPHostInfo
			if err := json.Unmarshal(rawHost, &h); err != nil {
				return fmt.Errorf("failed to parse manifest: %w", err)
			}
			h.Name = name
			h.Type = HostTypeHTTP
			manifest.Hosts[name] = h
		case HostTypePostgresql:
			var h PostgresqlHostInfo
			if err := json.Unmarshal(rawHost, &h); err != nil {
				return fmt.Errorf("failed to parse manifest: %w", err)
			}
			h.Name = name
			manifest.Hosts[name] = h
		case HostTypeDgraph:
			var h DgraphHostInfo
			if err := json.Unmarshal(rawHost, &h); err != nil {
				return fmt.Errorf("failed to parse manifest: %w", err)
			}
			h.Name = name
			manifest.Hosts[name] = h
		default:
			return fmt.Errorf("unknown host type: [%s]", hostType.String())
		}
	}

	manifest.Collections = m.Collections

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
