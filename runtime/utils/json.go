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
	"bytes"

	"github.com/goccy/go-json"
)

// JsonSerialize serializes the given value to JSON.
// Unlike json.Marshal, it does not escape HTML characters.
// It uses goccy/go-json to improve performance.
func JsonSerialize(v any) ([]byte, error) {
	return json.MarshalWithOption(v, json.DisableHTMLEscape())
}

// JsonDeserialize deserializes the given JSON data into the given value.
// Unlike json.Unmarshal, it does not automatically use a float64 type for numbers.
// It uses goccy/go-json to improve performance.
func JsonDeserialize(data []byte, v any) error {
	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(v)
}

type KeyValuePair struct {
	Key   string
	Value any
}

// MakeJsonObject creates a JSON object from the given key-value pairs.
func MakeJsonObject(pairs []KeyValuePair, pretty bool) []byte {
	var buf bytes.Buffer

	if pretty {
		buf.WriteString("{\n")
		for i, kv := range pairs {
			keyBytes, _ := json.Marshal(kv.Key)
			valBytes, _ := json.Marshal(kv.Value)

			buf.WriteString("  ")
			buf.Write(keyBytes)
			buf.WriteString(": ")
			buf.Write(valBytes)
			if i < len(pairs)-1 {
				buf.WriteByte(',')
			}
			buf.WriteByte('\n')
		}
		buf.WriteString("}\n")
	} else {
		buf.WriteByte('{')
		for i, kv := range pairs {
			keyBytes, _ := json.Marshal(kv.Key)
			valBytes, _ := json.Marshal(kv.Value)

			buf.Write(keyBytes)
			buf.WriteByte(':')
			buf.Write(valBytes)
			if i < len(pairs)-1 {
				buf.WriteByte(',')
			}
		}
		buf.WriteByte('}')
	}

	return buf.Bytes()
}
