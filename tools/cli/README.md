# @hypermode/cli

The Modus CLI

[![oclif](https://img.shields.io/badge/cli-oclif-brightgreen.svg)](https://oclif.io)
[![Version](https://img.shields.io/npm/v/@hypermode/cli.svg)](https://npmjs.org/package/@hypermode/cli)
[![Downloads/week](https://img.shields.io/npm/dw/@hypermode/cli.svg)](https://npmjs.org/package/@hypermode/cli)

<!-- toc -->
* [@hypermode/cli](#hypermodecli)
* [Usage](#usage)
* [Commands](#commands)
<!-- tocstop -->

# Usage

<!-- usage -->
```sh-session
$ npm install -g @modus/cli
$ modus COMMAND
running command...
$ modus (--version)
@modus/cli/0.0.0 linux-x64 node-v22.8.0
$ modus --help [COMMAND]
USAGE
  $ modus COMMAND
...
```
<!-- usagestop -->

# Commands

<!-- commands -->
* [`modus autocomplete [SHELL]`](#modus-autocomplete-shell)
* [`modus build [PATH]`](#modus-build-path)
* [`modus deploy`](#modus-deploy)
* [`modus new`](#modus-new)
* [`modus run [PATH]`](#modus-run-path)
* [`modus sdk install [VERSION]`](#modus-sdk-install-version)
* [`modus sdk list`](#modus-sdk-list)
* [`modus sdk remove install [VERSION]`](#modus-sdk-remove-install-version)
* [`modus upgrade`](#modus-upgrade)

## `modus autocomplete [SHELL]`

Display autocomplete installation instructions.

```
[1mUsage:[22m [1m[94mmodus[39m[22m autocomplete
```

_See code: [@oclif/plugin-autocomplete](https://github.com/oclif/plugin-autocomplete/blob/v3.2.5/src/commands/autocomplete/index.ts)_

## `modus build [PATH]`

Build a Modus project

```
[1mUsage:[22m [1m[94mmodus[39m[22m build
```

_See code: [src/commands/build/index.ts](https://github.com/HypermodeAI/modus/blob/v0.0.0/src/commands/build/index.ts)_

## `modus deploy`

Deploy a Modus app to Hypermode

```
[1mUsage:[22m [1m[94mmodus[39m[22m deploy
```

_See code: [src/commands/deploy/index.ts](https://github.com/HypermodeAI/modus/blob/v0.0.0/src/commands/deploy/index.ts)_

## `modus new`

Create a new Modus project

```
[1mUsage:[22m [1m[94mmodus[39m[22m new
```

_See code: [src/commands/new/index.ts](https://github.com/HypermodeAI/modus/blob/v0.0.0/src/commands/new/index.ts)_

## `modus run [PATH]`

Run a Modus app locally

```
[1mUsage:[22m [1m[94mmodus[39m[22m run
```

_See code: [src/commands/run/index.ts](https://github.com/HypermodeAI/modus/blob/v0.0.0/src/commands/run/index.ts)_

## `modus sdk install [VERSION]`

Install a specific SDK version

```
[1mUsage:[22m [1m[94mmodus[39m[22m sdk:install
```

_See code: [src/commands/sdk/install/index.ts](https://github.com/HypermodeAI/modus/blob/v0.0.0/src/commands/sdk/install/index.ts)_

## `modus sdk list`

List installed SDK versions

```
[1mUsage:[22m [1m[94mmodus[39m[22m sdk:list
```

_See code: [src/commands/sdk/list/index.ts](https://github.com/HypermodeAI/modus/blob/v0.0.0/src/commands/sdk/list/index.ts)_

## `modus sdk remove install [VERSION]`

Uninstall a specific SDK version

```
[1mUsage:[22m [1m[94mmodus[39m[22m sdk:remove:install
```

_See code: [src/commands/sdk/remove/install.ts](https://github.com/HypermodeAI/modus/blob/v0.0.0/src/commands/sdk/remove/install.ts)_

## `modus upgrade`

Upgrade a Modus component

```
[1mUsage:[22m [1m[94mmodus[39m[22m upgrade
```

_See code: [src/commands/upgrade/index.ts](https://github.com/HypermodeAI/modus/blob/v0.0.0/src/commands/upgrade/index.ts)_
<!-- commandsstop -->
