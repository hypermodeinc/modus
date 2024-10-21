/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package wasmextractor

import (
	"fmt"
	"strings"
)

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

func NewFunction(name string) *Function {
	return &Function{Name: name}
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
