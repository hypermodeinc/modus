/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package schemagen

import "strings"

// prefixes that are used to identify query fields, and will be trimmed from the field name
var queryTrimPrefixes = []string{"get", "list"}

// prefixes that are used to identify mutation fields
var mutationPrefixes = []string{
	"mutate",
	"post", "patch", "put", "delete",
	"add", "update", "insert", "upsert",
	"create", "edit", "save", "remove", "alter", "modify",
	"start", "stop",
}

func isMutation(fnName string) bool {
	return getPrefix(fnName, mutationPrefixes) != ""
}

func getFieldName(fnName string) string {
	prefix := getPrefix(fnName, queryTrimPrefixes)
	fieldName := strings.TrimPrefix(fnName, prefix)
	return strings.ToLower(fieldName[:1]) + fieldName[1:]
}

func getPrefix(fnName string, prefixes []string) string {
	for _, prefix := range prefixes {
		// check for exact match
		fnNameLowered := strings.ToLower(fnName)
		if fnNameLowered == prefix {
			return prefix
		}

		// check for a prefix, but only if the prefix is NOT followed by a lowercase letter
		// for example, we want to match "addPost" but not "additionalPosts"
		prefixLen := len(prefix)
		if len(fnName) > prefixLen && strings.HasPrefix(fnNameLowered, prefix) {
			c := fnName[prefixLen]
			if c < 'a' || c > 'z' {
				return prefix
			}
		}
	}

	return ""
}
