/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import "time"

type Person struct {
	Uid       string   `json:"uid,omitempty"`
	FirstName string   `json:"firstName,omitempty"`
	LastName  string   `json:"lastName,omitempty"`
	DType     []string `json:"dgraph.type,omitempty"`
}

type PeopleData struct {
	People []*Person `json:"people"`
}

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
