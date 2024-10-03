/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type ModelInfo struct {
	Name        string `json:"-"`
	SourceModel string `json:"sourceModel"`
	Provider    string `json:"provider"`
	Host        string `json:"host"`
	Path        string `json:"path"`
	Dedicated   bool   `json:"dedicated"`
}

func (m ModelInfo) Hash() string {
	// Concatenate the attributes into a single string
	data := fmt.Sprintf("%v|%v|%v|%v", m.Name, m.SourceModel, m.Provider, m.Host)
	// Don't include the "dedicated" attribute if host is NOT "hypermode" or
	// if it's NOT "dedicated" (default)
	if m.Host == "hypermode" && m.Dedicated {
		data += fmt.Sprintf("|%v", m.Dedicated)
	}

	// Compute the SHA-256 hash
	hash := sha256.Sum256([]byte(data))

	// Convert the hash to a hexadecimal string
	hashStr := hex.EncodeToString(hash[:])

	return hashStr
}
