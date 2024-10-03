/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifest

import (
	"encoding/json"

	v1_manifest "github.com/hypermodeinc/manifest/compat/v1"
)

const (
	// for backward compatibility
	V1AuthHeaderVariableName = "__V1_AUTH_HEADER_VALUE__"
)

func parseManifestJsonV1(data []byte, manifest *HypermodeManifest) error {
	// Parse the v1 manifest
	var v1_man v1_manifest.HypermodeManifest
	if err := json.Unmarshal(data, &v1_man); err != nil {
		return err
	}

	manifest.Version = 1

	// Copy the v1 models to the current structure.
	manifest.Models = make(map[string]ModelInfo, len(v1_man.Models))
	for _, model := range v1_man.Models {
		manifest.Models[model.Name] = ModelInfo{
			Name:        model.Name,
			SourceModel: model.SourceModel,
			Provider:    model.Provider,
			Host:        model.Host,
		}
	}

	// Copy the v1 hosts to the current structure.
	manifest.Hosts = make(map[string]HostInfo, len(v1_man.Hosts))
	for _, host := range v1_man.Hosts {
		h := HTTPHostInfo{
			Name: host.Name,
			// In v1 the endpoint was used for both endpoint and baseURL purposes.
			// We'll retain that behavior here so the usage doesn't need to change in the Runtime.
			Endpoint: host.Endpoint,
			BaseURL:  host.Endpoint,
		}
		if host.AuthHeader != "" {
			h.Headers = map[string]string{
				// Use a special variable name for the old auth header value.
				// The runtime will replace this with the old auth header secret value if it exists.
				host.AuthHeader: "{{" + V1AuthHeaderVariableName + "}}",
			}
		}
		manifest.Hosts[host.Name] = h
	}

	return nil
}
