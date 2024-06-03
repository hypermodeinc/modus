/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import (
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"
)

var atomicCounter uint64

func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func ParseNameAndVersion(s string) (name string, version string) {
	i := strings.LastIndex(s, "@")
	if i == -1 {
		return s, ""
	}
	return s[:i], s[i+1:]
}

func EnvVarFlagEnabled(envVarName string) bool {
	v := os.Getenv(envVarName)
	b, err := strconv.ParseBool(v)
	return err == nil && b
}

func HypermodeDebugEnabled() bool {
	return EnvVarFlagEnabled("HYPERMODE_DEBUG")
}

func HypermodeTraceEnabled() bool {
	return EnvVarFlagEnabled("HYPERMODE_TRACE")
}

func TrimStringBefore(s string, sep string) string {
	parts := strings.SplitN(s, sep, 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return s
}

func GenerateUUIDV7() string {
	return uuid.Must(uuid.NewV7()).String()
}

func NextUint64() uint64 {
	return atomic.AddUint64(&atomicCounter, 1)
}
