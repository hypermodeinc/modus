/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package modusdb

import (
	"time"
)

const internalSchema = `
id: string @index(exact) .
name: string @index(term) .
version: string @index(term) .
language: string @index(term) .
sdk_version: string @index(term) .
build_id: string @index(term) .
build_time: string @index(term) .
git_repo: string @index(term) .
git_commit: string @index(term) .

model_hash: string @index(exact) .
input: string @index(exact) .
output: string @index(exact) .
started_at: dateTime @index(day) .
duration_ms: int @index(int) .
plugin: uid @reverse .
function: string @index(exact) .

type Plugin {
	id
	name
	version
	language
	sdk_version
	build_id
	build_time
	git_repo
	git_commit
}

type Inference {
	id
	model_hash
	input
	output
	started_at
	duration_ms
	plugin
	function
}
`

type Plugin struct {
	Uid        string      `json:"uid,omitempty"`
	Id         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Version    string      `json:"version,omitempty"`
	Language   string      `json:"language,omitempty"`
	SdkVersion string      `json:"sdk_version,omitempty"`
	BuildId    string      `json:"build_id,omitempty"`
	BuildTime  string      `json:"build_time,omitempty"`
	GitRepo    string      `json:"git_repo,omitempty"`
	GitCommit  string      `json:"git_commit,omitempty"`
	Inferences []Inference `json:"inferences,omitempty"`
	DType      []string    `json:"dgraph.type,omitempty"`
}

type PluginData struct {
	Plugins []Plugin `json:"plugins"`
}

type Inference struct {
	Uid        string    `json:"uid,omitempty"`
	Id         string    `json:"id,omitempty"`
	ModelHash  string    `json:"model_hash,omitempty"`
	Input      string    `json:"input,omitempty"`
	Output     string    `json:"output,omitempty"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	DurationMs int64     `json:"duration_ms,omitempty"`
	Plugin     *Plugin   `json:"plugin,omitempty"`
	Function   string    `json:"function,omitempty"`
	DType      []string  `json:"dgraph.type,omitempty"`
}

type InferenceData struct {
	Inferences []Inference `json:"inferences"`
}
