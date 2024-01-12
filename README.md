# Hypermode Runtime

This repository contains the source code for the _Hypermode Runtime_.

The runtime loads and executes _plugins_ containing _Hypermode Functions_.

To get started with Hypermode Functions written in AssemblyScript, visit the
[`hypermode-as`](https://github.com/gohypermode/hypermode-as) repository.

## Docker Setup

To build a Docker image for the Hypermode Runtime:

```
docker build -t hypermode/runtime .
```

Then you can run that image.  Port `8686` should be exposed.

```
docker run -p 8686:8686 -v <PLUGINS_PATH>:/plugins hypermode/runtime --dgraph=<DGRAPH_ALPHA_URL>
```

Replace the following:
- `<PLUGINS_PATH>` should be the local path to the folder where you will load plugins from.
  - You can use paths such as `./plugins` or `~/plugins` etc. depending on where you want to keep your plugin files.
- `<DGRAPH_ALPHA_URL>` should be the URL to the Dgraph Alpha endpoint you are connecting the runtime to.
  - To connect to Dgraph running in another docker container, use `host.docker.internal`.

Optionally, you may also wish to give the container a specific name using the `--name` flag.
For example, to start a new Docker container named `hmruntime`, looking for plugins in a local `./plugins` folder,
and connecting to a local Dgraph docker image:

```
docker run --name hmruntime -p 8686:8686 -v ./plugins:/plugins hypermode/runtime --dgraph=http://host.docker.internal:8080
```

_Note, if you have previously created a container with the same name, then delete it first with `docker rm hmruntime`._

## Building without Docker

If needed, you can compile and run the Hypermode Runtime without using Docker.
This is most common for local development.

Be sure that you have Go installed in your dev environment, at the version specified in the [.go-version](./go-verson) file, or higher.
Then you can either run the Runtime code directly from source:

```
go run .
```

Or, you can build the `hmruntime` executable and then run that:

```
go build
./hmruntime
```

### Command Line Arguments

When starting the runtime, you may sometimes need to use the following command line arguments:

- `--port` - The port that the runtime will listen for HTTP requests on.  Defaults to `8686`.
- `--dgraph` - The URL to the Dgraph Alpha endpoint.  Defaults to `http://localhost:8080`.
- `--plugins` (or `--plugin`) - The folder that the runtime will look for plugins in.  Defaults to `./plugins`.

_Note: You can use either `-` or `--` prefixes, and you can add parameters with either a space or `=`._

## Working locally with plugins

Regardless of whether you use Docker or not, it is often useful to be developing both the runtime
and a plugin at the same time.  This is especially true if you are developing a new host function
for the runtime, and need to expose it via the `hypermode-as` library.

To facilitate this, you can point the runtime's plugins path to the root folder of any plugin's
source code.  The runtime will use the `build/debug.wasm` file, and will pick up changes
automatically when rebuilding the plugin.

For example, you may have the `runtime` and `hypermode-as` repos in the same parent directory,
and are working on a plugin in the `examples` folder, such as `hmplugin1`.  You can start the
runtime like so:

```
go run . --plugin ../hypermode-as/examples/hmplugin1
```

Or, if you are working on more than one plugin simultaneously you can use their parent directory:

```
go run . --plugins ../hypermode-as/examples
```

However, be aware that if there are conflicts between function names in the plugins,
the last one loaded byt the runtime will take precedence.  Thus, it's usually better to work
on one plugin at a time.

## Dgraph Setup

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
