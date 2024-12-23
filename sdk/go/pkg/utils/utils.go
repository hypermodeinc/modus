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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func GetLocalTime(tz string) time.Time {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Now()
	}
	return time.Now().In(loc)
}

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

	// Remove the extraneous trailing newline
	bytes := buf.Bytes()[:buf.Len()-1]

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

func ConvertInterfaceTo[T any](v any) (T, error) {
	switch t := v.(type) {
	case T:
		return t, nil
	case json.Number:
		n, err := ParseJsonNumberAs[T](t)
		if err != nil {
			var zero T
			return zero, fmt.Errorf("could not convert JSON number to %T: %v", zero, err)
		}
		return n, nil
	default:
		var zero T
		return zero, fmt.Errorf("could not convert %T to %T", v, zero)
	}
}

func ParseJsonNumberAs[T any](n json.Number) (T, error) {
	var zero T

	switch any(zero).(type) {
	case int:
		x, err := strconv.ParseInt(string(n), 10, 32) // 32-bit wasm only
		if err != nil {
			return zero, err
		}
		return any(int(x)).(T), nil

	case int8:
		x, err := strconv.ParseInt(string(n), 10, 8)
		if err != nil {
			return zero, err
		}
		return any(int8(x)).(T), nil

	case int16:
		x, err := strconv.ParseInt(string(n), 10, 16)
		if err != nil {
			return zero, err
		}
		return any(int16(x)).(T), nil

	case int32:
		x, err := strconv.ParseInt(string(n), 10, 32)
		if err != nil {
			return zero, err
		}
		return any(int32(x)).(T), nil

	case int64:
		x, err := strconv.ParseInt(string(n), 10, 64)
		if err != nil {
			return zero, err
		}
		return any(int64(x)).(T), nil

	case uint:
		x, err := strconv.ParseUint(string(n), 10, 32) // 32-bit wasm only
		if err != nil {
			return zero, err
		}
		return any(uint(x)).(T), nil

	case uint8:
		x, err := strconv.ParseUint(string(n), 10, 8)
		if err != nil {
			return zero, err
		}
		return any(uint8(x)).(T), nil

	case uint16:
		x, err := strconv.ParseUint(string(n), 10, 16)
		if err != nil {
			return zero, err
		}
		return any(uint16(x)).(T), nil

	case uint32:
		x, err := strconv.ParseUint(string(n), 10, 32)
		if err != nil {
			return zero, err
		}
		return any(uint32(x)).(T), nil

	case uint64:
		x, err := strconv.ParseUint(string(n), 10, 64)
		if err != nil {
			return zero, err
		}
		return any(uint64(x)).(T), nil

	case float32:
		x, err := strconv.ParseFloat(string(n), 32)
		if err != nil {
			return zero, err
		}
		return any(float32(x)).(T), nil

	case float64:
		x, err := strconv.ParseFloat(string(n), 64)
		if err != nil {
			return zero, err
		}
		return any(float64(x)).(T), nil

	default:
		return zero, errors.New("unsupported type")
	}
}

type RawJsonString string

func (s RawJsonString) MarshalJSON() ([]byte, error) {
	if len(s) == 0 {
		return []byte(`"null"`), nil
	}
	return []byte(s), nil
}

func (s *RawJsonString) UnmarshalJSON(data []byte) error {
	if data == nil {
		return errors.New("unexpected nil JSON input")
	}
	*s = RawJsonString(data)
	return nil
}

type StringOrRawJson string

func (s StringOrRawJson) MarshalJSON() ([]byte, error) {
	// empty string
	if len(s) == 0 {
		return []byte(`""`), nil
	}

	// raw object, array, or string
	switch s[0] {
	case '{', '[', '"':
		return []byte(string(s)), nil
	}

	// raw literal
	switch s {
	case "null", "true", "false":
		return []byte(string(s)), nil
	}

	// raw number
	if _, err := strconv.ParseFloat(string(s), 64); err == nil {
		return []byte(string(s)), nil
	}

	// string that needs to be encoded
	return json.Marshal(string(s))
}

func (s *StringOrRawJson) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("unexpected empty JSON input")
	}

	// string that needs to be decoded
	if data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		*s = StringOrRawJson(str)
		return nil
	}

	// raw json
	*s = StringOrRawJson(data)
	return nil
}
