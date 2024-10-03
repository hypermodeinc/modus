/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package gitinfo

import (
	"os/exec"
	"strings"

	"github.com/hypermodeAI/functions-go/tools/hypbuild/config"
	"github.com/hypermodeAI/functions-go/tools/hypbuild/metadata"
)

func TryCollectGitInfo(config *config.Config, meta *metadata.Metadata) {
	if !isGitRepo(config.SourceDir) {
		return
	}

	meta.GitCommit = getGitCommit(config.SourceDir)
	meta.GitRepo = getGitRepo(config.SourceDir)
}

func isGitRepo(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = dir
	if result, err := cmd.Output(); err == nil {
		return strings.TrimSpace(string(result)) == "true"
	}
	return false
}

func getGitCommit(dir string) string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	if result, err := cmd.Output(); err == nil {
		return strings.TrimSpace(string(result))
	}
	return ""
}

func getGitRepo(dir string) string {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = dir
	if result, err := cmd.Output(); err == nil {
		url := strings.TrimSpace(string(result))

		// Convert ssh url to https
		if strings.HasPrefix(url, "git@") {
			url = strings.ReplaceAll(url, ":", "/")
			url = strings.Replace(url, "git@", "https://", 1)
		}

		// Remove .git suffix
		url, _ = strings.CutSuffix(url, ".git")

		return url
	}
	return ""
}
