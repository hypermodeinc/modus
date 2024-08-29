/*
 * Copyright 2024 Hypermode, Inc.
 */

package metadata

import (
	"fmt"
	"strings"

	"hmruntime/utils"

	"github.com/buger/jsonparser"
)

const MetadataVersion = 2

type TypeMap map[string]*TypeDefinition
type FunctionMap map[string]*Function

type Metadata struct {
	Plugin    string      `json:"plugin"`
	Module    string      `json:"module"`
	SDK       string      `json:"sdk"`
	BuildId   string      `json:"buildId"`
	BuildTime string      `json:"buildTs"`
	GitRepo   string      `json:"gitRepo,omitempty"`
	GitCommit string      `json:"gitCommit,omitempty"`
	FnExports FunctionMap `json:"fnExports,omitempty"`
	FnImports FunctionMap `json:"fnImports,omitempty"`
	Types     TypeMap     `json:"types,omitempty"`
}

type Function struct {
	Name       string       `json:"-"`
	Parameters []*Parameter `json:"parameters,omitempty"`
	Results    []*Result    `json:"results,omitempty"`
}

type TypeDefinition struct {
	Name   string   `json:"-"`
	Id     uint32   `json:"id,omitempty"`
	Fields []*Field `json:"fields,omitempty"`
}

type Parameter struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Default  *any   `json:"default,omitempty"`
	Optional bool   `json:"-"` // deprecated
}

type Result struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (p *Parameter) UnmarshalJSON(data []byte) error {

	// We need to manually unmarshal the JSON to distinguish between a null default
	// value and the absence of a default value.

	name, err := jsonparser.GetString(data, "name")
	if err != nil {
		return err
	}
	p.Name = name

	typ, err := jsonparser.GetString(data, "type")
	if err != nil {
		return err
	}
	p.Type = typ

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

func (m *Metadata) NameAndVersion() (name string, version string) {
	return parseNameAndVersion(m.Plugin)
}

func (m *Metadata) Name() string {
	name, _ := m.NameAndVersion()
	return name
}

func (m *Metadata) Version() string {
	_, version := m.NameAndVersion()
	return version
}

func (m *Metadata) SdkNameAndVersion() (name string, version string) {
	return parseNameAndVersion(m.SDK)
}

func (m *Metadata) SdkName() string {
	name, _ := m.SdkNameAndVersion()
	return name
}

func (m *Metadata) SdkVersion() string {
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

func NewPluginMetadata() *Metadata {
	return &Metadata{
		FnExports: make(FunctionMap),
		FnImports: make(FunctionMap),
		Types:     make(TypeMap),
	}
}

func (m *FunctionMap) AddFunction(name string) *Function {
	f := &Function{
		Name:       name,
		Parameters: make([]*Parameter, 0),
		Results:    make([]*Result, 0),
	}

	(*m)[name] = f
	return f
}

func (m *TypeMap) AddType(name string) *TypeDefinition {
	t := &TypeDefinition{
		Name:   name,
		Fields: make([]*Field, 0),
	}

	(*m)[name] = t
	return t
}

func (f *Function) WithParameter(name string, typ string, dflt ...any) *Function {
	p := &Parameter{Name: name, Type: typ}
	if len(dflt) > 0 {
		p.Default = &dflt[0]
	}
	f.Parameters = append(f.Parameters, p)
	return f
}

func (f *Function) WithResult(typ string) *Function {
	r := &Result{Type: typ}
	f.Results = append(f.Results, r)
	return f
}

func (f *Function) WithNamedResult(name string, typ string) *Function {
	r := &Result{
		Name: name,
		Type: typ,
	}
	f.Results = append(f.Results, r)
	return f
}

func (t *TypeDefinition) WithId(id uint32) *TypeDefinition {
	t.Id = id
	return t
}

func (t *TypeDefinition) WithField(name string, typ string) *TypeDefinition {
	f := &Field{Name: name, Type: typ}
	t.Fields = append(t.Fields, f)
	return t
}

func (f *Function) String() string {
	p := strings.Trim(fmt.Sprintf("%v", f.Parameters), "[]")
	r := strings.Trim(fmt.Sprintf("%v", f.Results), "[]")
	return strings.TrimSpace(fmt.Sprintf("%s(%s) %s", f.Name, p, r))
}

func (p *Parameter) String() string {
	return fmt.Sprintf("%s %s", p.Name, p.Type)
}

func (r *Result) String() string {
	return fmt.Sprintf("%s %s", r.Name, r.Type)
}
