/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

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
	Uid        string   `json:"uid,omitempty"`
	Id         string   `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Version    string   `json:"version,omitempty"`
	Language   string   `json:"language,omitempty"`
	SdkVersion string   `json:"sdk_version,omitempty"`
	BuildId    string   `json:"build_id,omitempty"`
	BuildTime  string   `json:"build_time,omitempty"`
	GitRepo    string   `json:"git_repo,omitempty"`
	GitCommit  string   `json:"git_commit,omitempty"`
	DType      []string `json:"dgraph.type,omitempty"`
}

type PluginData struct {
	Plugins []*Plugin `json:"plugins"`
}
