/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package secrets

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hypermodeAI/manifest"
)

type localSecretsProvider struct {
}

func (sp *localSecretsProvider) initialize(ctx context.Context) {
}

func (sp *localSecretsProvider) getHostSecrets(host manifest.HostInfo) (map[string]string, error) {
	prefix := "MODUS_" + strings.ToUpper(strings.ReplaceAll(host.HostName(), "-", "_")) + "_"
	secrets := make(map[string]string)
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, prefix) {
			pair := strings.SplitN(e, "=", 2)
			secrets[pair[0][len(prefix):]] = pair[1]
		}
	}

	return secrets, nil
}

func (sp *localSecretsProvider) getSecretValue(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", fmt.Errorf("environment variable %s was not found", name)
	}

	return v, nil
}

func (sp *localSecretsProvider) hasSecret(name string) bool {
	return os.Getenv(name) != ""
}
