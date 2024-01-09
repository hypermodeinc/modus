/*
 * Copyright 2023 Hypermode, Inc.
 */
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf16"
	"unsafe"

	"github.com/dgraph-io/gqlparser/ast"
	"github.com/google/uuid"
	"github.com/radovskyb/watcher"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

// TODO: abstract AssemblyScript-specific details

var runtime wazero.Runtime

// map that holds the compiled modules for each plugin
var compiledModules = make(map[string]wazero.CompiledModule)

// map that holds the function info for each resolver
var functionsMap = make(map[string]functionInfo)

// channel and flag used to signal the HTTP server
var serverReady chan bool = make(chan bool)
var serverWaiting = true

var dgraphUrl *string
var pluginsPath *string

func main() {
	ctx := context.Background()

	// Parse command-line flags
	var port = flag.Int("port", 8686, "The HTTP port to listen on.")
	dgraphUrl = flag.String("dgraph", "http://localhost:8080", "The Dgraph url to connect to.")
	pluginsPath = flag.String("plugins", "../plugins/as", "The path to the plugins directory.")
	flag.Parse()

	// Initialize the WebAssembly runtime
	var err error
	runtime, err = initWasmRuntime(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer runtime.Close(ctx)

	// Load plugins
	err = loadPlugins(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Watch for registration requests
	monitorRegistration(ctx)

	// Watch for schema changes
	monitorGqlSchema(ctx)

	// Watch for plugin changes
	err = watchPluginDirectory(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Start the HTTP server when we're ready
	<-serverReady
	serverWaiting = false
	fmt.Printf("Listening on port %d...\n", *port)
	http.HandleFunc("/graphql-worker", handleRequest)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	log.Fatalln(err)

	// TODO: Shutdown gracefully
}

func loadPlugins(ctx context.Context) error {
	entries, err := os.ReadDir(*pluginsPath)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}

	for _, entry := range entries {

		// Determine if the entry represents a plugin.
		var pluginName string
		entryName := entry.Name()
		if entry.IsDir() {
			pluginName = entryName
			path := fmt.Sprintf("%s/%s/build/debug.wasm", *pluginsPath, pluginName)
			if _, err := os.Stat(path); err != nil {
				continue
			}
		} else if strings.HasSuffix(entryName, ".wasm") {
			pluginName = strings.TrimSuffix(entryName, ".wasm")
		} else {
			continue
		}

		// Load the plugin
		err := loadPluginModule(ctx, pluginName)
		if err != nil {
			log.Printf("Failed to load plugin '%s': %v\n", pluginName, err)
		}
	}

	return nil
}

func registerFunctions(gqlSchema string) error {

	// Get the function schema from the GraphQL schema.
	funcSchemas, err := getFunctionSchema(gqlSchema)
	if err != nil {
		return err
	}

	// Build a map of resolvers to function info, including the plugin name.
	// If there are function name conflicts between plugins, the last plugin loaded wins.
	for pluginName, cm := range compiledModules {
		for _, schema := range funcSchemas {
			for _, fn := range cm.ExportedFunctions() {
				fnName := fn.ExportNames()[0]
				if strings.EqualFold(fnName, schema.FunctionName()) {
					info := functionInfo{pluginName, schema}
					resolver := schema.Resolver()
					oldInfo, existed := functionsMap[resolver]
					if existed && reflect.DeepEqual(oldInfo, info) {
						continue
					}
					functionsMap[resolver] = info
					if existed {
						fmt.Printf("Re-registered %s to use %s in %s\n", resolver, fnName, pluginName)
					} else {
						fmt.Printf("Registered %s to use %s in %s\n", resolver, fnName, pluginName)
					}
				}
			}
		}
	}

	// Cleanup any previously registered functions that are no longer in the schema or loaded modules.
	for resolver, info := range functionsMap {
		foundSchema := false
		for _, schema := range funcSchemas {
			if strings.EqualFold(info.FunctionName(), schema.FunctionName()) {
				foundSchema = true
				break
			}
		}
		_, foundModule := compiledModules[info.PluginName]
		if !foundSchema || !foundModule {
			delete(functionsMap, resolver)
			fmt.Printf("Unregistered old function '%s' for resolver '%s'\n", info.FunctionName(), resolver)
		}
	}

	// If the HTTP server is waiting, signal that we're ready.
	if serverWaiting {
		serverReady <- true
	}

	return nil
}

func initWasmRuntime(ctx context.Context) (wazero.Runtime, error) {

	// Create the runtime
	cfg := wazero.NewRuntimeConfig().
		WithCloseOnContextDone(true)
	runtime := wazero.NewRuntimeWithConfig(ctx, cfg)

	// Connect WASI host functions
	err := instantiateWasiFunctions(ctx, runtime)
	if err != nil {
		return nil, err
	}

	// Connect Hypermode host functions
	err = instantiateHostFunctions(ctx, runtime)
	if err != nil {
		return nil, err
	}

	return runtime, nil
}

func loadPluginModule(ctx context.Context, name string) error {
	cmOld, reloading := compiledModules[name]
	if reloading {
		fmt.Printf("Reloading plugin '%s'\n", name)
	} else {
		fmt.Printf("Loading plugin '%s'\n", name)
	}

	path, err := getPathForPlugin(name)
	if err != nil {
		return fmt.Errorf("failed to get path for plugin: %w", err)
	}

	plugin, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to load the plugin: %w", err)
	}

	// Compile the plugin into a module.
	cm, err := runtime.CompileModule(ctx, plugin)
	if err != nil {
		return fmt.Errorf("failed to compile the plugin: %w", err)
	}

	// Store the compiled module for later retrieval.
	compiledModules[name] = cm

	// When reloading, close the old module.
	if reloading {
		cmOld.Close(ctx)
	}

	return nil
}

func unloadPluginModule(ctx context.Context, name string) error {
	cmOld, found := compiledModules[name]
	if !found {
		return fmt.Errorf("plugin not found '%s'", name)
	}

	fmt.Printf("Unloading plugin '%s'\n", name)
	delete(compiledModules, name)
	cmOld.Close(ctx)

	return nil
}

type buffers struct {
	Stdout *strings.Builder
	Stderr *strings.Builder
}

func getModuleInstance(ctx context.Context, pluginName string) (wasm.Module, buffers, error) {

	// Create string buffers to capture stdout and stderr.
	// Still write to the console, but also capture the output in the buffers.
	buf := buffers{&strings.Builder{}, &strings.Builder{}}
	wOut := io.MultiWriter(os.Stdout, buf.Stdout)
	wErr := io.MultiWriter(os.Stderr, buf.Stderr)

	// Get the compiled module.
	compiled, ok := compiledModules[pluginName]
	if !ok {
		return nil, buf, fmt.Errorf("no compiled module found for plugin '%s'", pluginName)
	}

	// Configure the module instance.
	cfg := wazero.NewModuleConfig().
		WithName(pluginName + "_" + uuid.NewString()).
		WithStdout(wOut).WithStderr(wErr)

	// Instantiate the plugin as a module.
	// NOTE: This will also invoke the plugin's `_start` function,
	// which will call any top-level code in the plugin.
	mod, err := runtime.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		return nil, buf, fmt.Errorf("failed to instantiate the plugin module: %w", err)
	}

	return mod, buf, nil
}

func getPathForPlugin(name string) (string, error) {

	// Normally the plugin will be directly in the plugins directory, by filename.
	path := *pluginsPath + "/" + name + ".wasm"
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	// For local development, the plugin will be in a subdirectory and we'll use the debug.wasm file.
	path = *pluginsPath + "/" + name + "/build/debug.wasm"
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("compiled wasm file not found for plugin '%s'", name)
}

func getPluginNameFromPath(path string) (string, error) {
	if !strings.HasSuffix(path, ".wasm") {
		return "", fmt.Errorf("path does not point to a wasm file: %s", path)
	}

	parts := strings.Split(path, "/")

	// For local development
	if strings.HasSuffix(path, "/build/debug.wasm") {
		return parts[len(parts)-3], nil
	} else if strings.HasSuffix(path, "/build/release.wasm") {
		return "", nil
	}

	return strings.TrimSuffix(parts[len(parts)-1], ".wasm"), nil
}

func watchPluginDirectory(ctx context.Context) error {
	w := watcher.New()
	w.AddFilterHook(watcher.RegexFilterHook(regexp.MustCompile(`^.+\.wasm$`), false))

	go func() {
		for {
			select {
			case evt := <-w.Event:

				pluginName, err := getPluginNameFromPath(evt.Path)
				if err != nil {
					log.Printf("failed to get plugin name: %v\n", err)
				}
				if pluginName == "" {
					continue
				}

				switch evt.Op {
				case watcher.Create, watcher.Write:
					err = loadPluginModule(ctx, pluginName)
					if err != nil {
						log.Printf("failed to load plugin: %v\n", err)
					}
				case watcher.Remove:
					err = unloadPluginModule(ctx, pluginName)
					if err != nil {
						log.Printf("failed to unload plugin: %v\n", err)
					}
				}

				// Signal that we need to register functions
				register <- true

			case err := <-w.Error:
				log.Fatalf("failure while watching plugin directory: %v\n", err)
			case <-w.Closed:
				return
			case <-ctx.Done():
				w.Close()
				return
			}
		}
	}()

	err := w.AddRecursive(*pluginsPath)
	if err != nil {
		return fmt.Errorf("failed to watch plugins directory: %w", err)
	}

	go func() {
		err = w.Start(time.Second * 1)
		if err != nil {
			log.Fatalf("failed to start file watcher: %v\n", err)
		}
	}()

	return nil
}

func callFunction(ctx context.Context, mod wasm.Module, info functionInfo, inputs map[string]any) (any, error) {
	fnName := info.FunctionName()
	fn := mod.ExportedFunction(fnName)
	def := fn.Definition()
	paramTypes := def.ParamTypes()
	schema := info.Schema

	// Get parameters to pass as input to the plugin function
	// Note that we can't use def.ParamNames() because they are only available when the plugin
	// is compiled in debug mode. They're striped by optimization in release mode.
	// Instead, we can use the argument names from the schema.
	// Also note, that the order of the arguments from schema should match order of params in wasm.
	params := make([]uint64, len(paramTypes))
	for i, arg := range schema.FunctionArgs() {
		val := inputs[arg.Name]
		if val == nil {
			return nil, fmt.Errorf("parameter %s is missing", arg.Name)
		}

		param, err := convertParam(ctx, mod, *arg.Type, paramTypes[i], val)
		if err != nil {
			return nil, fmt.Errorf("parameter %s is invalid: %w", arg.Name, err)
		}

		params[i] = param
	}

	// Call the wasm function
	fmt.Printf("Calling function '%s' for resolver '%s'\n", fnName, schema.Resolver())
	res, err := fn.Call(ctx, params...)
	if err != nil {
		return nil, err
	}

	// Get the result
	mem := mod.Memory()
	result, err := convertResult(mem, *schema.FieldDef.Type, def.ResultTypes()[0], res[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %w", err)
	}

	return result, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	// Decode the request body
	var req graphRequest
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Failed to decode request body: ", err)
		return
	}

	// Get the function info for the resolver
	info, ok := functionsMap[req.Resolver]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("No function registered for resolver '%s'", req.Resolver)
		return
	}

	// Get a module instance for this request.
	// Each request will get its own instance of the plugin module,
	// so that we can run multiple requests in parallel without risk
	// of corrupting the module's memory.
	ctx := r.Context()
	mod, buf, err := getModuleInstance(ctx, info.PluginName)
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, err)
		return
	}
	defer mod.Close(ctx)

	fnName := info.FunctionName()
	if req.Args != nil {

		// Call the function, passing in the args from the request
		result, err := callFunction(ctx, mod, info, req.Args)
		if err != nil {
			err := fmt.Errorf("error calling function '%s': %w", fnName, err)
			log.Println(err)
			writeErrorResponse(w, err, buf.Stdout.String(), buf.Stderr.String())
			return
		}

		// Handle no result due to void return type
		if result == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Determine if the result is already JSON
		isJson := false
		fieldType := info.Schema.FieldDef.Type.NamedType
		if _, ok := result.(string); ok && fieldType != "String" {
			isJson = true
		}

		// Write the result
		err = writeDataAsJson(w, result, isJson)
		if err != nil {
			log.Println(err)
		}

	} else if req.Parents != nil {

		results := make([]any, len(req.Parents))

		// Call the function for each parent
		for i, parent := range req.Parents {
			results[i], err = callFunction(ctx, mod, info, parent)
			if err != nil {
				err := fmt.Errorf("error calling function '%s': %w", fnName, err)
				log.Println(err)
				writeErrorResponse(w, err, buf.Stdout.String(), buf.Stderr.String())
				return
			}
		}

		// Write the result
		isJson := info.Schema.FieldDef.Type.NamedType == ""
		err = writeDataAsJson(w, results, isJson)
		if err != nil {
			log.Println(err)
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Request must have either args or parents.")
	}
}

func writeErrorResponse(w http.ResponseWriter, err error, msgs ...string) {
	w.WriteHeader(http.StatusInternalServerError)

	// Dgraph lambda expects a JSON response similar to a GraphQL error response
	w.Header().Set("Content-Type", "application/json")
	resp := HMErrorResponse{Errors: []HMError{}}

	// Emit messages first
	for _, msg := range msgs {
		for _, line := range strings.Split(msg, "\n") {
			if len(line) > 0 {
				resp.Errors = append(resp.Errors, HMError{Message: line})
			}
		}
	}

	// Emit the error last
	resp.Errors = append(resp.Errors, HMError{Message: err.Error()})

	json.NewEncoder(w).Encode(resp)
}

func writeDataAsJson(w http.ResponseWriter, data any, isJson bool) error {

	if isJson {

		switch data := data.(type) {
		case string:
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(data))
		case []string:
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte{'['})
			for i, s := range data {
				if i > 0 {
					w.Write([]byte{','})
				}
				w.Write([]byte(s))
			}
			w.Write([]byte{']'})
		default:
			err := fmt.Errorf("unexpected result type: %T", data)
			log.Println(err)
			writeErrorResponse(w, err)
		}

		return nil
	}

	output, err := json.Marshal(data)
	if err != nil {
		err := fmt.Errorf("failed to serialize result data: %s", err)
		log.Println(err)
		writeErrorResponse(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return nil
}

func convertParam(ctx context.Context, mod wasm.Module, schemaType ast.Type, wasmType wasm.ValueType, val any) (uint64, error) {

	switch schemaType.NamedType {

	case "Boolean":
		b, ok := val.(bool)
		if !ok {
			return 0, fmt.Errorf("input value is not a bool")
		}

		// Note, booleans are passed as i32 in wasm
		if wasmType != wasm.ValueTypeI32 {
			return 0, fmt.Errorf("parameter is not defined as a bool on the function")
		}

		if b {
			return 1, nil
		} else {
			return 0, nil
		}

	case "Int":
		n, err := val.(json.Number).Int64()
		if err != nil {
			return 0, fmt.Errorf("input value is not an int")
		}

		switch wasmType {
		case wasm.ValueTypeI32:
			return wasm.EncodeI32(int32(n)), nil
		case wasm.ValueTypeI64:
			return wasm.EncodeI64(n), nil
		default:
			return 0, fmt.Errorf("parameter is not defined as an int on the function")
		}

	case "Float":
		n, err := val.(json.Number).Float64()
		if err != nil {
			return 0, fmt.Errorf("input value is not a float")
		}

		switch wasmType {
		case wasm.ValueTypeF32:
			return wasm.EncodeF32(float32(n)), nil
		case wasm.ValueTypeF64:
			return wasm.EncodeF64(n), nil
		default:
			return 0, fmt.Errorf("parameter is not defined as a float on the function")
		}

	case "String", "ID", "":
		s, ok := val.(string)
		if !ok {
			return 0, fmt.Errorf("input value is not a string")
		}

		// Note, strings are passed as a pointer to a string in wasm memory
		if wasmType != wasm.ValueTypeI32 {
			return 0, fmt.Errorf("parameter is not defined as a string on the function")
		}

		ptr := writeString(ctx, mod, s)
		return uint64(ptr), nil

	default:
		return 0, fmt.Errorf("unknown parameter type: %s", schemaType.NamedType)
	}
}

func convertResult(mem wasm.Memory, schemaType ast.Type, wasmType wasm.ValueType, res uint64) (any, error) {

	switch schemaType.NamedType {

	case "Boolean":
		if wasmType != wasm.ValueTypeI32 {
			return nil, fmt.Errorf("return type is not defined as an bool on the function")
		}

		return res != 0, nil

	case "Int":

		// TODO: Do we need to handle unsigned ints differently?

		switch wasmType {
		case wasm.ValueTypeI32:
			return wasm.DecodeI32(res), nil

		case wasm.ValueTypeI64:
			return int64(res), nil

		default:
			return nil, fmt.Errorf("return type is not defined as an int on the function")
		}

	case "Float":

		switch wasmType {
		case wasm.ValueTypeF32:
			return wasm.DecodeF32(res), nil

		case wasm.ValueTypeF64:
			return wasm.DecodeF64(res), nil

		default:
			return nil, fmt.Errorf("return type is not defined as a float on the function")
		}

	case "ID":
		return nil, fmt.Errorf("the ID scalar is not allowed for function return types (use String instead)")

	default:
		// The return type is either a string, or an object that should be serialized as JSON.
		// Strings are passed as a pointer to a string in wasm memory
		if wasmType != wasm.ValueTypeI32 {
			return nil, fmt.Errorf("return type was not a pointer")
		}

		return readString(mem, uint32(res))
	}
}

type graphRequest struct {
	AccessToken string `json:"X-Dgraph-AccessToken"`
	AuthHeader  struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"authHeader"`
	Args     map[string]any   `json:"args"`
	Parents  []map[string]any `json:"parents"`
	Resolver string           `json:"resolver"`
}

func writeString(ctx context.Context, mod wasm.Module, s string) uint32 {
	buf := encodeUTF16(s)
	ptr := allocateWasmMemory(ctx, mod, len(buf), asString)
	mod.Memory().Write(ptr, buf)
	return ptr
}

func readString(mem wasm.Memory, offset uint32) (string, error) {

	// AssemblyScript managed objects have their classid stored 8 bytes before the offset.
	// See https://www.assemblyscript.org/runtime.html#memory-layout

	// Read the class id.
	id, ok := mem.ReadUint32Le(offset - 8)
	if !ok {
		return "", fmt.Errorf("failed to read class id of the WASM object")
	}

	// Make sure the pointer is to a string.
	if id != uint32(asString) {
		return "", fmt.Errorf("pointer is not to a string")
	}

	// Read from the buffer and decode it as a string.
	buf, err := readBuffer(mem, offset)
	if err != nil {
		return "", err
	}

	return decodeUTF16(buf), nil
}

func readBuffer(mem wasm.Memory, offset uint32) ([]byte, error) {

	// The length of AssemblyScript managed objects is stored 4 bytes before the offset.
	// See https://www.assemblyscript.org/runtime.html#memory-layout

	// Read the length.
	len, ok := mem.ReadUint32Le(offset - 4)
	if !ok {
		return nil, fmt.Errorf("failed to read buffer length")
	}

	// Handle empty buffers.
	if len == 0 {
		return []byte{}, nil
	}

	// Now read the data into the buffer.
	buf, ok := mem.Read(offset, len)
	if !ok {
		return nil, fmt.Errorf("failed to read buffer data from WASM memory")
	}

	return buf, nil
}

// See https://www.assemblyscript.org/runtime.html#memory-layout
type asClass int64

const (
	asBytes  asClass = 1
	asString asClass = 2
)

func allocateWasmMemory(ctx context.Context, mod wasm.Module, len int, class asClass) uint32 {
	// Allocate a string to hold our buffer within the AssemblyScript module.
	// This uses the `__new` function exported by the AssemblyScript runtime, so it will be garbage collected.
	// See https://www.assemblyscript.org/runtime.html#interface
	newFn := mod.ExportedFunction("__new")
	res, _ := newFn.Call(ctx, uint64(len), uint64(class))
	return uint32(res[0])
}

func decodeUTF16(bytes []byte) string {

	// Make sure the buffer is valid.
	if len(bytes) == 0 || len(bytes)%2 != 0 {
		return ""
	}

	// Reinterpret []byte as []uint16 to avoid excess copying.
	// This works because we can presume the system is little-endian.
	ptr := unsafe.Pointer(&bytes[0])
	words := unsafe.Slice((*uint16)(ptr), len(bytes)/2)

	// Decode UTF-16 words to a UTF-8 string.
	str := string(utf16.Decode(words))
	return str
}

func encodeUTF16(str string) []byte {
	// Encode the UTF-8 string to UTF-16 words.
	words := utf16.Encode([]rune(str))

	// Reinterpret []uint16 as []byte to avoid excess copying.
	// This works because we can presume the system is little-endian.
	ptr := unsafe.Pointer(&words[0])
	bytes := unsafe.Slice((*byte)(ptr), len(words)*2)
	return bytes
}

type HMErrorResponse struct {
	Errors []HMError `json:"errors"`
}

type HMError struct {
	Message string `json:"message"`
}
