/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
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
	EndpointAuthNone   EndpointAuthType = "none"
	EndpointAuthBearer EndpointAuthType = "bearer"
)

type GraphqlEndpointInfo struct {
	Name string           `json:"-"`
	Path string           `json:"path"`
	Auth EndpointAuthType `json:"auth"`
}

func (e GraphqlEndpointInfo) EndpointName() string {
	return e.Name
}

func (e GraphqlEndpointInfo) EndpointType() EndpointType {
	return EndpointTypeGraphQL
}

func (e GraphqlEndpointInfo) EndpointAuth() EndpointAuthType {
	return e.Auth
}
