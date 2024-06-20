/*
 * Copyright 2024 Hypermode, Inc.
 */

package plugins

import (
	"strings"
	"time"

	"hmruntime/utils"

	"github.com/tetratelabs/wazero"
)

type Plugin struct {
	Module   *wazero.CompiledModule
	Metadata PluginMetadata
	FileName string
	Types    map[string]TypeDefinition
}

type PluginMetadata struct {
	Plugin    string              `json:"plugin"`
	Library   string              `json:"library"`
	BuildId   string              `json:"buildId"`
	BuildTime time.Time           `json:"buildTs"`
	GitRepo   string              `json:"gitRepo"`
	GitCommit string              `json:"gitCommit"`
	Functions []FunctionSignature `json:"functions"`
	Types     []TypeDefinition    `json:"types"`
}

type FunctionSignature struct {
	Name       string      `json:"name"`
	Parameters []Parameter `json:"parameters"`
	ReturnType TypeInfo    `json:"returnType"`
}

func (f *FunctionSignature) String() string {
	b := strings.Builder{}
	b.WriteString(f.Name)
	b.WriteString("(")
	for i, p := range f.Parameters {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(p.Name)
		b.WriteString(": ")
		b.WriteString(p.Type.Name)
	}
	b.WriteString("): ")
	b.WriteString(f.ReturnType.Name)
	return b.String()
}

func (f *FunctionSignature) Signature() string {
	b := strings.Builder{}
	b.WriteString("(")
	for i, p := range f.Parameters {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(p.Type.Name)
	}
	b.WriteString("):")
	b.WriteString(f.ReturnType.Name)
	return b.String()
}

type TypeDefinition struct {
	Id     uint32  `json:"id"`
	Size   uint32  `json:"size"`
	Path   string  `json:"path"`
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type Parameter struct {
	Name string   `json:"name"`
	Type TypeInfo `json:"type"`
}

type Field struct {
	Offset uint32   `json:"offset"`
	Name   string   `json:"name"`
	Type   TypeInfo `json:"type"`
}

type TypeInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
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
	name, _ := utils.ParseNameAndVersion(p.Metadata.Plugin)
	return name
}

func (p *Plugin) BuildId() string {
	return p.Metadata.BuildId
}

func (p *Plugin) Language() PluginLanguage {
	libName, _ := utils.ParseNameAndVersion(p.Metadata.Library)
	switch libName {
	case "@hypermode/functions-as":
		return AssemblyScript
	case "hypermodeAI/functions-go":
		return GoLang
	default:
		return UnknownLanguage
	}
}
