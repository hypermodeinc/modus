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
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/tidwall/gjson"
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

	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		switch key.String() {
		case "name":
			p.Name = value.String()
		case "type":
			p.Type = value.String()
		case "default":
			val := value.Value()
			if val == nil {
				p.Default = new(any)
			} else {
				p.Default = &val
			}
		}
		return true
	})

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

func (m *Metadata) GetTypeDefinition(typ string) (*TypeDefinition, error) {
	switch typ {
	case "[]byte":
		return &TypeDefinition{typ, 1, nil}, nil
	case "string":
		return &TypeDefinition{typ, 2, nil}, nil
	}

	def, ok := m.Types[typ]
	if !ok {
		return nil, fmt.Errorf("info for type %s not found in plugin %s", typ, m.Name())
	}

	return def, nil
}

func (m *Metadata) GetExportedFunctions() []*Function {
	var fns []*Function
	for _, fn := range m.FnExports {
		fns = append(fns, fn)
	}
	return fns
}

func parseNameAndVersion(s string) (name string, version string) {
	i := strings.LastIndex(s, "@")
	if i == -1 {
		return s, ""
	}
	return s[:i], s[i+1:]
}

func GetMetadataFromContext(ctx context.Context) (*Metadata, error) {
	v := ctx.Value(utils.MetadataContextKey)
	if v == nil {
		return nil, errors.New("metadata not found in context")
	}
	return v.(*Metadata), nil
}
