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
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/metadata"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/utils"
	"golang.org/x/tools/go/packages"
)

func transformStruct(name string, s *types.Struct, pkgs map[string]*packages.Package) *metadata.TypeDefinition {
	if s == nil {
		return nil
	}

	var structComments []string

	name = name[strings.Index(name, ".")+1:]
	fields := make([]*metadata.Field, s.NumFields())
	
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)

		var fieldComments []string

		for _, pkg := range pkgs {
			for _, file := range pkg.Syntax {
				for _, decl := range file.Decls {
					if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
						for _, spec := range genDecl.Specs {
							if typeSpec, ok := spec.(*ast.TypeSpec); ok {
								if typeSpec.Name.Name == name {
									if structType, ok := typeSpec.Type.(*ast.StructType); ok {
										if typeSpec.Doc != nil {
											for _, comment := range typeSpec.Doc.List {
												txt := comment.Text
												if strings.HasPrefix(txt, "// ") {
													structComments = append(structComments, strings.TrimPrefix(txt, "// "))
												}
											}
										}
										// fmt.Println("Type:", typeSpec.Name.Name)
										// fmt.Println("Comments:", typeSpec.Comment)
										// fmt.Println("Docs: ", typeSpec.Doc)
										for _, field := range structType.Fields.List {
											if field.Names != nil && field.Names[0].Name == f.Name() {
												if field.Doc != nil {
													for _, comment := range field.Doc.List {
														txt := comment.Text
														if strings.HasPrefix(txt, "// ") {
															fieldComments = append(fieldComments, strings.TrimPrefix(txt, "// "))
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if len(fieldComments) == 0 {
			fields[i] = &metadata.Field{
				Name: utils.CamelCase(f.Name()),
				Type: f.Type().String(),
			}
		} else {
			fields[i] = &metadata.Field{
				Name: utils.CamelCase(f.Name()),
				Type: f.Type().String(),
				Docs: &metadata.Docs{
					Description: strings.Join(fieldComments, "\n"),
				},
			}
		}
	}

	if len(structComments) == 0 {
		return &metadata.TypeDefinition{
			Name:   name,
			Fields: fields,
		}
	} else {
		return &metadata.TypeDefinition{
			Name:   name,
			Fields: fields,
			Docs: &metadata.Docs{
				Description: strings.Join(structComments, "\n"),
			},
		}
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
