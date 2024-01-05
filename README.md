# Hypermode Runtime and Plugins

This repository contains the server-side code for the Hypermode Runtime.

It also currently contains the client-side code for Hypermode Plugins,
but that will eventually be moved to a separate repository as the project matures.

## Projects

See the separate `README.md` files in each project for further details:

- [`hmruntime`](./hmruntime) - The Hypermode Runtime service, written in Go
- [`plugins/as`](./plugins/as) - Plugin library and example, written in AssemblyScript

## Getting Started

### Dgraph Setup

Currently, the Hypermode Runtime service emulates parts of the 
[Dgraph Lambda](https://dgraph.io/docs/graphql/lambda/lambda-overview/) protocol.
Like Lambda, it listens for HTTP on port `8686` on the `/graphql-worker` endpoint.
Thus, you can tell Dgraph to use it when starting Dgraph Alpha:

```
dgraph zero
dgraph alpha --graphql lambda-url=http://localhost:8686/graphql-worker
```

If you like, you can use the [Dgraph standalone Docker image](https://dgraph.io/docs/deploy/installation/single-host-setup/).
Just add the Lambda URL as an environment variable, using `host.docker.internal` to escape the container.

For example:

```
docker run --name <CONTAINER_NAME> \
  -d -p 8080:8080 -p 9080:9080 \
  -v <DGRAPH_DATA_PATH>:/dgraph \
  --env=DGRAPH_ALPHA_GRAPHQL=lambda-url=http://host.docker.internal:8686/graphql-worker \
  dgraph/standalone:latest
```

### Hypermode Plugins

First, ensure you have [Node.js](https://nodejs.org/) 18 or higher installed.


Recommended: Install [nvm](https://github.com/nvm-sh/nvm/blob/master/README.md) and use it to install Node.js 20:

```
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash
nvm install 20
nvm use 20
```

Then, compile the `hmplugin1` example plugin, by running the following:

```
cd plugins/as/hmplugin1
npm install
npm run build
```

The `build` subfolder will contain the compiled plugin.

### Load schema and data

Run this script:

```sh
./plugins/as/hmplugin1/loaddata.sh
```

It connects to Dgraph on `localhost:8080`, and applies the `schema.graphql` and `sampledata.graphql` files.

### Hypermode Runtime

Now build and run the Hypermode Runtime:

```
cd hmruntime
go build
./hmruntime
```

Alternatively, you can run the Go code directly:

```
cd hmruntime
go run .
```

### Run the example

Try some graphql queries on `http://localhost:8080/graphql`:

```graphql
{
  add(a: 123, b: 456)
}
```

```graphql
{
  getFullName(firstName: "John", lastName:"Doe")
}
```

These will invoke the respective Hypermode functions within `hmplugin1`.

Next, try adding some data:

```graphql

mutation {
  addPerson(input: [
    { firstName: "Harry", lastName: "Potter" },
    { firstName: "Tom", lastName: "Riddle" },
    { firstName: "Albus", lastName: "Dumbledore" }
    ]) {
    person {
      id
      firstName
      lastName
      fullName
    }
  }
}
```

In the response, notice how the `fullName` field is returned,
which is the output from calling the `getFullName` function in `hmplugin1`.

You can now also query for data:

```graphql
{
  queryPerson {
    id
    firstName
    lastName
    fullName
  }
}
```

Again, the `fullName` field is populated by calling `getFullName` in `hmplugin1`.
