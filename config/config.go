/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var Port int
var DgraphUrl string
var ModelHost string
var StoragePath string
var UseAwsSecrets bool
var UseAwsStorage bool
var S3Bucket string
var S3Path string
var NoReload bool
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
type AppData any

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
	flag.StringVar(&ModelHost, "modelHost", "", "The base DNS of the host endpoint to the model server.")
	flag.StringVar(&StoragePath, "storagePath", getDefaultStoragePath(), "The path to a directory used for local storage.")
	flag.BoolVar(&UseAwsSecrets, "useAwsSecrets", false, "Use AWS Secrets Manager for API keys and other secrets.")
	flag.BoolVar(&UseAwsStorage, "useAwsStorage", false, "Use AWS S3 for storage instead of the local filesystem.")
	flag.StringVar(&S3Bucket, "s3bucket", "", "The S3 bucket to use, if using AWS storage.")
	flag.StringVar(&S3Path, "s3path", "", "The path within the S3 bucket to use, if using AWS storage.")
	flag.BoolVar(&NoReload, "noreload", false, "Disable automatic plugin reloading.")
	flag.DurationVar(&RefreshInterval, "refresh", time.Second*5, "The refresh interval to check for plugins and schema changes.")
	flag.BoolVar(&UseJsonLogging, "jsonlogs", false, "Use JSON format for logging.")

	flag.Parse()
}

func getDefaultStoragePath() string {

	// On Windows, the default is %APPDATA%\Hypermode
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		return filepath.Join(appData, "Hypermode")
	}

	// On Unix and macOS, the default is $HOME/.hypermode
	homedir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homedir, ".hypermode")
}
