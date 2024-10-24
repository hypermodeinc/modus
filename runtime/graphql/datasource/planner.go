/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package datasource

import (
	"bytes"
	"context"
	"fmt"
	"slices"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/wundergraph/graphql-go-tools/v2/pkg/ast"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/plan"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

type HypDSPlanner struct {
	id        int
	ctx       context.Context
	config    HypDSConfig
	visitor   *plan.Visitor
	variables resolve.Variables
	fields    map[int]fieldInfo
	template  struct {
		function *fieldInfo
		data     []byte
	}
}

type fieldInfo struct {
	ref       int         `json:"-"`
	Name      string      `json:"name"`
	Alias     string      `json:"alias,omitempty"`
	TypeName  string      `json:"type,omitempty"`
	Fields    []fieldInfo `json:"fields,omitempty"`
	IsMapType bool        `json:"isMapType,omitempty"`
	fieldRefs []int       `json:"-"`
}

func (t *fieldInfo) AliasOrName() string {
	if t.Alias != "" {
		return t.Alias
	}
	return t.Name
}

func (p *HypDSPlanner) SetID(id int) {
	p.id = id
}

func (p *HypDSPlanner) ID() (id int) {
	return p.id
}

func (p *HypDSPlanner) UpstreamSchema(dataSourceConfig plan.DataSourceConfiguration[HypDSConfig]) (*ast.Document, bool) {
	return nil, false
}

func (p *HypDSPlanner) DownstreamResponseFieldAlias(downstreamFieldRef int) (alias string, exists bool) {
	return
}

func (p *HypDSPlanner) DataSourcePlanningBehavior() plan.DataSourcePlanningBehavior {
	return plan.DataSourcePlanningBehavior{
		// This needs to be true, so we can distinguish results for multiple function calls in the same query.
		// Example:
		// query SayHello {
		//     a: sayHello(name: "Sam")
		//     b: sayHello(name: "Bob")
		// }
		// In this case, the Load function will be called twice, once for "a" and once for "b",
		// and the alias will be used in the return value to distinguish the results.
		OverrideFieldPathFromAlias: true,

		// This ensures that the __typename field is visited so we can include it in the response when requested.
		IncludeTypeNameFields: true,
	}
}

func (p *HypDSPlanner) Register(visitor *plan.Visitor, configuration plan.DataSourceConfiguration[HypDSConfig], dspc plan.DataSourcePlannerConfiguration) error {
	p.visitor = visitor
	visitor.Walker.RegisterEnterDocumentVisitor(p)
	visitor.Walker.RegisterEnterFieldVisitor(p)
	visitor.Walker.RegisterLeaveDocumentVisitor(p)
	p.config = HypDSConfig(configuration.CustomConfiguration())
	return nil
}

func (p *HypDSPlanner) EnterDocument(operation, definition *ast.Document) {
	p.fields = make(map[int]fieldInfo, len(operation.Fields))
}

func (p *HypDSPlanner) EnterField(ref int) {

	// Capture information about every field in the query.
	f := p.captureField(ref)
	p.fields[ref] = *f

	// If the field is enclosed by a root node, then it represents the function we want to call.
	if p.enclosingTypeIsRootNode() {

		// Save the field for the function.
		p.template.function = f

		// Also capture the input data for the function.
		err := p.captureInputData(ref)
		if err != nil {
			logger.Err(p.ctx, err).Msg("Error capturing input data.")
			return
		}
	}
}

func (p *HypDSPlanner) LeaveDocument(operation, definition *ast.Document) {
	// Stitch the captured fields together to form a tree.
	p.stitchFields(p.template.function)
}

func (p *HypDSPlanner) stitchFields(f *fieldInfo) {
	if len(f.fieldRefs) == 0 {
		return
	}

	f.Fields = make([]fieldInfo, len(f.fieldRefs))
	for i, ref := range f.fieldRefs {
		field := p.fields[ref]
		p.stitchFields(&field)
		f.Fields[i] = field
	}
}

func (p *HypDSPlanner) enclosingTypeIsRootNode() bool {
	enclosingTypeDef := p.visitor.Walker.EnclosingTypeDefinition
	for _, node := range p.visitor.Operation.RootNodes {
		if node.Ref == enclosingTypeDef.Ref {
			return true
		}
	}
	return false
}

func (p *HypDSPlanner) captureField(ref int) *fieldInfo {
	operation := p.visitor.Operation
	definition := p.visitor.Definition
	walker := p.visitor.Walker

	f := &fieldInfo{
		ref:   ref,
		Name:  operation.FieldNameString(ref),
		Alias: operation.FieldAliasString(ref),
	}

	def, ok := walker.FieldDefinition(ref)
	if ok {
		f.TypeName = definition.FieldDefinitionTypeNameString(def)
		f.IsMapType = slices.Contains(p.config.MapTypes, f.TypeName)
	}

	if operation.FieldHasSelections(ref) {
		ssRef, ok := operation.FieldSelectionSet(ref)
		if ok {
			f.fieldRefs = operation.SelectionSetFieldSelections(ssRef)
		}
	}

	return f
}

func (p *HypDSPlanner) captureInputData(fieldRef int) error {
	operation := p.visitor.Operation
	variables := resolve.NewVariables()
	var buf bytes.Buffer
	buf.WriteByte('{')

	args := operation.FieldArguments(fieldRef)
	for i, arg := range args {
		if i > 0 {
			buf.WriteByte(',')
		}

		argValue := operation.ArgumentValue(arg)
		if argValue.Kind != ast.ValueKindVariable {
			continue
		}

		argName := operation.ArgumentNameString(arg)

		variableName := operation.VariableValueNameString(argValue.Ref)
		placeHolder, _ := variables.AddVariable(
			&resolve.ContextVariable{
				Path:     []string{variableName},
				Renderer: resolve.NewJSONVariableRenderer(),
			})

		escapedKey, err := utils.JsonSerialize(argName)
		if err != nil {
			return err
		}

		buf.Write(escapedKey)
		buf.WriteByte(':')
		buf.WriteString(placeHolder)
	}
	buf.WriteByte('}')
	p.template.data = buf.Bytes()
	p.variables = variables
	return nil
}

func (p *HypDSPlanner) ConfigureFetch() resolve.FetchConfiguration {
	fnJson, err := utils.JsonSerialize(p.template.function)
	if err != nil {
		logger.Error(p.ctx).Err(err).Msg("Error serializing json while configuring graphql fetch.")
		return resolve.FetchConfiguration{}
	}

	// Note: we have to build the rest of the template manually, because the data field may
	// contain placeholders for variables, such as $$0$$ which are not valid in JSON.
	// They are replaced with the actual values by the time Load is called.
	inputTemplate := fmt.Sprintf(`{"fn":%s,"data":%s}`, fnJson, p.template.data)

	return resolve.FetchConfiguration{
		Input:     inputTemplate,
		Variables: p.variables,
		DataSource: &ModusDataSource{
			WasmHost: p.config.WasmHost,
		},
		PostProcessing: resolve.PostProcessingConfiguration{
			SelectResponseDataPath:   []string{"data"},
			SelectResponseErrorsPath: []string{"errors"},
		},
	}
}

func (p *HypDSPlanner) ConfigureSubscription() plan.SubscriptionConfiguration {
	panic("subscription not implemented")
}
