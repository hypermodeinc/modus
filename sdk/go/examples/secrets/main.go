/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/secrets"
)

func GetSecretValue(name string) string {
	value, err := secrets.GetSecretValue(name)
	if err != nil {
		return ""
	}
	return value
}
