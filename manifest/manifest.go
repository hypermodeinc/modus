/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifest

var manifestFiles = map[string]Manifest{
	"hypermode.json": &HypermodeData,
}

var HypermodeData HypermodeManifest = HypermodeManifest{}

type Manifest any

type HypermodeManifest struct {
	Models               []Model               `json:"models"`
	Hosts                []Host                `json:"hosts"`
	EmbeddingSpecs       []EmbeddingSpec       `json:"embeddingSpecs"`
	TrainingInstructions []TrainingInstruction `json:"trainingInstructions"`
	Manifest
}

type ModelTask string

const (
	ClassificationTask ModelTask = "classification"
	EmbeddingTask      ModelTask = "embedding"
	GenerationTask     ModelTask = "generation"
)

type Model struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Task        ModelTask `json:"task"`
	SourceModel string    `json:"sourceModel"`
	Provider    string    `json:"provider"`
	Host        string    `json:"host"`
}

type Host struct {
	Name       string `json:"name"`
	Endpoint   string `json:"endpoint"`
	AuthHeader string `json:"authHeader"`
}

type EmbeddingSpec struct {
	EntityType string `json:"entityType"`
	Attribute  string `json:"attribute"`
	ModelName  string `json:"modelName"`
	Config     struct {
		Query    string `json:"query"`
		Template string `json:"template"`
	} `json:"config"`
}

type TrainingInstruction struct {
	ModelName string   `json:"modelName"`
	Labels    []string `json:"labels"`
}
