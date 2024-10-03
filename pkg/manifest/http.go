/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
)

const (
	HostTypeHTTP string = "http"
)

var (
	templateRegex = regexp.MustCompile(`{{\s*(?:base64\((.+?):(.+?)\)|(.+?))\s*}}`)
)

type HTTPHostInfo struct {
	Name            string            `json:"-"`
	Type            string            `json:"type"`
	Endpoint        string            `json:"endpoint"`
	BaseURL         string            `json:"baseURL"`
	Headers         map[string]string `json:"headers"`
	QueryParameters map[string]string `json:"queryParameters"`
}

func (h HTTPHostInfo) HostName() string {
	return h.Name
}

func (h HTTPHostInfo) HostType() string {
	return HostTypeHTTP
}

func (h HTTPHostInfo) GetVariables() []string {
	cap := 2 * (len(h.Headers) + len(h.QueryParameters))
	set := make(map[string]bool, cap)
	results := make([]string, 0, cap)

	for _, header := range h.Headers {
		vars := extractVariables(header)
		for _, v := range vars {
			if _, ok := set[v]; !ok {
				set[v] = true
				results = append(results, v)
			}
		}
	}

	for _, v := range h.QueryParameters {
		vars := extractVariables(v)
		for _, v := range vars {
			if _, ok := set[v]; !ok {
				set[v] = true
				results = append(results, v)
			}
		}
	}

	return results
}

func (h HTTPHostInfo) Hash() string {
	// Concatenate the attributes into a single string
	// NOTE: The Type field is not included in the HTTP host hash, for backwards compatibility.  It should be include for other host types.
	data := fmt.Sprintf("%v|%v|%v|%v|%v", h.Name, h.Endpoint, h.BaseURL, h.Headers, h.QueryParameters)

	// Compute the SHA-256 hash
	hash := sha256.Sum256([]byte(data))

	// Convert the hash to a hexadecimal string
	hashStr := hex.EncodeToString(hash[:])

	return hashStr
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
