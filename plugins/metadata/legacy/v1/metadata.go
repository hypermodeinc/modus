/*
 * Copyright 2024 Hypermode, Inc.
 */

package v1

import (
	"time"

	"hypruntime/utils"

	"github.com/buger/jsonparser"
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

	optional, err := jsonparser.GetBoolean(data, "optional")
	if err != nil && err != jsonparser.KeyPathNotFoundError {
		return err
	}
	p.Optional = optional

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
