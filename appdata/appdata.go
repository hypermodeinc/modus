/*
 * Copyright 2024 Hypermode, Inc.
 */

package appdata

var appDataFiles = map[string]AppData{
	"hypermode.json": &HypermodeData,
	"models.json":    &ModelData,
}

var HypermodeData HypermodeAppData = HypermodeAppData{}
var ModelData ModelsAppData = ModelsAppData{}

type AppData any

type HypermodeAppData struct {
	Models               []Model               `json:"models"`
	Hosts                []Host                `json:"hosts"`
	EmbeddingSpecs       []EmbeddingSpec       `json:"embeddingSpecs"`
	TrainingInstructions []TrainingInstruction `json:"trainingInstructions"`
	AppData
}

type ModelsAppData struct {
	AppData
}

type ModelTask string

const (
	ClassificationTask ModelTask = "classification"
	EmbeddingTask      ModelTask = "embedding"
	GenerationTask     ModelTask = "generation"
)

type Model struct {
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
