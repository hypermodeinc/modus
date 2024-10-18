/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func computeHash(elements ...any) string {
	data := &bytes.Buffer{}
	for i, e := range elements {
		if i > 0 {
			data.WriteByte('|')
		}
		if s, ok := e.(string); ok {
			data.WriteString(s)
		} else {
			fmt.Fprint(data, e)
		}
	}

	hash := sha256.Sum256(data.Bytes())
	return hex.EncodeToString(hash[:])
}

func extractVariables(s string) []string {
	matches := templateRegex.FindAllStringSubmatch(s, -1)
	if matches == nil {
		return []string{}
	}

	results := make([]string, 0, len(matches)*2)
	for _, match := range matches {
		for j := 1; j < len(match); j++ {
			if match[j] != "" {
				results = append(results, match[j])
			}
		}
	}

	return results
}
