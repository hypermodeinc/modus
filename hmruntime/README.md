# Hypermode Runtime

This folder contains the Hypermode Runtime, which is the backend server that manages
Hypermode Plugins and executes their Hypermode Functions.

## Dependencies

The following must be installed on your development workstation or build server:

- A [Go](https://go.dev/) compiler, of at least the version specified in [`go.mod`](./go.mod)
- The [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation/) (`protoc`)

## Building

To build the server, run the `build.sh` script.
This will generate the required protobuf stubs, then compile the program.

## Running

- To run the compiled program, invoke the `hmruntime` binary.
- To run from code (while developing), use `go run main.go` instead.

## Notes

Currently, the `hmplugin1` plugin is hardcoded, so be sure to compile it before
running this.  In the future, plugins will be developed independently of the runtime,
and loaded from a database or repository.
