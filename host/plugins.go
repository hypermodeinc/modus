/*
 * Copyright 2024 Hypermode, Inc.
 */

package host

import (
	"strings"
	"time"

	"github.com/tetratelabs/wazero"
)

type Plugin struct {
	Module   *wazero.CompiledModule
	Metadata PluginMetadata
}

type PluginMetadata struct {
	Name           string
	Version        string
	Language       PluginLanguage
	LibraryName    string
	LibraryVersion string
	BuildId        string
	BuildTime      time.Time
	GitRepo        string
	GitCommit      string
}

type PluginLanguage int

const (
	UnknownLanguage PluginLanguage = iota
	AssemblyScript
	GoLang
)

func getPluginLanguage(libraryName string) PluginLanguage {
	switch libraryName {
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

func getPluginMetadata(cm *wazero.CompiledModule) PluginMetadata {
	var metadata = PluginMetadata{}

	for _, sec := range (*cm).CustomSections() {
		name := sec.Name()
		data := sec.Data()

		switch name {
		case "build_id":
			metadata.BuildId = string(data)
		case "build_ts":
			metadata.BuildTime, _ = time.Parse(time.RFC3339, string(data))
		case "hypermode_library":
			metadata.LibraryName, metadata.LibraryVersion = parseNameAndVersion(string(data))
			metadata.Language = getPluginLanguage(metadata.LibraryName)
		case "hypermode_plugin":
			metadata.Name, metadata.Version = parseNameAndVersion(string(data))
		case "git_repo":
			metadata.GitRepo = string(data)
		case "git_commit":
			metadata.GitCommit = string(data)
		}
	}

	return metadata
}
