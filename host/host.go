/*
 * Copyright 2023 Hypermode, Inc.
 */

package host

import "github.com/tetratelabs/wazero"

// runtime instance for the WASM modules
var WasmRuntime wazero.Runtime

// map that holds the compiled modules for each plugin
var CompiledModules = make(map[string]wazero.CompiledModule)

// Channel used to signal that registration is needed
var RegistrationRequest chan bool = make(chan bool)

// HypermodeJson holds the hypermode.json data
var HypermodeJson HypermodeJsonStruct = HypermodeJsonStruct{}

// struct that holds the hypermode.json data
type HypermodeJsonStruct struct {
	ModelSpecs           []ModelSpec           `json:"modelSpecs"`
	EmbeddingSpecs       []EmbeddingSpec       `json:"embeddingSpecs"`
	TrainingInstructions []TrainingInstruction `json:"trainingInstructions"`
}

type ModelSpec struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Endpoint   string `json:"endpoint"`
	AuthHeader string `json:"authHeader"`
	ApiKey     string `json:"apiKey"`
}

type EmbeddingSpec struct {
	EntityType string `json:"entity-type"`
	Attribute  string `json:"attribute"`
	ModelID    string `json:"modelName"`
	Config     struct {
		Query    string `json:"query"`
		Template string `json:"template"`
	} `json:"config"`
}

type TrainingInstruction struct {
	ModelName string   `json:"modelName"`
	Labels    []string `json:"labels"`
}
