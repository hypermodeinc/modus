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
	"path/filepath"
	"sort"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/config"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/metadata"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/utils"
)

func PostProcess(config *config.Config, meta *metadata.Metadata) error {
	imports := meta.GetImports()
	types := getTypes(meta)

	body := &bytes.Buffer{}
	writeFuncUnpin(body)
	writeFuncNew(body, types, imports)
	writeFuncMake(body, types, imports)
	writeFuncReadMap(body, types, imports)
	writeFuncWriteMap(body, types, imports)

	header := &bytes.Buffer{}
	writePostProcessHeader(header, meta, imports)

	return writeBuffersToFile(filepath.Join(config.SourceDir, post_file), header, body)
}

func getTypes(meta *metadata.Metadata) []*metadata.TypeDefinition {
	types := utils.MapValues(meta.Types)
	sort.Slice(types, func(i, j int) bool {
		return types[i].Id < types[j].Id
	})
	return types
}

func writePostProcessHeader(b *bytes.Buffer, meta *metadata.Metadata, imports map[string]string) {
	b.WriteString("// Code generated by Modus. DO NOT EDIT.\n\n")
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"unsafe\"\n")
	for pkg, name := range imports {
		if pkg != "" && pkg != "unsafe" && pkg != meta.Module {
			if pkg == name || strings.HasSuffix(pkg, "/"+name) {
				fmt.Fprintf(b, "\t\"%s\"\n", pkg)
			} else {
				fmt.Fprintf(b, "\t%s \"%s\"\n", name, pkg)
			}
		}
	}
	b.WriteString(")\n")
}

func writeFuncUnpin(b *bytes.Buffer) {
	b.WriteString(`
var __pins = make(map[unsafe.Pointer]int)

//go:export __unpin
func __unpin(p unsafe.Pointer) {
	n := __pins[p]
	if n == 1 {
		delete(__pins, p)
	} else {
		__pins[p] = n - 1
	}
}`)

	b.WriteString("\n")
}

func writeFuncNew(b *bytes.Buffer, types []*metadata.TypeDefinition, imports map[string]string) {
	buf := &bytes.Buffer{}
	found := false

	buf.WriteString(`
//go:export __new
func __new(id int) unsafe.Pointer {
`)
	buf.WriteString("\tswitch id {\n")
	for _, t := range types {
		if utils.IsSliceType(t.Name) && utils.IsMapType(t.Name) {
			continue
		}

		ptrName := utils.GetNameForType(t.Name, imports)
		elementName := utils.GetUnderlyingType(ptrName)
		found = true
		fmt.Fprintf(buf, `	case %d:
		o := new(%s)
		p := unsafe.Pointer(o)
		__pins[p]++
		return p
`, t.Id, elementName)
	}
	buf.WriteString("\t}\n\n")
	buf.WriteString("\treturn nil\n}\n")

	if found {
		_, _ = buf.WriteTo(b)
	}
}

func writeFuncMake(b *bytes.Buffer, types []*metadata.TypeDefinition, imports map[string]string) {

	b.WriteString(`
//go:export __make
func __make(id, size int) unsafe.Pointer {
	switch id {
	case 1:
		o := make([]byte, size)
		p := unsafe.Pointer(&o)
		__pins[p]++
		return p
	case 2:
		o := string(make([]byte, size))
		p := unsafe.Pointer(&o)
		__pins[p]++
		return p
`)

	for _, t := range types {
		name := utils.GetNameForType(t.Name, imports)
		if utils.IsSliceType(name) || utils.IsMapType(name) {
			fmt.Fprintf(b, `	case %d:
		o := make(%s, size)
		p := unsafe.Pointer(&o)
		__pins[p]++
		return p
`, t.Id, name)
		}
	}

	b.WriteString("\t}\n\n")
	b.WriteString("\treturn nil\n}\n")
}

func writeFuncReadMap(b *bytes.Buffer, types []*metadata.TypeDefinition, imports map[string]string) {
	buf := &bytes.Buffer{}
	found := false

	buf.WriteString(`
//go:export __read_map
func __read_map(id int, m unsafe.Pointer) uint64 {
`)
	buf.WriteString("\tswitch id {\n")
	for _, t := range types {
		if utils.IsMapType(t.Name) {
			found = true
			typeName := utils.GetNameForType(t.Name, imports)
			fmt.Fprintf(buf, `	case %d:
		return __doReadMap(*(*%s)(m))
`, t.Id, typeName)
		}
	}
	buf.WriteString("\t}\n\n")
	buf.WriteString("\treturn 0\n}\n")

	buf.WriteString(`
func __doReadMap[M ~map[K]V, K comparable, V any](m M) uint64 {
	size := len(m)
	keys := make([]K, size)
	values := make([]V, size)
	
	i := 0
	for k, v := range m {
		keys[i] = k
		values[i] = v
		i++
	}

	pKeys := uint32(uintptr(unsafe.Pointer(&keys)))
	pValues := uint32(uintptr(unsafe.Pointer(&values)))
	return uint64(pKeys)<<32 | uint64(pValues)
}
`)

	if found {
		_, _ = buf.WriteTo(b)
	}
}

func writeFuncWriteMap(b *bytes.Buffer, types []*metadata.TypeDefinition, imports map[string]string) {
	buf := &bytes.Buffer{}
	found := false

	buf.WriteString(`
//go:export __write_map
func __write_map(id int, m, keys, values unsafe.Pointer) {
`)
	buf.WriteString("\tswitch id {\n")
	for _, t := range types {
		if strings.HasPrefix(t.Name, "map[") {
			found = true
			typeName := utils.GetNameForType(t.Name, imports)
			kt, vt := utils.GetMapSubtypes(t.Name)
			keyTypeName := utils.GetNameForType(kt, imports)
			valTypeName := utils.GetNameForType(vt, imports)
			fmt.Fprintf(buf, `	case %d:
		__doWriteMap(*(*%s)(m), *(*[]%s)(keys), *(*[]%s)(values))
`, t.Id, typeName, keyTypeName, valTypeName)
		}
	}
	buf.WriteString("\t}\n")
	buf.WriteString("}\n")

	buf.WriteString(`
func __doWriteMap[M ~map[K]V, K comparable, V any](m M, keys[]K, values[]V) {
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = values[i]
	}
}
`)

	if found {
		_, _ = buf.WriteTo(b)
	}
}
