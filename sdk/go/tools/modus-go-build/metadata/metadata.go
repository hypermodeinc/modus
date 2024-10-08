/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package metadata

import (
	"fmt"
	"go/token"
	"sort"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/utils"
	"github.com/rs/xid"
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
	Id     uint32   `json:"id"`
	Name   string   `json:"-"`
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

func NewMetadata() *Metadata {
	return &Metadata{
		BuildId:   xid.New().String(),
		BuildTime: utils.TimeNow(),
		FnExports: make(FunctionMap),
		FnImports: make(FunctionMap),
		Types:     make(TypeMap),
	}
}

func (f *Function) String(m *Metadata) string {

	imports := m.GetImports()

	b := strings.Builder{}
	b.WriteString(f.Name)
	b.WriteString("(")
	for i, p := range f.Parameters {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(p.Name)
		b.WriteString(" ")
		b.WriteString(utils.GetNameForType(p.Type, imports))
	}
	b.WriteString(")")

	switch len(f.Results) {
	case 0:
		return b.String()
	case 1:
		r := f.Results[0]
		if r.Name == "" {
			b.WriteString(" ")
			b.WriteString(utils.GetNameForType(r.Type, imports))
			return b.String()
		}
	}

	b.WriteString(" (")
	for i, r := range f.Results {
		if i > 0 {
			b.WriteString(", ")
		}
		if r.Name != "" {
			b.WriteString(r.Name)
			b.WriteString(" ")
		}

		b.WriteString(utils.GetNameForType(r.Type, imports))
	}
	b.WriteString(")")

	return b.String()
}

func (t *TypeDefinition) String(m *Metadata) string {

	imports := m.GetImports()

	if len(t.Fields) == 0 {
		return utils.GetNameForType(t.Name, imports)
	}

	b := strings.Builder{}
	b.WriteString(utils.GetNameForType(t.Name, imports))
	b.WriteString(" { ")
	for i, f := range t.Fields {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(f.Name)
		b.WriteString(" ")
		b.WriteString(utils.GetNameForType(f.Type, imports))
	}
	b.WriteString(" }")

	return b.String()
}

func (m *FunctionMap) SortedKeys() []string {
	keys := utils.MapKeys(*m)
	sort.Strings(keys)
	return keys
}

func (m *TypeMap) SortedKeys(mod string) []string {
	imports := make(map[string]string)
	imports[mod] = ""
	keys := utils.MapKeys(*m)
	sort.Slice(keys, func(i, j int) bool {
		a := strings.ToLower(utils.GetNameForType(keys[i], imports))
		b := strings.ToLower(utils.GetNameForType(keys[j], imports))
		return a < b
	})
	return keys
}

func (m *Metadata) GetImports() map[string]string {
	names := make(map[string]string)
	imports := make(map[string]string)
	imports[m.Module] = ""

	for _, t := range m.Types {
		for _, p := range utils.GetPackageNamesForType(t.Name) {
			if p == m.Module || p == "unsafe" {
				continue
			}
			if _, ok := imports[p]; !ok {
				name := p[strings.LastIndex(p, "/")+1:]

				// make sure the name is not a reserved word
				if token.IsKeyword(name) {
					name = fmt.Sprintf("%s_", name)
				}

				// make sure the name is unique
				n := name
				for i := 1; ; i++ {
					if _, ok := names[n]; !ok {
						break
					}
					n = fmt.Sprintf("%s%d", name, i)
				}

				imports[p] = n
				names[n] = p
			}
		}
	}

	return imports
}
