/*
 * Copyright 2024 Hypermode, Inc.
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
