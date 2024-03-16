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
	FilePath string
}

type PluginMetadata struct {
	Name           string
	Version        string
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
	return p.Metadata.Name
}

func (p *Plugin) Language() PluginLanguage {
	switch p.Metadata.LibraryName {
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

func getPluginMetadata(cm *wazero.CompiledModule) (metadata PluginMetadata, found bool) {
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
			metadata.LibraryName, metadata.LibraryVersion = parseNameAndVersion(string(data))
			found = true
		case "hypermode_plugin":
			metadata.Name, metadata.Version = parseNameAndVersion(string(data))
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
