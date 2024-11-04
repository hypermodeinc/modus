/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package postgresql

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/utils"
)

type UUID [16]uint8

func (u *UUID) String() string {
	return hex.EncodeToString(u[:])
}

func (u *UUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.String() + `"`), nil
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	if len(data) < 34 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid UUID: %s", string(data))
	}

	_, err := hex.Decode(u[:], data[1:len(data)-1])
	return err
}

func InterfaceSliceToUUID(data []interface{}) (*UUID, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf("invalid UUID length: %d", len(data))
	}

	var uuid UUID
	for i, v := range data {
		n, ok := v.(json.Number)
		if !ok {
			return nil, fmt.Errorf("invalid UUID byte: %v", v)
		}
		u, err := utils.ParseJsonNumberAs[uint8](n)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID byte: %v", v)
		}
		uuid[i] = u
	}

	return &uuid, nil
}
