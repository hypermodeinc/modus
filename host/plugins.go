/*
 * Copyright 2024 Hypermode, Inc.
 */

package host

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/tetratelabs/wazero"
)

type Plugin struct {
	Module   *wazero.CompiledModule
	Metadata PluginMetadata
	FileName string
}

type PluginMetadata struct {
	Plugin    string              `json:"plugin"`
	Library   string              `json:"library"`
	BuildId   string              `json:"buildId"`
	BuildTime time.Time           `json:"buildTs"`
	GitRepo   string              `json:"gitRepo"`
	GitCommit string              `json:"gitCommit"`
	Functions []FunctionSignature `json:"functions"`
}

type FunctionSignature struct {
	Name       string              `json:"name"`
	Parameters []FunctionParameter `json:"parameters"`
	ReturnType string              `json:"returnType"`
}

type FunctionParameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type PluginLanguage int

const (
	UnknownLanguage PluginLanguage = iota
	AssemblyScript
	GoLang
)

func (lang PluginLanguage) String() string {
	switch lang {
	case AssemblyScript:
		return "AssemblyScript"
	case GoLang:
		return "Go"
	default:
		return "Unknown"
	}
}

func (p *Plugin) Name() string {
	name, _ := parseNameAndVersion(p.Metadata.Plugin)
	return name
}

func (p *Plugin) BuildId() string {
	return p.Metadata.BuildId
}

func (p *Plugin) Language() PluginLanguage {
	libName, _ := parseNameAndVersion(p.Metadata.Library)
	switch libName {
	case "@hypermode/functions-as":
		return AssemblyScript
	case "gohypermode/functions-go":
		return GoLang
	default:
		return UnknownLanguage
	}
}

func parseNameAndVersion(s string) (name string, version string) {
	i := strings.LastIndex(s, "@")
	if i == -1 {
		return s, ""
	}
	return s[:i], s[i+1:]
}

func getCustomSectionData(cm *wazero.CompiledModule, name string) (data []byte, found bool) {
	for _, sec := range (*cm).CustomSections() {
		if sec.Name() == name {
			data = sec.Data()
			found = true
			break
		}
	}
	return data, found
}

var errPluginMetadataNotFound = fmt.Errorf("no metadata found in plugin")

func getPluginMetadata(cm *wazero.CompiledModule) (PluginMetadata, error) {
	metadataJson, found := getCustomSectionData(cm, "hypermode_meta")
	if !found {
		return PluginMetadata{}, errPluginMetadataNotFound
	}

	metadata := PluginMetadata{}
	err := json.Unmarshal(metadataJson, &metadata)
	if err != nil {
		return PluginMetadata{}, fmt.Errorf("failed to parse plugin metadata: %w", err)
	}

	return metadata, nil
}

func getPluginMetadata_old(cm *wazero.CompiledModule) (metadata PluginMetadata, found bool) {
	for _, sec := range (*cm).CustomSections() {
		name := sec.Name()
		data := sec.Data()

		switch name {
		case "build_id":
			metadata.BuildId = string(data)
			found = true
		case "build_ts":
			metadata.BuildTime, _ = time.Parse(time.RFC3339, string(data))
			found = true
		case "hypermode_library":
			metadata.Library = string(data)
			found = true
		case "hypermode_plugin":
			metadata.Plugin = string(data)
			found = true
		case "git_repo":
			metadata.GitRepo = string(data)
			found = true
		case "git_commit":
			metadata.GitCommit = string(data)
			found = true
		}
	}

	return metadata, found
}
