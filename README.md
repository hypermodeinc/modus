# Hypermode Runtime

This repository contains the source code for the _Hypermode Runtime_.
The Runtime loads and executes _plugins_ containing _Hypermode Functions_.

To get started with Hypermode Functions written in AssemblyScript, visit the
[`functions-as`](https://github.com/gohypermode/functions-as) repository.

## Command Line Parameters

When starting the Runtime, you can use the following command line parameters:

- `--port` - The HTTP port to listen on.  Defaults to `8686`.
- `--modelHost` - The base DNS of the host endpoint to the model server.
- `--storagePath` - The path to a directory used for local storage.
  - Linux / OSX default: `$HOME/.hypermode`
  - Windows default: `%APPDATA%\Hypermode`
- `--useAwsSecrets` - Use AWS Secrets Manager for API keys and other secrets.
- `--useAwsStorage` - Use AWS S3 for storage instead of the local filesystem.
- `--s3bucket` - The S3 bucket to use, if using AWS storage.
- `--s3path` - The path within the S3 bucket to use, if using AWS storage.
- `--refresh` - The refresh interval to reload any changes.  Defaults to `5s`.
- `--jsonlogs` - Use JSON format for logging.

_Note: You can use either `-` or `--` prefixes, and you use either a space or `=` to provide values._

## Environment Variables

Some model hosts have pre-established well-known environment variable names.
The Hypermode Runtime will use them if they are available:

- `OPENAI_API_KEY` - If set, will be used by an model who's host is `"openai"`.

Alternatively, any model can have its key passed as an environment variable,
using the following convention:

- `HYP_MODEL_KEY_<MODEL_NAME>` - Can be used to specify keys for any model.
  - _Substitute `<MODEL_NAME>` for the upper-cased name of the model._
  - _Replace any hyphen (`-`) characters with underscore (`_`) characters._

### Using a `.env` file

While environment variables can be set through traditional mechanisms, they can also (optionally)
be set in a file called `.env`, placed in the plugins path.  The file should be a simple list of
key/value pairs.  For example:

```
HYP_MODEL_KEY_FOO=abc123
HYP_MODEL_KEY_BAR=xyz456
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
  ./hmruntime
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
docker run --name hmruntime \
  -p 8686:8686 \
  -v ~/.hypermode:/root/.hypermode hypermode/runtime
```

_Note, if you have previously created a container with the same name, then delete it first with `docker rm hmruntime`._

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
./hmruntime \
  --useAwsSecrets \
  --useAwsStorage \
  --s3bucket=sandbox-runtime-storage \
  --s3path=shared
```

You can omit the exports if the environment variables are already set.
You can also use any S3 bucket or path you like.  If a path is not specified, the Runtime will look for files in the root of the bucket.

_The shared sandbox is intended for temporary use.  In production, each customer's backend gets a separate path within a single bucket._

## Using Model Inference History Locally
To set up the model inference history table, run the following commands
```sh
cd tools/local && docker-compose up
```
This will start a local postgres instance. To add the inference history table and the schema, from a separate terminal in the root runtime directory, run the following command
```sh
export POSTGRESQL_URL='postgresql://postgres:postgres@localhost:5433/my-runtime-db?sslmode=disable' && migrate -database ${POSTGRESQL_URL} -path db/migrations up
```

Now any model inference will be logged in the `local_instance` table in the `my-runtime-db` database.

#### Troubleshooting

**error posting to model endpoint: error getting model key: error getting secret: NoCredentialProviders: no valid providers in chain. Deprecated.**

Your are missing `export AWS_SDK_LOAD_CONFIG=true`

or your AWS region is wrong, or you do not have an AWS secret set for the ModelSpec name. Add via:

`aws secretsmanager create-secret --name '<ModelSpec.name>' --secret-string '<apikey>'
`
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
