@hypermode/cli
=================

The Hypermode CLI


[![oclif](https://img.shields.io/badge/cli-oclif-brightgreen.svg)](https://oclif.io)
[![Version](https://img.shields.io/npm/v/@hypermode/cli.svg)](https://npmjs.org/package/@hypermode/cli)
[![Downloads/week](https://img.shields.io/npm/dw/@hypermode/cli.svg)](https://npmjs.org/package/@hypermode/cli)


<!-- toc -->
* [Usage](#usage)
* [Commands](#commands)
<!-- tocstop -->
# Usage
<!-- usage -->
```sh-session
$ npm install -g @hypermode/cli
$ hyp COMMAND
running command...
$ hyp (--version)
@hypermode/cli/0.0.0 linux-x64 node-v22.4.0
$ hyp --help [COMMAND]
USAGE
  $ hyp COMMAND
...
```
<!-- usagestop -->
# Commands
<!-- commands -->
* [`hyp autocomplete [SHELL]`](#hyp-autocomplete-shell)
* [`hyp build [PATH]`](#hyp-build-path)
* [`hyp deploy`](#hyp-deploy)
* [`hyp help PERSON`](#hyp-help-person)
* [`hyp new`](#hyp-new)
* [`hyp run`](#hyp-run)
* [`hyp uninstall`](#hyp-uninstall)
* [`hyp upgrade`](#hyp-upgrade)

## `hyp autocomplete [SHELL]`

Display autocomplete installation instructions.

```
USAGE
  $ hyp autocomplete [SHELL] [-r]

ARGUMENTS
  SHELL  (zsh|bash|powershell) Shell type

FLAGS
  -r, --refresh-cache  Refresh cache (ignores displaying instructions)

DESCRIPTION
  Display autocomplete installation instructions.

EXAMPLES
  $ hyp autocomplete

  $ hyp autocomplete bash

  $ hyp autocomplete zsh

  $ hyp autocomplete powershell

  $ hyp autocomplete --refresh-cache
```

_See code: [@oclif/plugin-autocomplete](https://github.com/oclif/plugin-autocomplete/blob/v3.2.4/src/commands/autocomplete/index.ts)_

## `hyp build [PATH]`

Build a Hypermode project

```
USAGE
  $ hyp build [PATH]

ARGUMENTS
  PATH  ./my-project-|-Directory to build

DESCRIPTION
  Build a Hypermode project

EXAMPLES
  $ hyp build ./my-project
```

_See code: [src/commands/build/index.ts](https://github.com/HypermodeAI/cli/blob/v0.0.0/src/commands/build/index.ts)_

## `hyp deploy`

Deploy a Hypermode app to GitHub

```
USAGE
  $ hyp deploy -f <value>

FLAGS
  -f, --from=<value>  (required) Who is saying hello

DESCRIPTION
  Deploy a Hypermode app to GitHub

EXAMPLES
  $ hyp deploy friend --from oclif
  hello friend from oclif! (./src/commands/hello/index.ts)
```

_See code: [src/commands/deploy/index.ts](https://github.com/HypermodeAI/cli/blob/v0.0.0/src/commands/deploy/index.ts)_

## `hyp help PERSON`

Print help for certain command

```
USAGE
  $ hyp help PERSON -f <value>

ARGUMENTS
  PERSON  <command>-|-Print help for command

FLAGS
  -f, --from=<value>  (required) Who is saying hello

DESCRIPTION
  Print help for certain command

EXAMPLES
  $ hyp help friend --from oclif
  hello friend from oclif! (./src/commands/hello/index.ts)
```

_See code: [src/commands/help/index.ts](https://github.com/HypermodeAI/cli/blob/v0.0.0/src/commands/help/index.ts)_

## `hyp new`

Create a new Hypermode project

```
USAGE
  $ hyp new [--name <value>] [--pkgmgr <value>] [--sdk <value>]

FLAGS
  --name=<value>    Project name
  --pkgmgr=<value>  Package manager to use
  --sdk=<value>     SDK to use

DESCRIPTION
  Create a new Hypermode project

EXAMPLES
  $ hyp new --sdk go
```

_See code: [src/commands/new/index.ts](https://github.com/HypermodeAI/cli/blob/v0.0.0/src/commands/new/index.ts)_

## `hyp run`

Run a Hypermode app locally

```
USAGE
  $ hyp run -f <value>

FLAGS
  -f, --from=<value>  (required) Who is saying hello

DESCRIPTION
  Run a Hypermode app locally

EXAMPLES
  $ hyp run friend --from oclif
  hello friend from oclif! (./src/commands/hello/index.ts)
```

_See code: [src/commands/run/index.ts](https://github.com/HypermodeAI/cli/blob/v0.0.0/src/commands/run/index.ts)_

## `hyp uninstall`

Remove Hypermode Runtime and SDK

```
USAGE
  $ hyp uninstall

DESCRIPTION
  Remove Hypermode Runtime and SDK
```

_See code: [src/commands/uninstall/index.ts](https://github.com/HypermodeAI/cli/blob/v0.0.0/src/commands/uninstall/index.ts)_

## `hyp upgrade`

Upgrade a Hypermode component

```
USAGE
  $ hyp upgrade -f <value>

FLAGS
  -f, --from=<value>  (required) Who is saying hello

DESCRIPTION
  Upgrade a Hypermode component

EXAMPLES
  $ hyp upgrade friend --from oclif
  hello friend from oclif! (./src/commands/hello/index.ts)
```

_See code: [src/commands/upgrade/index.ts](https://github.com/HypermodeAI/cli/blob/v0.0.0/src/commands/upgrade/index.ts)_
<!-- commandsstop -->
