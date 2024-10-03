# Modus

## TODO - rewrite this whole file in the context of Modus

This repository contains the source code for the _Hypermode Runtime_.
The Runtime loads and executes _plugins_ containing _Hypermode Functions_.

To get started with Hypermode Functions written in AssemblyScript, visit the
[`functions-as`](https://github.com/hypermodeAI/functions-as) repository.

## Command Line Parameters

When starting the Runtime, you can use the following command line parameters:

- `--port` - The HTTP port to listen on. Defaults to `8686`.
- `--modelHost` - The base DNS of the host endpoint to the model server.
- `--storagePath` - The path to a directory used for local storage.
  - Linux / OSX default: `$HOME/.hypermode`
  - Windows default: `%APPDATA%\Hypermode`
- `--useAwsSecrets` - Use AWS Secrets Manager for API keys and other secrets.
- `--useAwsStorage` - Use AWS S3 for storage instead of the local filesystem.
- `--s3bucket` - The S3 bucket to use, if using AWS storage.
- `--s3path` - The path within the S3 bucket to use, if using AWS storage.
- `--refresh` - The refresh interval to reload any changes. Defaults to `5s`.
- `--jsonlogs` - Use JSON format for logging.

_Note: You can use either `-` or `--` prefixes, and you use either a space or `=` to provide values._

## Environment Variables

The following environment variables are used by the Hypermode Runtime:

- `ENVIRONMENT` - The name of the environment, such as `dev`, `stage`, or `prod`. Defaults to `dev` when not specified.
- `NAMESPACE` - Used in non-dev (ie. stage/prod) environments to retrieve the Kubernetes namespace for the Hypermode backend.
- `HYPERMODE_METADATA_DB` - Used in dev environment only, to set a connection string to a local Postgres database to use for Hypermode metadata.

### Local Secrets

When running locally, in addition to the above, secrets are passed by environment variable using the following convention:

```
HYPERMODE_<HOST_NAME>_<SECRET_NAME>
```

- Substitute `<HOST_NAME>` for the upper-cased name of the model, replacing hyphens (`-`) with underscores (`_`).
- Substitute `<SECRET_NAME>` for the upper-cased name of the secret variable used in the manifest.

For example, if the manifest contains the host entry:

```json
"my-api": {
  "endpoint":"https://api.example.com/",
  "headers": {
    "Authorization": "Bearer {{AUTH_TOKEN}}"
  }
},
```

... then the corresponding environment variable name would be `HYPERMODE_MY_API_AUTH_TOKEN`.
The value of that environment variable would be used as the bearer token for any requests to that host.

### Using a `.env` file

While environment variables can be set through traditional mechanisms, they can also (optionally)
be set in a file called `.env`, placed in the plugins path. The file should be a simple list of
key/value pairs for each environment variable you want to set when the Runtime is loaded.

For example:

```
HYPERMODE_FOO=abc123
HYPERMODE_BAR=xyz456
```

## Building the Runtime

Ensure that you have Go installed in your dev environment.
The required minimum major version is specified in the [go.mod](./go.mod) file.
It is recommended to use the latest patch version available, and keep current.

Then, you can do any of the following:

- You can run directly from source code:

  ```sh
  go run .
  ```

- You can compile the source code and run the output:

  ```sh
  go build
  ./hypruntime
  ```

- You can run and debug the source code in VS Code, using the VS Code debugger.

## Docker Setup

To build a Docker image for the Hypermode Runtime:

```sh
docker build -t hypermode/runtime .
```

When running the image via Docker Desktop, keep in mind:

- You may wish to give the container a specific name using `--name`.
- Port `8686` should be exposed.
- If using local storage, you'll need to map the container's `/root/.hypermode` folder to your own `~/.hypermode` folder.
- You may need to pass command line parameters and/or set environment variables.

For example:

```sh
docker run --name hypruntime \
  -p 8686:8686 \
  -v ~/.hypermode:/root/.hypermode hypermode/runtime
```

_Note, if you have previously created a container with the same name, then delete it first with `docker rm hypruntime`._

## AWS Setup

If configured to do so, the Hypermode Runtime may access AWS resources.
If you are debugging locally, set up an AWS profile.

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

After SSO login you can start the Runtime, either from the VS Code debugger
using the `Hypermode Runtime (AWS)` launch profile, or from the command line as follows:

```sh
export AWS_SDK_LOAD_CONFIG=true
export AWS_PROFILE=sandbox
./hypruntime \
  --useAwsSecrets \
  --useAwsStorage \
  --s3bucket=sandbox-runtime-storage \
  --s3path=shared
```

You can omit the exports if the environment variables are already set.
You can also use any S3 bucket or path you like. If a path is not specified, the Runtime will look for files in the root of the bucket.

_The shared sandbox is intended for temporary use. In production, each customer's backend gets a separate path within a single bucket._

## Using Model Inference History Locally

To set up the model inference history table, run the following command, which will start a local postgres instance.

```sh
cd tools/local
docker-compose up
```

Next, you will need to apply the database schema using the [golang-migrate](https://github.com/golang-migrate/migrate) utility.

On MacOS, you can install this utility with the following:

```sh
brew install golang-migrate
```

Then, you can apply the migration as follows:

```sh
export POSTGRESQL_URL='postgresql://postgres:postgres@localhost:5433/my-runtime-db?sslmode=disable'
migrate -database ${POSTGRESQL_URL} -path db/migrations up
```

Now any model inference will be logged in the `inferences` table in the `my-runtime-db` database.

## Unit Testing

Unit tests are created using Go's [built-in unit test support](https://go.dev/doc/tutorial/add-a-test).

To run all tests in the project:

```sh
go test ./...
```

Or, you can just run tests in specific folders. For example:

```sh
go test ./functions
```

Tests can also be run from VS Code's Testing panel, and are run automatically for pull requests.
