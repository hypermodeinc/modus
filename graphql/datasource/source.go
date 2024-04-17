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

	"github.com/rs/xid"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

const DataSourceName = "HypermodeFunctionsDataSource"

var errCallingFunction = fmt.Errorf("error calling function")

type callInfo struct {
	Function   string         `json:"fn"`
	Alias      string         `json:"alias"`
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
	info, ok := functions.Functions[ci.Function]
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
	errs := makeErrors(buf, ci)
	if len(errs) > 0 {
		var err error
		jsonErrors, err = json.Marshal(errs)
		if err != nil {
			return err
		}
	}

	if fnErr != nil {
		// If there are errors in the output buffers, return those by themselves
		if len(errs) > 0 {
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

	// Build and write the response
	if len(errs) > 0 {
		fmt.Fprintf(writer, `{"data":{"%s":%s},"errors":%s}`, ci.Alias, jsonData, jsonErrors)
	} else {
		fmt.Fprintf(writer, `{"data":{"%s":%s}}`, ci.Alias, jsonData)
	}

	return nil
}

func makeErrors(buf host.OutputBuffers, ci callInfo) []resolve.GraphQLError {
	errors := make([]resolve.GraphQLError, 0)
	for _, s := range strings.Split(buf.Stdout.String(), "\n") {
		if s != "" {
			errors = append(errors, makeError(s, ci))
		}
	}
	for _, s := range strings.Split(buf.Stderr.String(), "\n") {
		if s != "" {
			errors = append(errors, makeError(s, ci))
		}
	}
	return errors
}

func makeError(msg string, ci callInfo) resolve.GraphQLError {
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
		Path:    []string{ci.Alias},
	}

	if level != "" {
		e.Extensions = map[string]interface{}{
			"level": level,
		}
	}

	return e
}
