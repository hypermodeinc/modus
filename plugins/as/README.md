# Hypermode AssemblyScript Plugins

This folder contains multiple projects:

- [`hypermode-as`](./hypermode-as) is a library that all plugins are required to import,
  which provides support for working with the Hypermode runtime.

- [`hmplugin1`](./hmplugin1) is an example Hypermode plugin.

## Dependencies

The following must be installed on your development workstation or build server:

- [Node.js](https://nodejs.org/) version 18 or newer
- The [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation/) (`protoc`)

## Building

Each project has its own `build.sh` script, or you can use the `buildall.sh` script
in this folder to build everything.  If building separately, note that the `hypermode-as`
library must be built _first_.
