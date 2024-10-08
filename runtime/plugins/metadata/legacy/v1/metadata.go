/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package v1

import (
	"time"

	"github.com/tidwall/gjson"
)

type Metadata struct {
	Plugin    string            `json:"plugin"`
	SDK       string            `json:"sdk"`
	Library   string            `json:"library"` // deprecated
	BuildId   string            `json:"buildId"`
	BuildTime time.Time         `json:"buildTs"`
	GitRepo   string            `json:"gitRepo"`
	GitCommit string            `json:"gitCommit"`
	Functions []*Function       `json:"functions"`
	Types     []*TypeDefinition `json:"types"`
}

type Function struct {
	Name       string       `json:"name"`
	Parameters []*Parameter `json:"parameters"`
	ReturnType *TypeInfo    `json:"returnType"`
}

type TypeDefinition struct {
	Id     uint32   `json:"id"`
	Path   string   `json:"path"`
	Name   string   `json:"name"`
	Fields []*Field `json:"fields"`
}

type Parameter struct {
	Name     string    `json:"name"`
	Type     *TypeInfo `json:"type"`
	Optional bool      `json:"optional"` // deprecated
	Default  *any      `json:"default"`
}

type Field struct {
	Name string    `json:"name"`
	Type *TypeInfo `json:"type"`
}

type TypeInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (p *Parameter) UnmarshalJSON(data []byte) error {

	// We need to manually unmarshal the JSON to distinguish between a null default
	// value and the absence of a default value.

	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		switch key.String() {
		case "name":
			p.Name = value.String()
		case "type":
			p.Type = &TypeInfo{
				Name: value.Get("name").String(),
				Path: value.Get("path").String(),
			}
		case "optional":
			p.Optional = value.Bool()
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
