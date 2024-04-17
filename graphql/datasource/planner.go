/*
 * Copyright 2024 Hypermode, Inc.
 */

package datasource

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/logger"

	"github.com/wundergraph/graphql-go-tools/v2/pkg/ast"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/plan"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

type Planner[T Configuration] struct {
	ctx       context.Context
	config    Configuration
	visitor   *plan.Visitor
	variables resolve.Variables
	template  struct {
		function templateField
		data     []byte
	}
}

type templateField struct {
	Name   string          `json:"name"`
	Alias  string          `json:"alias,omitempty"`
	Fields []templateField `json:"fields,omitempty"`
}

func (t templateField) AliasOrName() string {
	if t.Alias != "" {
		return t.Alias
	}
	return t.Name
}

func (p *Planner[T]) UpstreamSchema(dataSourceConfig plan.DataSourceConfiguration[T]) (*ast.Document, bool) {
	return nil, false
}

func (p *Planner[T]) DownstreamResponseFieldAlias(downstreamFieldRef int) (alias string, exists bool) {
	return
}

func (p *Planner[T]) DataSourcePlanningBehavior() plan.DataSourcePlanningBehavior {
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
	}
}

func (p *Planner[T]) Register(visitor *plan.Visitor, configuration plan.DataSourceConfiguration[T], dspc plan.DataSourcePlannerConfiguration) error {
	p.visitor = visitor
	visitor.Walker.RegisterEnterFieldVisitor(p)
	p.config = Configuration(configuration.CustomConfiguration())
	return nil
}

func (p *Planner[T]) EnterField(ref int) {
	if !p.enclosingTypeIsRootNode() {
		return
	}

	p.captureTemplateField(ref, &p.template.function)
	err := p.captureInputData(ref)
	if err != nil {
		logger.Err(p.ctx, err).Msg("Error capturing input data.")
		return
	}
}

func (p *Planner[T]) enclosingTypeIsRootNode() bool {
	enclosingType := p.visitor.Walker.EnclosingTypeDefinition
	for _, node := range p.visitor.Operation.RootNodes {
		if node.Ref == enclosingType.Ref {
			return true
		}
	}
	return false
}

func (p *Planner[T]) captureTemplateField(ref int, tf *templateField) {
	operation := p.visitor.Operation
	tf.Name = operation.FieldNameString(ref)
	tf.Alias = operation.FieldAliasString(ref)
	if operation.FieldHasSelections(ref) {
		ssRef, ok := operation.FieldSelectionSet(ref)
		if ok {
			fieldRefs := operation.SelectionSetFieldSelections(ssRef)
			tf.Fields = make([]templateField, len(fieldRefs))
			for i, fieldRef := range fieldRefs {
				p.captureTemplateField(fieldRef, &tf.Fields[i])
			}
		}
	}
}

func (p *Planner[T]) captureInputData(fieldRef int) error {
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

		escapedKey, err := json.Marshal(argName)
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

func (p *Planner[T]) ConfigureFetch() resolve.FetchConfiguration {
	fnJson, err := json.Marshal(p.template.function)
	if err != nil {
		logger.Error(p.ctx).Err(err).Msg("Error marshalling json while configuring graphql fetch.")
		return resolve.FetchConfiguration{}
	}

	// Note: we have to build the rest of the template manually, because the data field may
	// contain placeholders for variables, such as $$0$$ which are not valid in JSON.
	// They are replaced with the actual values by the time Load is called.
	inputTemplate := fmt.Sprintf(`{"fn":%s,"data":%s}`, fnJson, p.template.data)

	return resolve.FetchConfiguration{
		Input:      inputTemplate,
		Variables:  p.variables,
		DataSource: Source{},
		PostProcessing: resolve.PostProcessingConfiguration{
			SelectResponseDataPath:   []string{"data"},
			SelectResponseErrorsPath: []string{"errors"},
		},
	}
}

func (p *Planner[T]) ConfigureSubscription() plan.SubscriptionConfiguration {
	panic("subscription not implemented")
}
