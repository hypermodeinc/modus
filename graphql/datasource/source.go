/*
 * Copyright 2024 Hypermode, Inc.
 */

package datasource

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"hmruntime/functions"
	"hmruntime/host"
	"hmruntime/logger"
	"hmruntime/utils"

	"github.com/buger/jsonparser"
	"github.com/rs/xid"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

const DataSourceName = "HypermodeFunctionsDataSource"

var errCallingFunction = fmt.Errorf("error calling function")

type callInfo struct {
	Function   templateField  `json:"fn"`
	Parameters map[string]any `json:"data"`
}

type Source struct{}

func (s Source) Load(ctx context.Context, input []byte, writer io.Writer) error {
	err := s.load(ctx, input, writer)
	if err != nil && !errors.Is(err, errCallingFunction) {
		// note: function call errors are already logged, so we don't log them again here
		logger.Err(ctx, err).Msg("Failed to load data.")
	}
	return err
}

func (s Source) load(ctx context.Context, input []byte, writer io.Writer) error {
	// Get the call info
	var ci callInfo
	dec := json.NewDecoder(bytes.NewReader(input))
	dec.UseNumber()
	err := dec.Decode(&ci)

	if err != nil {
		return fmt.Errorf("error getting function input: %w", err)
	}

	// Get the function info
	info, ok := functions.Functions[ci.Function.Name]
	if !ok {
		return fmt.Errorf("no function registered named %s", ci.Function)
	}

	// Add plugin to the context
	ctx = context.WithValue(ctx, utils.PluginContextKey, info.Plugin)

	// Add execution ID to the context
	executionId := xid.New().String()
	ctx = context.WithValue(ctx, utils.ExecutionIdContextKey, executionId)

	// TODO: We should return the execution id(s) in the response with X-Hypermode-ExecutionID headers.
	// There might be multiple execution ids if the request triggers multiple function calls.

	// Get a module instance for this request.
	// Each request will get its own instance of the plugin module,
	// so that we can run multiple requests in parallel without risk
	// of corrupting the module's memory.
	mod, buf, err := host.GetModuleInstance(ctx, info.Plugin)
	if err != nil {
		return fmt.Errorf("error getting module instance: %w", err)
	}
	defer mod.Close(ctx)

	// Call the function and get any errors that were written to the output buffers
	result, fnErr := functions.CallFunction(ctx, mod, info, ci.Parameters)

	var jsonErrors []byte
	errors := transformErrors(buf, ci)
	if len(errors) > 0 {
		var err error
		jsonErrors, err = json.Marshal(errors)
		if err != nil {
			return err
		}
	}

	if fnErr != nil {
		// If there are errors in the output buffers, return those by themselves
		if len(errors) > 0 {
			fmt.Fprintf(writer, `{"errors":%s}`, jsonErrors)
			return nil
		}

		// If not, then just return the error from the function call
		return fmt.Errorf("%w '%s': %w", errCallingFunction, ci.Function, err)
	}

	// Get the data as json from the result
	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	// Transform the data
	jsonData, err = transformData(jsonData, ci.Function)
	if err != nil {
		return err
	}

	// Build and write the response
	if len(errors) > 0 {
		fmt.Fprintf(writer, `{"data":%s,"errors":%s}`, jsonData, jsonErrors)
	} else {
		fmt.Fprintf(writer, `{"data":%s}`, jsonData)
	}

	return nil
}

func transformData(data []byte, tf templateField) ([]byte, error) {
	val, err := transformValue(data, tf)
	if err != nil {
		return nil, err
	}

	out := []byte(`{}`)
	return jsonparser.Set(out, val, tf.AliasOrName())
}

func transformValue(data []byte, tf templateField) ([]byte, error) {
	if len(tf.Fields) == 0 || len(data) == 0 {
		return data, nil
	}

	buf := bytes.Buffer{}

	switch data[0] {
	case '{': // object
		buf.WriteByte('{')
		for i, f := range tf.Fields {
			val, dataType, _, err := jsonparser.Get(data, f.Name)
			if err != nil {
				return nil, err
			}
			if dataType == jsonparser.String {
				val, err = json.Marshal(string(val))
				if err != nil {
					return nil, err
				}
			}
			val, err = transformValue(val, f)
			if err != nil {
				return nil, err
			}
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteByte('"')
			buf.WriteString(f.AliasOrName())
			buf.WriteString(`":`)
			buf.Write(val)
		}
		buf.WriteByte('}')

	case '[': // array
		buf.WriteByte('[')
		_, err := jsonparser.ArrayEach(data, func(val []byte, _ jsonparser.ValueType, _ int, _ error) {
			if buf.Len() > 1 {
				buf.WriteByte(',')
			}
			val, err := transformValue(val, tf)
			if err != nil {
				return
			}
			buf.Write(val)
		})
		if err != nil {
			return nil, err
		}

		buf.WriteByte(']')

	default:
		return nil, fmt.Errorf("expected object or array")
	}

	return buf.Bytes(), nil
}

func transformErrors(buf host.OutputBuffers, ci callInfo) []resolve.GraphQLError {
	errors := make([]resolve.GraphQLError, 0)
	for _, s := range strings.Split(buf.Stdout.String(), "\n") {
		if s != "" {
			errors = append(errors, transformError(s, ci))
		}
	}
	for _, s := range strings.Split(buf.Stderr.String(), "\n") {
		if s != "" {
			errors = append(errors, transformError(s, ci))
		}
	}
	return errors
}

func transformError(msg string, ci callInfo) resolve.GraphQLError {
	level := ""
	a := strings.SplitAfterN(msg, ": ", 2)
	if len(a) == 2 {
		switch a[0] {
		case "Debug: ":
			level = "debug"
			msg = a[1]
		case "Info: ":
			level = "info"
			msg = a[1]
		case "Warning: ":
			level = "warning"
			msg = a[1]
		case "Error: ":
			level = "error"
			msg = a[1]
		case "abort: ":
			level = "fatal"
			msg = a[1]
		}
	}

	e := resolve.GraphQLError{
		Message: msg,
		Path:    []string{ci.Function.AliasOrName()},
	}

	if level != "" {
		e.Extensions = map[string]interface{}{
			"level": level,
		}
	}

	return e
}
