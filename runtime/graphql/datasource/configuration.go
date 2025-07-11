/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package datasource

import (
	"github.com/hypermodeinc/modus/runtime/wasmhost"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/plan"
)

type ModusDataSourceConfig struct {
	WasmHost          wasmhost.WasmHost
	FieldsToFunctions map[string]string
	MapTypes          []string
	Metadata          *plan.DataSourceMetadata
}
