/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import (
	"bytes"
	"encoding/json"
)

// JsonSerialize serializes the given value to JSON.
// Unlike json.Marshal, it does not escape HTML characters.
func JsonSerialize(v any) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	return buf.Bytes(), err
}

// JsonDeserialize deserializes the given JSON data into the given value.
// Unlike json.Unmarshal, it does not automatically use a float64 type for numbers.
func JsonDeserialize(data []byte, v any) error {
	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(v)
}
