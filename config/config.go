/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

import (
	"flag"
	"sync"
	"time"
)

var Port int
var DgraphUrl string
var PluginsPath string
var NoReload bool
var S3Bucket string
var RefreshInterval time.Duration
var UseJsonLogging bool

var Mu = &sync.Mutex{}

var SupportedJsons = map[string]AppData{
	"hypermode.json": &HypermodeData,
	"models.json":    &ModelData,
}

// HypermodeJson holds the hypermode.json data
var HypermodeData HypermodeAppData = HypermodeAppData{}

// ModelJson holds the models.json data
var ModelData ModelsAppData = ModelsAppData{}

// AppData interface
type AppData interface{}

// struct that holds the hypermode.json data
type HypermodeAppData struct {
	Models               []Model               `json:"models"`
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
	GeneratorTask      ModelTask = "generator"
)

type Model struct {
	Name        string    `json:"name"`
	Task        ModelTask `json:"task"`
	SourceModel string    `json:"sourceModel"`
	Provider    string    `json:"provider"`
	Host        string    `json:"host"`
	Endpoint    string    `json:"endpoint"`
	AuthHeader  string    `json:"authHeader"`
	ApiKey      string    `json:"apiKey"`
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

func ParseCommandLineFlags() {
	flag.IntVar(&Port, "port", 8686, "The HTTP port to listen on.")
	flag.StringVar(&DgraphUrl, "dgraph", "http://localhost:8080", "The Dgraph url to connect to.")
	flag.StringVar(&PluginsPath, "plugins", "", "The path to the plugins directory.")
	flag.StringVar(&PluginsPath, "plugin", "", "alias for -plugins")
	flag.BoolVar(&NoReload, "noreload", false, "Disable automatic plugin reloading.")
	flag.StringVar(&S3Bucket, "s3bucket", "", "The S3 bucket to use, if using AWS for plugin storage.")
	flag.DurationVar(&RefreshInterval, "refresh", time.Second*5, "The refresh interval to check for plugins and schema changes.")
	flag.BoolVar(&UseJsonLogging, "jsonlogs", false, "Use JSON format for logging.")

	flag.Parse()
}
