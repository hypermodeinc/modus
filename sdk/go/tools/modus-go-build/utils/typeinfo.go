/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"sort"
	"strings"
)

func GetNameForType(t string, imports map[string]string) string {
	sep := strings.LastIndex(t, ".")
	if sep == -1 {
		return t
	}

	if IsPointerType(t) {
		return "*" + GetNameForType(GetUnderlyingType(t), imports)
	}

	if IsListType(t) {
		return arrayPrefix(t) + GetNameForType(GetArraySubtype(t), imports)
	}

	if IsMapType(t) {
		kt, vt := GetMapSubtypes(t)
		return "map[" + GetNameForType(kt, imports) + "]" + GetNameForType(vt, imports)
	}

	pkgPath := t[:sep]
	pkgName := imports[pkgPath]
	typeName := t[sep+1:]
	if pkgName == "" {
		return typeName
	} else {
		return pkgName + "." + typeName
	}
}

func GetPackageNamesForType(t string) []string {
	if IsPointerType(t) {
		return GetPackageNamesForType(GetUnderlyingType(t))
	}

	if IsListType(t) {
		return GetPackageNamesForType(GetArraySubtype(t))
	}

	if IsMapType(t) {
		kt, vt := GetMapSubtypes(t)
		kp, vp := GetPackageNamesForType(kt), GetPackageNamesForType(vt)
		m := make(map[string]bool, len(kp)+len(vp))
		for _, p := range kp {
			m[p] = true
		}
		for _, p := range vp {
			m[p] = true
		}
		pkgs := MapKeys(m)
		sort.Strings(pkgs)
		return pkgs
	}

	if i := strings.LastIndex(t, "."); i != -1 {
		return []string{t[:i]}
	}

	return nil
}

func GetArraySubtype(t string) string {
	return t[strings.Index(t, "]")+1:]
}

func GetMapSubtypes(t string) (string, string) {
	const prefix = "map["
	if !strings.HasPrefix(t, prefix) {
		return "", ""
	}

	n := 1
	for i := len(prefix); i < len(t); i++ {
		switch t[i] {
		case '[':
			n++
		case ']':
			n--
			if n == 0 {
				return t[len(prefix):i], t[i+1:]
			}
		}
	}

	return "", ""
}

func GetUnderlyingType(t string) string {
	return strings.TrimPrefix(t, "*")
}

func IsListType(t string) bool {
	// covers both slices and arrays
	return strings.HasPrefix(t, "[")
}

func IsSliceType(t string) bool {
	return strings.HasPrefix(t, "[]")
}

func IsArrayType(t string) bool {
	return IsListType(t) && !IsSliceType(t)
}

func arrayPrefix(t string) string {
	return t[:strings.Index(t, "]")+1]
}

func IsMapType(t string) bool {
	return strings.HasPrefix(t, "map[")
}

func IsPointerType(t string) bool {
	return strings.HasPrefix(t, "*")
}

func IsStringType(t string) bool {
	return t == "string"
}

func IsStructType(t string) bool {
	return !IsPointerType(t) && !IsPrimitiveType(t) && !IsListType(t) && !IsMapType(t) && !IsStringType(t)
}

func IsPrimitiveType(t string) bool {
	switch t {
	case "bool",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"uintptr", "byte", "rune",
		"float32", "float64",
		"complex64", "complex128",
		"unsafe.Pointer", "time.Duration":
		return true
	}

	return false
}
