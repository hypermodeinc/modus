<!-- markdownlint-disable first-line-heading -->
<div align="center">

[![modus](https://github.com/user-attachments/assets/1a6020bd-d041-4dd0-b4a9-ce01dc015b65)](https://github.com/hypermodeinc/modus)

[![GitHub License](https://img.shields.io/github/license/hypermodeinc/modus)](https://github.com/hypermodeinc/modus?tab=Apache-2.0-1-ov-file#readme)
[![chat](https://img.shields.io/discord/1267579648657850441)](https://discord.gg/NJQ4bJpffF)
[![GitHub Repo stars](https://img.shields.io/github/stars/hypermodeinc/modus)](https://github.com/hypermodeinc/modus/stargazers)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/m/hypermodeinc/modus)](https://github.com/hypermodeinc/modus/commits/main/)

</div>

<p align="center">
   <a href="https://docs.hypermode.com/modus/quickstart">Get Started</a>
   <span> · </span>
   <a href="https://docs.hypermode.com/">Docs</a>
   <span> · </span>
   <a href="https://discord.com/invite/MAZgkhP6C6">Discord</a>
<p>

**Modus is an open-source, serverless framework for building agentic systems and AI apps in
Go and AssemblyScript**

It simplifies the integration of models, context, and data. We're continuing to add additional
capabilities to better support the construction and calling of tools.

You write a function.

```ts
export function sayHello(name: string): string {
  return `Hello, ${name}!`;
}
```

Then, Modus:

- extracts the metadata of your functions
- compiles your code with optimizations based on the host environment
- caches the compiled module in memory for fast retrieval
- prepares an invocation plan for each function
- extracts connections, models, and other configuration details from the app’s manifest
- generates an API schema and activates the endpoint

You query the endpoint

```graphql
query SayHello {
  sayHello(name: "World")
}
```

In a few milliseconds, Modus:

- loads your compiled code into a sandboxed execution environment with a dedicated memory space
- runs your code, aided by host functions that power the Modus APIs
- securely queries data and AI models as needed, without exposing credentials to your code
- responds via the API result and releases the execution environment

**Now you have a production ready scalable endpoint for your AI-enabled app. AI-ready when you’re
ready. Launch and iterate.**

## Quickstart

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

[![Modus Demo](https://github.com/user-attachments/assets/12ac5db9-ca32-418c-a70d-67aef797a9e3)](https://www.youtube.com/watch?v=3CcJTXTmz88)

</div>

## What's it good for?

We believe that agentic flows are still fundamentally apps. As such, we've designed Modus
primarily as an general-purpose framework, that just happens to treat models, agents, and tools as
first-class components.

Modus is optimized for applications that require sub-second response times. We’ve made trade-offs to
prioritize speed and simplicity.

## Programming Languages

Since Modus is based on WebAssembly, you can write Modus apps in various programming languages. Each
language offers the full capabilities of the Modus framework.

Currently, the supported languages you may choose from are:

- [Go](https://go.dev/) - A general-purpose programming language originally designed by Google.
- [AssemblyScript](https://www.assemblyscript.org/) - A TypeScript-like language designed for
  WebAssembly.

Additional programming languages will be supported in the future.

## Hosting

Modus is a free, open-source framework, you’re welcome to run your Modus apps on your own hardware
or on any hosting platform that meets your needs.

However, over at [Hypermode](https://hypermode.com/) we have built additional tooling and
infrastructure to make constructing and deploying your agents and AI apps easy. Hypermode
includes models, observability, debugging, and function hosting. Just push your code to Github and
Hypermode will automatically build and deploy it on global infrastructure.

## Open Source

Modus is developed by [Hypermode](https://hypermode.com/) as an open-source project, integral but
independent from Hypermode.

We welcome external contributions. See the [CONTRIBUTING.md](./CONTRIBUTING.md) file if you would
like to get involved.

## Acknowledgements

It's taken a lot of hard work to bring Modus to life, but we couldn't have done it alone. Modus is
built upon _many_ open source components and projects. We'd especially like to express our gratitude
to the authors and teams of our core dependencies:

- Takeshi Yoneda, author of [Wazero](https://wazero.io/), and other contributors to the Wazero
  project - and to [Tetrate](https://tetrate.io/) for continuing its support of Wazero. Modus uses
  Wazero to execute WebAssembly modules.
- Jens Neuse, Stefan Avram, and the rest of the team at [Wundergraph](https://wundergraph.com/).
  Modus uses Wundergraph's [GraphQL Go Tools](https://github.com/wundergraph/graphql-go-tools)
  library to process incoming GraphQL API requests.
- Max Graey, Daniel Wirtz, and other contributors to the
  [AssemblyScript](https://www.assemblyscript.org/) project. Modus chose AssemblyScript as one of
  its core languages because it is ideal for web developers getting started with Web Assembly.
- The [Go language](https://go.dev/) team, and also the maintainers of
  [TinyGo](https://tinygo.org/). The Modus Runtime is written in Go, and the Modus Go SDK uses
  TinyGo.

## License

Modus and its components are Copyright 2024 Hypermode Inc., and licensed under the terms of the
Apache License, Version 2.0. See the [LICENSE](./LICENSE) file for a complete copy of the license.

If you have any questions about Modus licensing, or need an alternate license or other arrangement,
please contact us at hello@hypermode.com.
