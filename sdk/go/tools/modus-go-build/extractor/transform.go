/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package extractor

import (
	"go/types"

	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/metadata"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/utils"
	"golang.org/x/tools/go/packages"
)

func transformStruct(name string, s *types.Struct) *metadata.TypeDefinition {
	if s == nil {
		return nil
	}

	fields := make([]*metadata.Field, s.NumFields())

	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)
		fields[i] = &metadata.Field{
			Name: utils.CamelCase(f.Name()),
			Type: f.Type().String(),
		}
	}

	return &metadata.TypeDefinition{
		Name:   name,
		Fields: fields,
	}
}

func transformFunc(name string, f *types.Func, pkgs map[string]*packages.Package) *metadata.Function {
	if f == nil {
		return nil
	}

	sig := f.Type().(*types.Signature)
	params := sig.Params()
	results := sig.Results()

	ret := metadata.Function{
		Name: name,
		Docs: getFuncDocumentation(pkgs, f),
	}

	if params != nil {
		ret.Parameters = make([]*metadata.Parameter, params.Len())
		for i := 0; i < params.Len(); i++ {
			p := params.At(i)
			ret.Parameters[i] = &metadata.Parameter{
				Name: p.Name(),
				Type: p.Type().String(),
			}
		}
	}

	if results != nil {
		ret.Results = make([]*metadata.Result, results.Len())
		for i := 0; i < results.Len(); i++ {
			r := results.At(i)
			ret.Results[i] = &metadata.Result{
				Name: r.Name(),
				Type: r.Type().String(),
			}
		}
	}

	return &ret
}
