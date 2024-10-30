<!-- markdownlint-disable first-line-heading -->
<div align="center">

  [![modus](https://github.com/user-attachments/assets/a0652fa6-6836-49f6-9c8c-c007cafab000)](https://github.com/hypermodeinc/modus)

  [![GitHub License](https://img.shields.io/github/license/hypermodeinc/modus)](https://github.com/hypermodeinc/modus?tab=Apache-2.0-1-ov-file#readme)
  [![chat](https://img.shields.io/discord/1267579648657850441)](https://discord.gg/NJQ4bJpffF)
  [![GitHub Repo stars](https://img.shields.io/github/stars/hypermodeinc/modus)](https://github.com/hypermodeinc/modus/stargazers)
  [![GitHub commit activity](https://img.shields.io/github/commit-activity/m/hypermodeinc/modus)](https://github.com/hypermodeinc/modus/commits/main/)
  [![GitHub commits since latest release](https://img.shields.io/github/commits-since/hypermodeinc/modus/latest)](https://github.com/hypermodeinc/modus/commits/main/)

</div>

<p align="center">
   <a href="https://docs.hypermode.com/modus/quickstart">Get Started</a> · 
   <a href="https://docs.hypermode.com/">Docs</a> · 
   <a href="https://discord.com/invite/MAZgkhP6C6">Discord</a>
<p>


**Modus is an open-source, serverless framework for building intelligent functions and APIs, powered by WebAssembly.**

This repository contains the Modus source code, including the Modus runtime, SDKs, CLI, and introductory examples. It simplifies the integration of AI models, data, and business logic with sandboxed execution, providing a code-first workbench for building AI-enabled apps. 

It's really fast.

> ❕ Note: prior to version 0.13, Modus was simply referred to as "Hypermode". Please bear with us if you find name conflicts while we are in transition.

## Quick Install

Install the Modus CLI

```bash
npm install -g @hypermode/modus-cli
```

Initialize your Modus app

```bash
modus new
```

Run your app locally with fast refresh

```bash
modus dev
```

## Demo

<div align="center">

[![Getting started with Modus video](https://img.youtube.com/vi/3CcJTXTmz88/0.jpg)](https://www.youtube.com/watch?v=3CcJTXTmz88)

</div>

## What's it good for?

We designed Modus primarily as a general-purpose framework, it just happens to treat models as a first-class component. With Modus you can use models, as appropriate, without additional complexity.

Modus is optimized for applications that require sub-second response times. We’ve made trade-offs to prioritize speed and simplicity.

## Programming Languages

Since Modus is based on WebAssembly, you can write Modus apps in various programming languages.
Each language offers the full capabilities of the Modus framework.

Currently, the supported languages you may choose from are:

- [AssemblyScript](https://www.assemblyscript.org/) - A TypeScript-like language designed for WebAssembly.
  - If you are primarily used to writing front-end web apps, you'll feel at home with AssemblyScript.

- [Go](https://go.dev/) - A general-purpose programming language originally designed by Google.
  - If you are primarily used to writing back-end apps, you'll likely prefer to use Go.

Additional programming languages may be supported in the future.

## Hosting

We have designed Hypermode to be the best place to run your Modus app.
Hypermode hosting plans include features you might expect, such as support, telemetry, and high availability.
They also include specialty features such as model hosting that are purposefully designed to work in tandem with Modus apps.

As Modus is a free, open-source framework, you’re welcome to run your Modus apps on your own hardware or on any
hosting platform that meets your needs.

## Open Source

Modus is developed by [Hypermode](https://hypermode.com/) as an open-source project, integral but independent from Hypermode.

We welcome external contributions. See the [CONTRIBUTING.md](./CONTRIBUTING.md) file if you would like to get involved.

## Acknowledgements

It's taken a lot of hard work to bring Modus to life, but we couldn't have done it alone. Modus is built upon _many_ open source components and projects.  We'd especially like to express our gratitude to the authors and teams of our core dependencies:

- Takeshi Yoneda, author of [Wazero](https://wazero.io/), and other contributors to the Wazero project - and to [Tetrate](https://tetrate.io/) for continuing its support of Wazero.  Modus uses Wazero to execute WebAssembly modules.
- Jens Neuse, Stefan Avram, and the rest of the team at [Wundergraph](https://wundergraph.com/).  Modus uses Wundergraph's [GraphQL Go Tools](https://github.com/wundergraph/graphql-go-tools) library to process incoming GraphQL API requests.
- Max Graey, Daniel Wirtz, and other contributors to the [AssemblyScript](https://www.assemblyscript.org/) project.  Modus chose AssemblyScript as one of its core languages because it is ideal for web developers getting started with Web Assembly.
- The [Go language](https://go.dev/) team, and also the maintainers of [TinyGo](https://tinygo.org/).  The Modus Runtime is written in Go, and the Modus Go SDK uses TinyGo.

## License

Modus and its components are Copyright 2024 Hypermode Inc., and licensed under the terms of the Apache License, Version 2.0.
See the [LICENSE](./LICENSE) file for a complete copy of the license.

If you have any questions about Modus licensing, or need an alternate license or other arrangement, please contact us at hello@hypermode.com.
