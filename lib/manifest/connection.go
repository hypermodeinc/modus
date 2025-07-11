/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

type ConnectionType string

type ConnectionInfo interface {
	ConnectionName() string
	ConnectionType() ConnectionType
	Hash() string
	Variables() []string
}
