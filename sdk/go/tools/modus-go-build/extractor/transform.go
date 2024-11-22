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

	structDecl, structType := getStructDeclarationAndType(name, pkgs)
	if structDecl == nil || structType == nil {
		return nil
	}

	structDocs := getDocs(structDecl.Doc)

	fields := make([]*metadata.Field, s.NumFields())
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)

		fieldDocs := getDocs(structType.Fields.List[i].Doc)

		fields[i] = &metadata.Field{
			Name: utils.CamelCase(f.Name()),
			Type: f.Type().String(),
			Docs: fieldDocs,
		}
	}

	return &metadata.TypeDefinition{
		Name:   name,
		Fields: fields,
		Docs:   structDocs,
	}
}

func transformFunc(name string, f *types.Func, pkgs map[string]*packages.Package) *metadata.Function {
	if f == nil {
		return nil
	}

	sig := f.Type().(*types.Signature)
	params := sig.Params()
	results := sig.Results()

	funcDecl := getFuncDeclaration(f, pkgs)
	if funcDecl == nil {
		return nil
	}

	ret := metadata.Function{
		Name: name,
		Docs: getDocs(funcDecl.Doc),
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

func getStructDeclarationAndType(name string, pkgs map[string]*packages.Package) (*ast.GenDecl, *ast.StructType) {
	objName := name[strings.LastIndex(name, ".")+1:]
	pkgName := utils.GetPackageNamesForType(name)[0]
	pkg := pkgs[pkgName]

	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if typeSpec.Name.Name == objName {
							if structType, ok := typeSpec.Type.(*ast.StructType); ok {
								return genDecl, structType
							}
						}
					}
				}
			}
		}
	}

	return nil, nil
}

func getDocs(comments *ast.CommentGroup) *metadata.Docs {
	if comments == nil {
		return nil
	}

	var lines []string
	for _, comment := range comments.List {
		txt := comment.Text
		if strings.HasPrefix(txt, "// ") {
			txt = strings.TrimPrefix(txt, "// ")
			txt = strings.TrimSpace(txt)
			lines = append(lines, txt)
		} else if strings.HasPrefix(txt, "/*") {
			txt = strings.TrimPrefix(txt, "/*")
			txt = strings.TrimSuffix(txt, "*/")
			txt = strings.TrimSpace(txt)
			lines = append(lines, strings.Split(txt, "\n")...)
		}
	}

	if len(lines) == 0 {
		return nil
	}

	return &metadata.Docs{
		Lines: lines,
	}
}
