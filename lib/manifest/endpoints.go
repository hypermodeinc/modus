/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

type EndpointInfo interface {
	EndpointName() string
	EndpointType() EndpointType
	EndpointAuth() EndpointAuthType
}

type EndpointType string

const (
	EndpointTypeGraphQL EndpointType = "graphql"
)

type EndpointAuthType string

const (
	EndpointAuthNone        EndpointAuthType = "none"
	EndpointAuthBearerToken EndpointAuthType = "bearer-token"
)

type GraphqlEndpointInfo struct {
	Name string           `json:"-"`
	Type EndpointType     `json:"type"`
	Path string           `json:"path"`
	Auth EndpointAuthType `json:"auth"`
}

func (e GraphqlEndpointInfo) EndpointName() string {
	return e.Name
}

func (e GraphqlEndpointInfo) EndpointType() EndpointType {
	return e.Type
}

func (e GraphqlEndpointInfo) EndpointAuth() EndpointAuthType {
	return e.Auth
}
