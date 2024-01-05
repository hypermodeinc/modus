# Hypermode Runtime

This folder contains the Hypermode Runtime, which is the backend server that manages
Hypermode Plugins and executes their Hypermode Functions.

## Dependencies

The following must be installed on your development workstation or build server:

- A [Go](https://go.dev/) compiler, of at least the version specified in [`go.mod`](./go.mod)
- The [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation/) (`protoc`)

## Building

To build the Hypermode runtime server: `go build`

To build the docker image, from the root directory: `docker build -t hypermode/runtime .`

## Running

- To run the compiled program, invoke the `hmruntime` binary.
- To run from code (while developing), use `go run .` instead.
- To run using docker containers, first decide where plugins will be loaded from.
  - To just use the default `hmplugin1` sample plugin, use `docker run -p 8686:8686 hypermode/runtime --dgraph=http://host.docker.internal:8080`.
  - Or, mount a plugins directory on the host, use `docker run -p 8686:8686 -v <PLUGINS_PATH>:/plugins hypermode/runtime --dgraph=http://host.docker.internal:8080`.
    - For example `-v ./plugins/as:/plugins`
