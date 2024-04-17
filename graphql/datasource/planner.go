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
	ctx          context.Context
	config       Configuration
	visitor      *plan.Visitor
	variables    resolve.Variables
	rootFieldRef int
	template     struct {
		function string
		alias    string
		data     []byte
	}
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
	visitor.Walker.RegisterEnterDocumentVisitor(p)
	visitor.Walker.RegisterEnterFieldVisitor(p)
	p.config = Configuration(configuration.CustomConfiguration())
	return nil
}

func (p *Planner[T]) EnterDocument(operation, definition *ast.Document) {
	p.rootFieldRef = -1
}

func (p *Planner[T]) EnterField(ref int) {
	if p.rootFieldRef == -1 {
		p.rootFieldRef = ref
	} else {
		// This is a nested field, we don't need to do anything here
		return
	}

	walker := p.visitor.Walker
	typeName := walker.EnclosingTypeDefinition.NameString(p.visitor.Definition)

	if typeName != "Query" {
		return
	}

	operation := p.visitor.Operation
	p.template.function = operation.FieldNameString(ref)
	p.template.alias = operation.FieldAliasOrNameString(ref)

	var err error
	p.template.data, p.variables, err = getTemplateDataAndVariables(operation, ref)
	if err != nil {
		logger.Error(p.ctx).Err(err).Msg("error getting info from AST for operation")
		return
	}
}

func getTemplateDataAndVariables(operation *ast.Document, fieldRef int) ([]byte, resolve.Variables, error) {
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
			return nil, nil, err
		}

		buf.Write(escapedKey)
		buf.WriteByte(':')
		buf.WriteString(placeHolder)
	}
	buf.WriteByte('}')
	return buf.Bytes(), variables, nil
}

func (p *Planner[T]) ConfigureFetch() resolve.FetchConfiguration {
	// Note: we have to build the JSON manually here, because the data field may
	// contain placeholders for variables, such as $$0$$ which are not valid in JSON.
	// They are replaced with the actual values by the time Load is called.
	inputTemplate := fmt.Sprintf(
		`{"fn":"%s","alias":"%s","data":%s}`,
		p.template.function,
		p.template.alias,
		p.template.data)

	return resolve.FetchConfiguration{
		Input:      string(inputTemplate),
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
