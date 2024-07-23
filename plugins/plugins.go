/*
 * Copyright 2024 Hypermode, Inc.
 */

package plugins

import (
	"strings"
	"time"

	"hmruntime/utils"

	"github.com/buger/jsonparser"
	"github.com/tetratelabs/wazero"
)

type Plugin struct {
	Module   *wazero.CompiledModule
	Metadata PluginMetadata
	FileName string
	Types    map[string]TypeDefinition
}

type PluginMetadata struct {
	Id        string              `json:"-"` // from db when inserted
	Plugin    string              `json:"plugin"`
	SDK       string              `json:"sdk"`
	Library   string              `json:"library"` // deprecated
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
	Name     string   `json:"name"`
	Type     TypeInfo `json:"type"`
	Default  *any     `json:"default"`
	Optional bool     `json:"optional"` // deprecated
}

func (p *Parameter) UnmarshalJSON(data []byte) error {

	// We need to manually unmarshal the JSON to distinguish between a null default
	// value and the absence of a default value.

	name, err := jsonparser.GetString(data, "name")
	if err != nil {
		return err
	}
	p.Name = name

	typeData, _, _, err := jsonparser.Get(data, "type")
	if err != nil {
		return err
	}
	if err := utils.JsonDeserialize(typeData, &p.Type); err != nil {
		return err
	}

	defaultData, dt, _, err := jsonparser.Get(data, "default")
	switch dt {
	case jsonparser.NotExist:
		// no default value
		p.Default = nil
	case jsonparser.Null:
		// an explicit null default value
		p.Default = new(any)
	case jsonparser.String:
		// a default value that is a string
		s, err := jsonparser.ParseString(defaultData)
		if err != nil {
			return err
		}
		def := any(s)
		p.Default = &def
	default:
		// some other non-null default value
		if err != nil {
			return err
		}
		if err := utils.JsonDeserialize(defaultData, &p.Default); err != nil {
			return err
		}
	}

	return nil
}

type Field struct {
	Offset uint32   `json:"offset"`
	Name   string   `json:"name"`
	Type   TypeInfo `json:"type"`
}

type TypeInfo struct {
	Name     string         `json:"name"`
	Path     string         `json:"path"`
	Language PluginLanguage `json:"-"`
	Nullable bool           `json:"-"`
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

func (p *Plugin) NameAndVersion() (name string, version string) {
	return p.Metadata.NameAndVersion()
}

func (p *Plugin) Name() string {
	return p.Metadata.Name()
}

func (p *Plugin) Version() string {
	return p.Metadata.Version()
}

func (p *Plugin) BuildId() string {
	return p.Metadata.BuildId
}

func (p *Plugin) Language() PluginLanguage {
	return p.Metadata.Language()
}

func (m *PluginMetadata) Language() PluginLanguage {
	switch m.SdkName() {
	case "functions-as":
		return AssemblyScript
	case "functions-go":
		return GoLang
	default:
		return UnknownLanguage
	}
}

func (m *PluginMetadata) NameAndVersion() (name string, version string) {
	return parseNameAndVersion(m.Plugin)
}

func (m *PluginMetadata) Name() string {
	name, _ := m.NameAndVersion()
	return name
}

func (m *PluginMetadata) Version() string {
	_, version := m.NameAndVersion()
	return version
}

func (m *PluginMetadata) SdkNameAndVersion() (name string, version string) {
	return parseNameAndVersion(m.SDK)
}

func (m *PluginMetadata) SdkName() string {
	name, _ := m.SdkNameAndVersion()
	return name
}

func (m *PluginMetadata) SdkVersion() string {
	_, version := m.SdkNameAndVersion()
	return version
}

func parseNameAndVersion(s string) (name string, version string) {
	i := strings.LastIndex(s, "@")
	if i == -1 {
		return s, ""
	}
	return s[:i], s[i+1:]
}
