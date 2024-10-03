/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	HostTypePostgresql string = "postgresql"
)

type PostgresqlHostInfo struct {
	Name    string `json:"-"`
	Type    string `json:"type"`
	ConnStr string `json:"connString"`
}

func (p PostgresqlHostInfo) HostName() string {
	return p.Name
}

func (PostgresqlHostInfo) HostType() string {
	return HostTypePostgresql
}

func (h PostgresqlHostInfo) GetVariables() []string {
	return extractVariables(h.ConnStr)
}

func (h PostgresqlHostInfo) Hash() string {
	// Concatenate the attributes into a single string
	data := fmt.Sprintf("%v|%v|%v", h.Name, h.Type, h.ConnStr)

	// Compute the SHA-256 hash
	hash := sha256.Sum256([]byte(data))

	// Convert the hash to a hexadecimal string
	hashStr := hex.EncodeToString(hash[:])

	return hashStr
}
