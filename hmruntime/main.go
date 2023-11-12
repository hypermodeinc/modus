package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"unicode/utf16"
	"unsafe"

	"github.com/dgraph-io/gqlparser/ast"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// TODO: abstract AssemblyScript-specific details

var runtime wazero.Runtime

type functionInfo struct {
	Module   *wasm.Module
	Function *wasm.Function
	Schema   *functionSchemaInfo
}

// map of resolver to registered function info
// TODO: this probably isn't robust enough for production
var functionsMap = make(map[string]functionInfo)

func main() {
	ctx := context.Background()

	// Parse command-line flags
	var port = flag.Int("port", 8686, "The HTTP port to listen on.")
	flag.Parse()

	// Initialize the WebAssembly runtime
	runtime = initWasmRuntime(ctx)
	defer runtime.Close(ctx)

	// Load plugins
	// TODO: This will need work:
	// - Plugins should probably be loaded from a repository, not from disk.
	// - We'll need to figure out how to handle plugin updates.
	// - We'll need to figure out hot/warm/cold plugin loading.
	_, err := loadPlugin(ctx, "hmplugin1") // for now just hardcoded
	if err != nil {
		log.Fatalln(err)
	}

	// Start the HTTP server
	fmt.Printf("Listening on port %d...\n", *port)
	http.HandleFunc("/graphql-worker", handleRequest)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	log.Fatalln(err)

	// TODO: Shutdown gracefully
}

func initWasmRuntime(ctx context.Context) wazero.Runtime {

	// Create the runtime
	cfg := wazero.NewRuntimeConfig().
		WithCloseOnContextDone(true)
	runtime := wazero.NewRuntimeWithConfig(ctx, cfg)

	// Enable WASI support
	wasi.MustInstantiate(ctx, runtime)

	// TODO: Define host functions

	return runtime
}

func loadPlugin(ctx context.Context, name string) (wasm.Module, error) {

	// Load the plugin plugin.
	// TODO: Load the plugin from some repository instead of disk.
	path := "../plugins/as/" + name + "/build/release.wasm"
	plugin, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load the plugin: %v", err)
	}

	cfg := wazero.NewModuleConfig().
		WithStdout(os.Stdout).WithStderr(os.Stderr)

	// Instantiate the plugin as a module.
	// NOTE: This will also invoke the plugin's `_start` function,
	// which will call any top-level code in the plugin.
	mod, err := runtime.InstantiateWithConfig(ctx, plugin, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate the plugin: %v", err)
	}

	fmt.Printf("Loaded plugin \"%s\"\n", name)

	// Get the GraphQL schema from Dgraph and use it to register the functions in this plugin.
	schema, err := getGQLSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get GraphQL schema: %v", err)
	}
	infos := getFunctionSchemaInfos(schema)
	for _, info := range infos {
		err = registerFunction(ctx, mod, info)
		if err != nil {
			fmt.Printf("Failed to register function \"%s\": %v\n", info.FunctionName(), err)
		}
	}

	return mod, nil
}

func registerFunction(ctx context.Context, mod wasm.Module, schema functionSchemaInfo) error {

	// Find the function in the module.
	fnName := schema.FunctionName()
	fn := mod.ExportedFunction(fnName)
	if fn == nil {
		return fmt.Errorf("function %s not found in module", fnName)
	}

	// Validate the function info.
	info := functionInfo{&mod, &fn, &schema}
	err := validateFunction(info)
	if err != nil {
		return fmt.Errorf("function %s is invalid: %v", fnName, err)
	}

	// Save the function info into the map.
	// TODO: this presumes there's no naming conflicts
	resolver := schema.Resolver()
	functionsMap[resolver] = info

	fmt.Printf("Registered function \"%s\" for resolver \"%s\"\n", fnName, resolver)
	return nil
}

func validateFunction(info functionInfo) error {
	// TODO: validate that the function definition matches the schema
	return nil
}

func callFunction(ctx context.Context, info functionInfo, inputs map[string]any) (any, error) {

	mod := *info.Module
	mem := mod.Memory()

	fn := *info.Function
	def := fn.Definition()
	paramTypes := def.ParamTypes()

	schema := *info.Schema
	fnName := schema.FunctionName()
	resolver := schema.Resolver()

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
			return nil, fmt.Errorf("parameter %s is invalid: %v", arg.Name, err)
		}

		params[i] = param
	}

	// Call the wasm function
	fmt.Printf("Calling function \"%s\" for resolver \"%s\"\n", fnName, resolver)
	res, err := fn.Call(ctx, params...)
	if err != nil {
		return nil, err
	}

	// Get the result
	result, err := convertResult(mem, *schema.FieldDef.Type, def.ResultTypes()[0], res[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %v", err)
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

	info := functionsMap[req.Resolver]
	if info.Function == nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("No function registered for resolver \"%s\"", req.Resolver)
		return
	}

	fnName := info.Schema.FunctionName()
	ctx := r.Context()

	if req.Args != nil {

		// Call the function, passing in the args from the request
		result, err := callFunction(ctx, info, req.Args)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error calling function \"%s\": %v", fnName, err)
			return
		}

		// Handle no result due to void return type
		if result == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Write the result
		isJson := info.Schema.FieldDef.Type.NamedType == ""
		err = writeDataAsJson(w, result, isJson)
		if err != nil {
			log.Println(err)
		}

	} else if req.Parents != nil {

		results := make([]any, len(req.Parents))

		// Call the function for each parent
		for i, parent := range req.Parents {
			results[i], err = callFunction(ctx, info, parent)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("Error calling function \"%s\": %v", fnName, err)
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
			w.WriteHeader(http.StatusInternalServerError)
			return fmt.Errorf("failed to serialize result data")
		}

		return nil
	}

	output, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("failed to serialize result data:", err)
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

	case "String", "Id", "":
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

		if res == 1 {
			return true, nil
		} else {
			return false, nil
		}

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

	case "String", "Id", "":
		// Note, strings are passed as a pointer to a string in wasm memory
		if wasmType != wasm.ValueTypeI32 {
			return nil, fmt.Errorf("return type was not a pointer")
		}

		return readString(mem, uint32(res))

	default:
		return nil, fmt.Errorf("unknown return type")
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
