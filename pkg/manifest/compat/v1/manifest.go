/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifest

type HypermodeManifest struct {
	Models []ModelInfo `json:"models"`
	Hosts  []HostInfo  `json:"hosts"`
}

type ModelTask string

const (
	ClassificationTask ModelTask = "classification"
	EmbeddingTask      ModelTask = "embedding"
	GenerationTask     ModelTask = "generation"
)

type ModelInfo struct {
	Name        string    `json:"name"`
	Task        ModelTask `json:"task"`
	SourceModel string    `json:"sourceModel"`
	Provider    string    `json:"provider"`
	Host        string    `json:"host"`
}

type HostInfo struct {
	Name       string `json:"name"`
	Endpoint   string `json:"endpoint"`
	AuthHeader string `json:"authHeader"`
}
