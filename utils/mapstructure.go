/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import "github.com/go-viper/mapstructure/v2"

func MapToStruct(m map[string]any, result any) error {

	config := &mapstructure.DecoderConfig{
		Result: result,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(m)
}
