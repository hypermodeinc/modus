# Hypermode Example Plugin 1

This is an example Hypermode plugin, written in [AssemblyScript](https://www.assemblyscript.org/).

## Dependencies

First, install [Node.js](https://nodejs.org/) version 18 or newer
on your development workstation or build server

Then, install package dependencies:

```sh
npm install
```

## Building

To build the plugin:

```sh
npm run build
```

This will create both a `debug` and `release` build of the app.
Output files are located in the `build` folder.

Currently, only the `release.wasm` file is used.
In the future, the other files will help support testing and debugging.

## Schema and Sample Data

The `loaddata.sh` script can be used to populate schema and sample data.
It connects to Dgraph on `localhost:8080`, and applies the `schema.graphql` and `sampledata.graphql` files.
