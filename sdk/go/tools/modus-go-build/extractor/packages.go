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
	"os"

	"golang.org/x/tools/go/packages"
)

func loadPackages(dir string) (map[string]*packages.Package, error) {
	mode := packages.NeedName |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo

	cfg := &packages.Config{
		Mode: mode,
		Dir:  dir,
		Env:  append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm"),
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, err
	}

	pkgMap := make(map[string]*packages.Package)
	for _, pkg := range pkgs {
		expandPackages(pkg, pkgMap)
	}

	return pkgMap, nil
}

func expandPackages(pkg *packages.Package, pkgMap map[string]*packages.Package) {
	for _, imp := range pkg.Imports {
		if _, ok := pkgMap[imp.PkgPath]; ok {
			continue
		}
		expandPackages(imp, pkgMap)
	}
	pkgMap[pkg.PkgPath] = pkg
}
