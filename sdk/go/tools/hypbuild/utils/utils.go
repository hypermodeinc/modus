/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"
)

func IsDebugModeEnabled() bool {
	b, _ := strconv.ParseBool(os.Getenv("HYPERMODE_DEBUG"))
	return b
}

func JsonSerialize(v any, ident bool) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if ident {
		enc.SetIndent("", "  ")
	}

	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}

	// Remove the newline at the end
	bytes := buf.Bytes()
	bytes = bytes[:len(bytes)-1]

	return bytes, nil
}

func MapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func TimeNow() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

func CamelCase(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToLower(str[:1]) + str[1:]
}
