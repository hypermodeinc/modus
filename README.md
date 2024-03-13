# Hypermode Runtime

This repository contains the source code for the _Hypermode Runtime_.

The runtime loads and executes _plugins_ containing _Hypermode Functions_.

To get started with Hypermode Functions written in AssemblyScript, visit the
[`hypermode-as`](https://github.com/gohypermode/hypermode-as) repository.

## Docker Setup

To build a Docker image for the Hypermode Runtime:

```sh
docker build -t hypermode/runtime .
```

Then you can run that image.  Port `8686` should be exposed.

```sh
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

```sh
docker run --name hmruntime -p 8686:8686 -v ./plugins:/plugins hypermode/runtime --plugins=/plugins --dgraph=http://host.docker.internal:8080
```

_Note, if you have previously created a container with the same name, then delete it first with `docker rm hmruntime`._

## Building without Docker

If needed, you can compile and run the Hypermode Runtime without using Docker.
This is most common for local development.

Be sure that you have Go installed in your dev environment, at the version specified in the [.go-version](./go-verson) file, or higher.

You can run the code directly using VSCode's debugger.

Alternatively you can either run the Runtime code directly from source:

```sh
go run . --plugin ../hypermode-as/examples/hmplugin1
```

Or, you can build the `hmruntime` executable and then run that:

```sh
go build
./hmruntime --plugin ../hypermode-as/examples/hmplugin1
```

## Command Line Arguments

When starting the runtime, you may need to use the following command line arguments:

- `--plugins` (or `--plugin`) - The folder that the runtime will look for plugins in.  ***Required.***
- `--port` - The port that the runtime will listen for HTTP requests on.  Defaults to `8686`.
- `--dgraph` - The URL to the Dgraph Alpha endpoint.  Defaults to `http://localhost:8080`.
- `--modelHost` - The host portion of the url to the model endpoint. This is used for cloud deployments for kserve hosted models .
- `--noreload` - Disables automatic reloading of plugins.
- `--s3bucket` - The S3 bucket to use, if using AWS for plugin storage.
- `--useAwsSecrets` - Directs the Runtime to use AWS Secret Manager for secrets such as model keys.
- `--refresh` - The refresh interval to check for plugins and schema changes.  Defaults to `5s`.
- `--jsonlogs` - Switches log output to JSON format.

_Note: You can use either `-` or `--` prefixes, and you can add parameters with either a space or `=`._

## Environment Variables

Some model hosts have pre-established well-known environment variable names.
The Hypermode Runtime will use them if they are available:

- `OPENAI_API_KEY` - If set, will be used by an model who's host is `"openai"`.

Alternatively, any model can have its key passed as an environment variable,
using the following convention:

- `HYP_MODEL_KEY_<MODEL_NAME>` - Can be used to specify keys for any model.
  - _Substitute `<MODEL_NAME>` for the upper-cased name of the model._

### Using a `.env` file

While environment variables can be set through traditional mechanisms, they can also (optionally)
be set in a file called `.env`, placed in the plugins path.  The file should be a simple list of
key/value pairs.  For example:

```
HYP_MODEL_KEY_FOO=abc123
HYP_MODEL_KEY_BAR=xyz456
```

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

```sh
go run . --plugin ../hypermode-as/examples/hmplugin1
```

Or, if you are working on more than one plugin simultaneously you can use their parent directory:

```sh
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

```sh
dgraph zero
dgraph alpha --graphql lambda-url=http://localhost:8686/graphql-worker
```

If you like, you can use the [Dgraph standalone Docker image](https://dgraph.io/docs/deploy/installation/single-host-setup/).
Just add the Lambda URL as an environment variable, using `host.docker.internal` to escape the container.

For example:

```sh
docker run --name <CONTAINER_NAME> \
  -d -p 8080:8080 -p 9080:9080 \
  -v <DGRAPH_DATA_PATH>:/dgraph \
  --env=DGRAPH_ALPHA_GRAPHQL=lambda-url=http://host.docker.internal:8686/graphql-worker \
  dgraph/standalone:latest
```

## AWS Setup
Runtime may access AWS resources, so we need to use an AWS profile.

```sh
export AWS_PROFILE=sandbox
```

Declare the `sandbox` profile in `~/.aws/config` as follows:

```ini
[profile sandbox]
sso_start_url = https://d-92674fec42.awsapps.com/start#
sso_region = us-west-2
sso_account_id = 436061841671
sso_role_name = AWSPowerUserAccess
region = us-west-2
```

Then run `aws sso login --profile sandbox` to login to the profile.

After SSO login you can start the runtime, either from the VS Code "Run and Debug" panel,
or from the command line as follows:

```sh
export AWS_SDK_LOAD_CONFIG=true
export AWS_PROFILE=sandbox
./hmruntime --plugins <your plugin folder>
```

_You can omit the exports if the environment variables are already set._

#### Troubleshooting

**error posting to model endpoint: error getting model key: error getting secret: NoCredentialProviders: no valid providers in chain. Deprecated.**

Your are missing `export AWS_SDK_LOAD_CONFIG=true`

or your AWS region is wrong, or you do not have an AWS secret set for the ModelSpec name. Add via:

`aws secretsmanager create-secret --name '<ModelSpec.name>' --secret-string '<apikey>'
`
### Using S3 for plugin storage

You can optionally use S3 for plugin storage.  This configuration is usually for staging or production,
but you can use it locally as well.

First, configure the AWS setup as described above.  Then start the runtime with the following command line arguments:

- `--s3bucket <bucket>` - An standard S3 bucket within the pre-configured AWS account. 
- `--plugins <folder>` - A folder within that bucket where `.wasm` files are contained.

Note that the `--plugins` argument is re-purposed.  It now refers to a folder inside the S3 bucket, 
rather than a local directory on disk.

For example:

```sh
export AWS_PROFILE=sandbox; ./hmruntime --s3bucket sandbox-runtime-plugins --plugins shared-plugins
```
You can omit the export if the environment variable is already set.

_In staging and production, a single S3 bucket is shared, but each backend has its own plugins folder._

#### Using S3 from VS Code Debugging Session

From the VS Code "Run and Debug" pane, you can choose the launch profile called `runtime (AWS S3 plugin storage)`.

When running using this profile, you will be prompted for the S3 bucket and folder name.
You can use the default values provided, or override them to load plugins from another location.

### Unit Testing

Unit tests are created using Go's [built-in unit test support](https://go.dev/doc/tutorial/add-a-test).

To run all tests in the project:

```sh
go test ./...
```

Or, you can just run tests in specific folders.  For example:

```sh
go test ./functions
```

Tests can also be run from VS Code's Testing panel, and are run automatically for pull requests.
