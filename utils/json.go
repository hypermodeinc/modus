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
// It also returns results without an extra newline at the end.
func JsonSerialize(v any) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}

	// Remove the newline at the end, if there is one
	bytes := buf.Bytes()
	if bytes[len(bytes)-1] == '\n' {
		bytes = bytes[:len(bytes)-1]
	}

	return bytes, nil
}

// JsonDeserialize deserializes the given JSON data into the given value.
// Unlike json.Unmarshal, it does not automatically use a float64 type for numbers.
func JsonDeserialize(data []byte, v any) error {
	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(v)
}
