/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/config"
	"golang.org/x/tools/go/packages"
)

func PreProcess(config *config.Config) error {

	// cleanup from previous runs first
	err := cleanup(config.SourceDir)
	if err != nil {
		return err
	}

	pkg, err := getMainPackage(config.SourceDir)
	if err != nil {
		return err
	}

	functions := getFunctionsNeedingWrappers(pkg)
	imports := getRequiredImports(functions)

	body := &bytes.Buffer{}
	writeMainFunc(body, pkg)
	if err := writeFuncWrappers(body, pkg, imports, functions); err != nil {
		return err
	}

	header := &bytes.Buffer{}
	writePreProcessHeader(header, imports)

	return writeBuffersToFile(filepath.Join(config.SourceDir, pre_file), header, body)
}

func getMainPackage(dir string) (*packages.Package, error) {
	mode := packages.NeedName |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo

	cfg := &packages.Config{Mode: mode, Dir: dir}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, err
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("expected exactly one root package, got %d", len(pkgs))
	}

	pkg := pkgs[0]

	if pkg.Name != "main" {
		return nil, fmt.Errorf("expected root package name to be 'main', got %s", pkg.Name)
	}

	return pkg, nil
}

func writeMainFunc(b *bytes.Buffer, pkg *packages.Package) {
	// see if there is a main function
	for _, f := range pkg.Syntax {
		for _, d := range f.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok {
				if fn.Name.Name == "main" {
					// found, nothing to do
					return
				}
			}
		}
	}

	// add a main function
	b.WriteString("func main() {}\n\n")
}

type funcInfo struct {
	function *ast.FuncDecl
	imports  map[string]string
	aliases  map[string]string
}

func getFunctionsNeedingWrappers(pkg *packages.Package) []*funcInfo {
	var results []*funcInfo
	for _, f := range pkg.Syntax {

		imports := make(map[string]string, len(f.Imports))
		for _, imp := range f.Imports {
			pkgPath := strings.Trim(imp.Path.Value, `"`)
			var name string
			if imp.Name != nil {
				name = imp.Name.Name
			} else {
				name = pkgPath[strings.LastIndex(pkgPath, "/")+1:]
			}
			imports[name] = pkgPath
		}

		for _, d := range f.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok {
				if fn.Name.IsExported() && getExportedFuncName(fn) == "" {

					var fields []*ast.Field
					if fn.Type.Params != nil {
						fields = append(fields, fn.Type.Params.List...)
					}
					if fn.Type.Results != nil {
						fields = append(fields, fn.Type.Results.List...)
					}

					usedImports := make(map[string]string)
					for _, p := range fields {
						importNames := getImportNames(p.Type)
						for _, name := range importNames {
							usedImports[name] = imports[name]
						}
					}

					info := &funcInfo{
						function: fn,
						imports:  usedImports,
					}

					results = append(results, info)
				}
			}
		}
	}
	return results
}

func getImportNames(e ast.Expr) []string {
	switch t := e.(type) {
	case *ast.StarExpr:
		return getImportNames(t.X)
	case *ast.ArrayType:
		return getImportNames(t.Elt)
	case *ast.MapType:
		return append(getImportNames(t.Key), getImportNames(t.Value)...)
	case *ast.StructType:
		names := make([]string, 0, len(t.Fields.List))
		for _, f := range t.Fields.List {
			names = append(names, getImportNames(f.Type)...)
		}
		return names
	case *ast.SelectorExpr:
		if id, ok := t.X.(*ast.Ident); ok {
			return []string{id.Name}
		}
	}
	return nil
}

func getRequiredImports(fns []*funcInfo) map[string]string {
	imports := make(map[string]string)
	names := make(map[string]string)

	for _, info := range fns {
		info.aliases = make(map[string]string, len(info.imports))
		for alias, pkgPath := range info.imports {
			name := pkgPath[strings.LastIndex(pkgPath, "/")+1:]

			// make sure the name is not a reserved word
			if token.IsKeyword(name) {
				name = fmt.Sprintf("%s_", name)
			}

			// make sure each package's name is unique
			n := name
			for i := 1; ; i++ {
				if p, ok := names[n]; !ok {
					break
				} else if p == pkgPath {
					break
				}
				n = fmt.Sprintf("%s%d", name, i)
			}

			imports[pkgPath] = n
			names[n] = pkgPath

			if alias != n {
				info.aliases[alias] = n
			}
		}
	}

	return imports
}

func writePreProcessHeader(b *bytes.Buffer, imports map[string]string) {
	b.WriteString("// Code generated by Modus. DO NOT EDIT.\n\n")
	b.WriteString("package main\n\n")

	if len(imports) == 0 {
		return
	}

	b.WriteString("import (\n")
	for pkg, name := range imports {
		if pkg == name || strings.HasSuffix(pkg, "/"+name) {
			b.WriteString(fmt.Sprintf("\t\"%s\"\n", pkg))
		} else {
			b.WriteString(fmt.Sprintf("\t%s \"%s\"\n", name, pkg))
		}
	}
	b.WriteString(")\n\n")
}

func writeFuncWrappers(b *bytes.Buffer, pkg *packages.Package, imports map[string]string, fns []*funcInfo) error {
	for _, info := range fns {
		fn := info.function
		name := fn.Name.Name
		params := fn.Type.Params
		results := fn.Type.Results

		hasErrorReturn := false
		if results != nil {
			for i := len(results.List) - 1; i >= 0; i-- {
				t := results.List[i].Type
				if ident, ok := t.(*ast.Ident); ok && ident.Name == "error" {
					if i == len(results.List)-1 && len(results.List[i].Names) <= 1 {
						hasErrorReturn = true
					} else {
						return fmt.Errorf("only one error can be returned, and it must be the last return value, for function %s", name)
					}
				}
			}
		}

		if hasErrorReturn {
			imports["github.com/hypermodeinc/modus/sdk/go/pkg/console"] = "console"

			// remove the error return value from the function signature
			results.List = results.List[:len(results.List)-1]
		}

		if fn.Doc != nil {
			var docStrings []string
			for _, comment := range fn.Doc.List {
				txt := comment.Text
				if strings.HasPrefix(txt, "// ") {
					docStrings = append(docStrings, txt)
				}
			}
			b.WriteString(strings.Join(docStrings, "\n"))
			b.WriteString("\n")
		}
		b.WriteString("//go:export ")
		b.WriteString(strings.ToLower(name[:1]))
		b.WriteString(name[1:])
		b.WriteString("\n")
		b.WriteString("func __modus_")
		b.WriteString(name)

		buf := &bytes.Buffer{}
		printer.Fprint(buf, pkg.Fset, fn.Type)
		decl := strings.TrimPrefix(buf.String(), "func")
		for a, n := range info.aliases {
			re := regexp.MustCompile(`\b` + a + `\.`)
			decl = re.ReplaceAllString(decl, n+".")
		}

		b.WriteString(decl)
		b.WriteString(" {\n")

		inputParams := &bytes.Buffer{}
		inputParams.WriteByte('(')
		for i, p := range params.List {
			if i > 0 {
				inputParams.WriteString(", ")
			}
			for j, n := range p.Names {
				if j > 0 {
					inputParams.WriteString(", ")
				}
				inputParams.WriteString(n.Name)
				if _, ok := p.Type.(*ast.Ellipsis); ok {
					inputParams.WriteString("...")
				}
			}
		}
		inputParams.WriteByte(')')

		numResults := 0
		if results != nil {
			numResults = results.NumFields()
		}

		if hasErrorReturn {
			b.WriteString("\t")
			if results != nil {
				for i, r := range results.List {
					n := len(r.Names)
					if n == 0 {
						n = 1
					}
					for j := 0; j < n; j++ {
						b.WriteByte('r')
						b.WriteString(strconv.Itoa(i + j))
						b.WriteString(", ")
					}
				}
			}
			b.WriteString("err := ")
			b.WriteString(name)
			b.Write(inputParams.Bytes())
			b.WriteByte('\n')

			b.WriteString("\tif err != nil {\n")
			b.WriteString("\t\tconsole.Error(err.Error())\n")
			b.WriteString("\t}\n")

			if numResults > 0 {
				b.WriteString("\treturn ")
				for i := 0; i < numResults; i++ {
					if i > 0 {
						b.WriteString(", ")
					}
					b.WriteByte('r')
					b.WriteString(strconv.Itoa(i))
				}
				b.WriteByte('\n')
			}

		} else {
			b.WriteString("\t")
			if results != nil && results.NumFields() > 0 {
				b.WriteString("return ")
			}
			b.WriteString(name)
			b.Write(inputParams.Bytes())
			b.WriteByte('\n')
		}
		b.WriteString("}\n\n")
	}

	return nil
}

func getExportedFuncName(fn *ast.FuncDecl) string {
	if fn.Body != nil && fn.Doc != nil {
		for _, c := range fn.Doc.List {
			parts := strings.Split(c.Text, " ")
			if len(parts) == 2 {
				switch parts[0] {
				case "//export", "//go:export", "//go:wasmexport":
					return parts[1]
				}
			}
		}
	}
	return ""
}
