/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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
func MakeJsonObject(pairs []KeyValuePair, pretty bool) ([]byte, error) {
	var buf bytes.Buffer

	if pretty {
		buf.WriteString("{\n")
		for i, kv := range pairs {
			keyBytes, err := json.Marshal(kv.Key)
			if err != nil {
				return nil, err
			}
			valBytes, err := json.Marshal(kv.Value)
			if err != nil {
				return nil, err
			}

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
			keyBytes, err := json.Marshal(kv.Key)
			if err != nil {
				return nil, err
			}
			valBytes, err := json.Marshal(kv.Value)
			if err != nil {
				return nil, err
			}

			buf.Write(keyBytes)
			buf.WriteByte(':')
			buf.Write(valBytes)
			if i < len(pairs)-1 {
				buf.WriteByte(',')
			}
		}
		buf.WriteByte('}')
	}

	return buf.Bytes(), nil
}
