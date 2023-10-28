# Hypermode Example Plugin 1

This is an example Hypermode plugin, written in [AssemblyScript](https://www.assemblyscript.org/).

## Dependencies

The following must be installed on your development workstation or build server:

- [Node.js](https://nodejs.org/) version 18 or newer
- The [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation/) (`protoc`)

## Building

_NOTE: The `hypermode-as` library must be built first, or you will get an error when building the plugin._

To build the plugin, run the `build.sh` script.

This will create both a `debug` and `release` build of the app.
Currently, only the `release.wasm` file is used.
In the future, the other files will help support testing and debugging.
