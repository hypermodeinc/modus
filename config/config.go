/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

import (
	"flag"
	"time"
)

var Port int
var DgraphUrl string
var PluginsPath string
var NoReload bool
var S3Bucket string
var RefreshInterval time.Duration

// HypermodeJson holds the hypermode.json data
var HypermodeData HypermodeAppData = HypermodeAppData{}

// struct that holds the hypermode.json data
type HypermodeAppData struct {
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

func ParseCommandLineFlags() {
	flag.IntVar(&Port, "port", 8686, "The HTTP port to listen on.")
	flag.StringVar(&DgraphUrl, "dgraph", "http://localhost:8080", "The Dgraph url to connect to.")
	flag.StringVar(&PluginsPath, "plugins", "", "The path to the plugins directory.")
	flag.StringVar(&PluginsPath, "plugin", "", "alias for -plugins")
	flag.BoolVar(&NoReload, "noreload", false, "Disable automatic plugin reloading.")
	flag.StringVar(&S3Bucket, "s3bucket", "", "The S3 bucket to use, if using AWS for plugin storage.")
	flag.DurationVar(&RefreshInterval, "refresh", time.Second*5, "The refresh interval to check for plugins and schema changes.")

	flag.Parse()
}
