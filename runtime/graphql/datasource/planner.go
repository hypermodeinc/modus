/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package datasource

import (
	"bytes"
	"context"
	"slices"
	"strings"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/sentryutils"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/tidwall/gjson"

	"github.com/wundergraph/graphql-go-tools/v2/pkg/ast"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/plan"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

type modusDataSourcePlanner struct {
	id        int
	ctx       context.Context
	config    plan.DataSourceConfiguration[ModusDataSourceConfig]
	visitor   *plan.Visitor
	variables resolve.Variables
	fields    map[int]fieldInfo
	template  struct {
		fieldInfo    *fieldInfo
		functionName string
		data         []byte
	}
}

type fieldInfo struct {
	ref        int         `json:"-"`
	Name       string      `json:"name"`
	Alias      string      `json:"alias,omitempty"`
	TypeName   string      `json:"type,omitempty"`
	ParentType string      `json:"parentType,omitempty"`
	Fields     []fieldInfo `json:"fields,omitempty"`
	IsMapType  bool        `json:"isMapType,omitempty"`
	fieldRefs  []int       `json:"-"`
	depth      int         `json:"-"`
}

func (t *fieldInfo) AliasOrName() string {
	if t.Alias != "" {
		return t.Alias
	}
	return t.Name
}

func (p *modusDataSourcePlanner) SetID(id int) {
	p.id = id
}

func (p *modusDataSourcePlanner) ID() (id int) {
	return p.id
}

func (p *modusDataSourcePlanner) UpstreamSchema(dataSourceConfig plan.DataSourceConfiguration[ModusDataSourceConfig]) (*ast.Document, bool) {
	return nil, false
}

func (p *modusDataSourcePlanner) DownstreamResponseFieldAlias(downstreamFieldRef int) (alias string, exists bool) {
	return
}

func (p *modusDataSourcePlanner) Register(visitor *plan.Visitor, configuration plan.DataSourceConfiguration[ModusDataSourceConfig], dspc plan.DataSourcePlannerConfiguration) error {
	p.visitor = visitor
	visitor.Walker.RegisterEnterDocumentVisitor(p)
	visitor.Walker.RegisterEnterFieldVisitor(p)
	visitor.Walker.RegisterLeaveDocumentVisitor(p)
	p.config = configuration

	return nil
}

func (p *modusDataSourcePlanner) EnterDocument(operation, definition *ast.Document) {
	p.fields = make(map[int]fieldInfo, len(operation.Fields))
}

func (p *modusDataSourcePlanner) EnterField(ref int) {
	ds := p.config.(plan.DataSource)
	config := p.config.CustomConfiguration()

	// Capture information about every field in the operation.
	f := p.captureField(ref)
	p.fields[ref] = *f

	// Capture input data from only the fields that represent function calls or event subscriptions.
	// These fields are identified by their depth (3), and the fact that they have
	// a parent type that is a root node in the data source configuration.
	//
	// The depth is 3 because the tree structure of the GraphQL operation looks like this:
	//   0: root node (query, mutation, or subscription)
	//   1: selection set node
	//   2: this field's node
	if f.depth == 3 && ds.HasRootNode(f.ParentType, f.Name) {
		p.template.fieldInfo = f
		p.template.functionName = config.FieldsToFunctions[f.Name]
		if err := p.captureInputData(ref); err != nil {
			const msg = "Error capturing graphql input data."
			sentryutils.CaptureError(p.ctx, err, msg)
			logger.Error(p.ctx, err).Msg(msg)
			return
		}
	}
}

func (p *modusDataSourcePlanner) LeaveDocument(operation, definition *ast.Document) {
	// Stitch the captured fields together to form a tree.
	p.stitchFields(p.template.fieldInfo)
}

func (p *modusDataSourcePlanner) stitchFields(f *fieldInfo) {
	if f == nil || len(f.fieldRefs) == 0 {
		return
	}

	f.Fields = make([]fieldInfo, len(f.fieldRefs))
	for i, ref := range f.fieldRefs {
		field := p.fields[ref]
		p.stitchFields(&field)
		f.Fields[i] = field
	}
}

func (p *modusDataSourcePlanner) captureField(ref int) *fieldInfo {
	operation := p.visitor.Operation
	definition := p.visitor.Definition
	walker := p.visitor.Walker
	config := p.config.CustomConfiguration()

	f := &fieldInfo{
		ref:   ref,
		depth: walker.Depth,
		Name:  operation.FieldNameString(ref),
		Alias: operation.FieldAliasString(ref),
	}

	def, ok := walker.FieldDefinition(ref)
	if ok {
		f.TypeName = definition.FieldDefinitionTypeNameString(def)
		f.ParentType = walker.EnclosingTypeDefinition.NameString(definition)
		f.IsMapType = slices.Contains(config.MapTypes, f.TypeName)
	}

	if operation.FieldHasSelections(ref) {
		ssRef, ok := operation.FieldSelectionSet(ref)
		if ok {
			f.fieldRefs = operation.SelectionSetFieldSelections(ssRef)
		}
	}

	return f
}

func (p *modusDataSourcePlanner) captureInputData(fieldRef int) error {
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

func (p *modusDataSourcePlanner) getInputTemplate() (string, error) {

	fieldInfoJson, err := utils.JsonSerialize(p.template.fieldInfo)
	if err != nil {
		return "", err
	}

	// Note: we have to build the rest of the template manually, because the data field may
	// contain placeholders for variables, such as $$0$$ which are not valid in JSON.
	// They are replaced with the actual values when the input is rendered.

	b := &strings.Builder{}
	b.WriteString(`{"field":`)
	b.Write(fieldInfoJson)

	if len(p.template.functionName) > 0 {
		b.WriteString(`,"function":`)
		b.Write(gjson.AppendJSONString(nil, p.template.functionName))
	}

	b.WriteString(`,"data":`)
	b.Write(p.template.data)

	b.WriteByte('}')

	return b.String(), nil
}

func (p *modusDataSourcePlanner) ConfigureFetch() resolve.FetchConfiguration {
	input, err := p.getInputTemplate()
	if err != nil {
		const msg = "Error creating input template for graphql data source."
		sentryutils.CaptureError(p.ctx, err, msg)
		logger.Error(p.ctx, err).Msg(msg)
		return resolve.FetchConfiguration{}
	}

	config := p.config.CustomConfiguration()

	return resolve.FetchConfiguration{
		Input:     input,
		Variables: p.variables,
		DataSource: &functionsDataSource{
			WasmHost: config.WasmHost,
		},
		PostProcessing: resolve.PostProcessingConfiguration{
			SelectResponseDataPath:   []string{"data"},
			SelectResponseErrorsPath: []string{"errors"},
		},
	}
}

func (p *modusDataSourcePlanner) ConfigureSubscription() plan.SubscriptionConfiguration {
	input, err := p.getInputTemplate()
	if err != nil {
		const msg = "Error creating input template for graphql data source."
		sentryutils.CaptureError(p.ctx, err, msg)
		logger.Error(p.ctx, err).Msg(msg)
		return plan.SubscriptionConfiguration{}
	}

	return plan.SubscriptionConfiguration{
		Input:      input,
		Variables:  p.variables,
		DataSource: &eventsDataSource{},
		PostProcessing: resolve.PostProcessingConfiguration{
			SelectResponseDataPath:   []string{"data"},
			SelectResponseErrorsPath: []string{"errors"},
		},
	}
}
