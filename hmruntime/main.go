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

	"hmruntime/protos"

	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"google.golang.org/protobuf/proto"
)

// TODO: standardize logging and output error handling throughout
// TODO: abstract AssemblyScript-specific details

var runtime wazero.Runtime

type FuncInfo struct {
	Module   *wasm.Module
	Function *protos.HMFunction
}

// map of resolver to registered function info
// TODO: this probably isn't robust enough for production
var functionsMap = make(map[string]FuncInfo)

func main() {
	// Parse command-line flags
	var port = flag.Int("port", 8686, "The HTTP port to listen on.")
	flag.Parse()

	// Initialize the WebAssembly runtime
	ctx := context.Background()
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
	runtime := wazero.NewRuntime(ctx)

	// Enable WASI support
	wasi.MustInstantiate(ctx, runtime)

	// Define host functions
	_, err := runtime.NewHostModuleBuilder("host").
		NewFunctionBuilder().WithFunc(registerFunction).Export("registerFunction").
		Instantiate(ctx)
	if err != nil {
		panic(err)
	}

	return runtime
}

func registerFunction(ctx context.Context, mod wasm.Module, msg uint32) {
	mem := mod.Memory()
	buf := readBuffer(mem, msg)

	var info protos.HMFunction
	err := proto.Unmarshal(buf, &info)
	if err != nil {
		fmt.Println("error deserializing function proto:", err)
		return
	}

	// Make sure the named function exists in the module and parameters match up.
	fn := mod.ExportedFunction(info.Name)
	if fn == nil {
		fmt.Printf("function %s not found in module\n", info.Name)
		return
	}

	def := fn.Definition()
	paramTypes := def.ParamTypes()
	resultTypes := def.ResultTypes()

	if len(info.Parameters) != len(paramTypes) {
		fmt.Printf("function %s has %d parameters, but %d were registered\n", info.Name, len(paramTypes), len(info.Parameters))
		return
	}

	// Validate number of return types.
	// NOTE: It's currently not possible to have more than one return type.
	if info.ReturnType == protos.HMType_VOID && len(resultTypes) != 0 {
		fmt.Printf("function %s has no return type, but %d were registered\n", info.Name, len(resultTypes))
		return
	} else if info.ReturnType != protos.HMType_VOID && len(resultTypes) != 1 {
		fmt.Printf("function %s has one return type, but %d were registered\n", info.Name, len(resultTypes))
		return
	}

	// Save the function and module info into the map.
	// TODO: this presumes there's no naming conflicts
	functionsMap[info.Resolver] = FuncInfo{&mod, &info}
}

func loadPlugin(ctx context.Context, name string) (wasm.Module, error) {

	// Load the plugin plugin.
	// TODO: Load the plugin from some repository instead of disk.
	path := "../plugins/as/" + name + "/build/release.wasm"
	plugin, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load the plugin: %v", err)
	}

	// Instantiate the plugin as a module.
	// This will invoke the plugin's _start function,
	// which will call back into registerFunction as needed.
	mod, err := runtime.Instantiate(ctx, plugin)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate the plugin: %v", err)
	}

	return mod, nil
}

func callFunction(ctx context.Context, mod wasm.Module, fn wasm.Function, info FuncInfo, inputs map[string]any) (any, error) {

	def := fn.Definition()
	paramTypes := def.ParamTypes()
	mem := mod.Memory()

	// Get parameters to pass as input to the plugin function
	// Note that we can't use def.ParamNames() because they are only available when the plugin
	// is compiled in debug mode. They're striped by optimization in release mode.
	// Instead, we can use the parameter names from registration.
	// Also note, that the order of registered params should match order of params in wasm.
	params := make([]uint64, len(paramTypes))
	for i, p := range info.Function.Parameters {
		val := inputs[p.Name]
		if val == nil {
			return nil, fmt.Errorf("parameter %s is missing", p.Name)
		}

		param, err := convertParam(ctx, mod, p.Type, paramTypes[i], val)
		if err != nil {
			return nil, fmt.Errorf("parameter %s is invalid: %v", p.Name, err)
		}

		params[i] = param
	}

	// Call the wasm function
	log.Printf("calling function \"%s\" for resolver \"%s\"\n", info.Function.Name, info.Function.Resolver)
	res, err := fn.Call(ctx, params...)
	if err != nil {
		return nil, err
	}

	// Get the result
	result, err := convertResult(mem, info.Function.ReturnType, def.ResultTypes()[0], res[0])
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
		fmt.Println("Failed to decode request body: ", err)
		return
	}

	info := functionsMap[req.Resolver]
	mod := *info.Module
	fnName := info.Function.Name

	fn := mod.ExportedFunction(fnName)
	if fn == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "The function \"%s\" was not found.", fnName)
		return
	}

	ctx := r.Context()

	if req.Args != nil {

		// Call the function, passing in the args from the request
		result, err := callFunction(ctx, mod, fn, info, req.Args)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("Failed to call function \"%s\": %v", fnName, err)
			return
		}

		// Handle no result due to void return type
		if result == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Write the result
		writeDataAsJson(w, result)

	} else if req.Parents != nil {

		results := make([]any, len(req.Parents))

		// Call the function for each parent
		for i, parent := range req.Parents {
			results[i], err = callFunction(ctx, mod, fn, info, parent)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Printf("Failed to call function \"%s\": %v", fnName, err)
				return
			}
		}

		// Write the result
		writeDataAsJson(w, results)

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Request must have either args or parents.")
	}
}

func writeDataAsJson(w http.ResponseWriter, data any) {
	output, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed to serialize result data: ", err)
		return
	}

	w.Write(output)
}

func convertParam(ctx context.Context, mod wasm.Module, hmType protos.HMType, wasmType wasm.ValueType, val any) (uint64, error) {
	switch hmType {

	case protos.HMType_INT:
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

	case protos.HMType_FLOAT:
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

	case protos.HMType_STRING:
		s, ok := val.(string)
		if !ok {
			return 0, fmt.Errorf("input value is not a string")
		}

		// Note, strings are passed as a pointer to a string in wasm memory
		if wasmType != wasm.ValueTypeI32 {
			return 0, fmt.Errorf("parameter is not defined as a string on the function")
		}

		buf := encodeUTF16(s)
		ptr := allocateWasmMemory(ctx, mod, len(buf), asString)
		mod.Memory().Write(ptr, buf)
		return uint64(ptr), nil

	case protos.HMType_BOOL:
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

	default:
		return 0, fmt.Errorf("unknown parameter type")
	}
}

func convertResult(mem wasm.Memory, hmType protos.HMType, wasmType wasm.ValueType, res uint64) (any, error) {

	// TODO: Do we need to handle unsigned ints differently?

	switch hmType {
	case protos.HMType_VOID:
		return nil, nil

	case protos.HMType_INT:
		switch wasmType {
		case wasm.ValueTypeI32:
			return wasm.DecodeI32(res), nil

		case wasm.ValueTypeI64:
			return int64(res), nil

		default:
			return nil, fmt.Errorf("return type is not defined as an int on the function")
		}

	case protos.HMType_FLOAT:

		switch wasmType {
		case wasm.ValueTypeF32:
			return wasm.DecodeF32(res), nil

		case wasm.ValueTypeF64:
			return wasm.DecodeF64(res), nil

		default:
			return nil, fmt.Errorf("return type is not defined as a float on the function")
		}

	case protos.HMType_STRING:
		// Note, strings are passed as a pointer to a string in wasm memory
		if wasmType != wasm.ValueTypeI32 {
			return nil, fmt.Errorf("return type is not defined as a string on the function")
		}

		buf := readBuffer(mem, uint32(res))
		return decodeUTF16(buf), nil

	case protos.HMType_BOOL:
		if wasmType != wasm.ValueTypeI32 {
			return nil, fmt.Errorf("return type is not defined as an bool on the function")
		}

		if res == 1 {
			return true, nil
		} else {
			return false, nil
		}

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

func readBuffer(mem wasm.Memory, offset uint32) []byte {

	// The length of AssemblyScript managed objects is stored 4 bytes before the offset.
	// See https://www.assemblyscript.org/runtime.html#memory-layout

	// Read the length.
	len, ok := mem.ReadUint32Le(offset - 4)
	if !ok {
		return nil
	}

	// Now read the data into the buffer.
	buf, ok := mem.Read(offset, len)
	if !ok {
		return nil
	}

	return buf
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
