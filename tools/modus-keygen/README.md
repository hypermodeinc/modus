# Keygen Tool

This is a small utility that will generate keys and tokens for use with Modus.

## Usage

You must have Go v1.25 or newer installed locally. We currently do not release this tool as a
binary.

You can run it in one of two ways:

- Clone this repo, then from the `/tools/modus-keygen` directory, run the utility via: `go run .`
- Install it with `go install github.com/hypermodeinc/modus/tools/modus-keygen@latest` then run
  `modus-keygen`

The keygen tool will print the following items to stdout:

- A JSON encoded object containing an RSA public key, suitable for passing to `MODUS_PEMS`
- A JWT token that is signed with that key, having a 1 year expiration date and no other claims,
  suitable for passing as a bearer token for authorization

The tool will also create two files in the current working directory:

- `private-key.pem` - A 2048 bit RSA private key
- `public-key.pem` - The corresponding public key

If the private key file already exists in the working directory, the tool will re-use it instead of
creating a new one.

**IMPORTANT** - Keep the private key secret, and do not lose it! You will need it to issue new JWT
tokens that can be validated with the public key.

The public key file contains the same key that was given as JSON. It is there for convenience only.

You can run the tool multiple times to generate additional JWTs from the same key pair.
