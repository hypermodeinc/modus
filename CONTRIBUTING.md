# Contributing to Modus

We're really glad you're here and would love for you to contribute to Modus! There are a variety of
ways to make Modus better, including bug fixes, features, docs, and blog posts (among others). Every
bit helps 🙏

Please help us keep the community safe while working on the project by upholding our
[Code of Conduct](/CODE_OF_CONDUCT.md) at all times.

Before jumping to a pull request, ensure you've looked at
[PRs](https://github.com/hypermodeinc/modus/pulls) and
[issues](https://github.com/hypermodeinc/modus/issues) (open and closed) for existing work related
to your idea.

If in doubt or contemplating a larger change, join the
[Hypermode Discord](https://discord.hypermode.com) and start a discussion in the
[#modus](https://discord.com/channels/1267579648657850441/1292948253796466730) channel.

## Codebase

The primary development language of Modus is Go. The **runtime** and supporting packages are all
written in Go.

**Modus SDKs** enable users to build apps in multiple languages, currently including Go and
AssemblyScript.

The **Modus CLI** drives the local development experience and is written in TypeScript.

### Development environment

The fastest path to setting up a development environment for Modus is through VS Code. The repo
includes a set of configs to set VS Code up automatically.

#### Trunk

We use [Trunk](https://docs.trunk.io/) for linting and formatting. Please
[install the Trunk CLI](https://docs.trunk.io/cli/install) on your local development machine. You
can then use commands such as `trunk check` or `trunk fmt` for linting and formatting.

If you are using VS Code, you should also install the
[Trunk VS Code Extension](https://marketplace.visualstudio.com/items?itemName=trunk.io). The
workspace settings file included in this repo will automatically format and lint your code with
Trunk as you save your files, reducing or removing the need to use the Trunk CLI manually.

### Clone the Modus repository

To contribute code, start by forking the Modus repository. In the top-right of the
[repo](https://github.com/hypermodeinc/modus), click **Fork**. Follow the instructions to create a
fork of the repo in your GitHub workspace.

### Building and running tests

Wherever possible, we use the built-in capabilities of each language. For example, unit tests for Go
components can be run with:

```bash
go test ./...
```

### Opening a pull request

When you're ready, open a pull request against the `main` branch in the Modus repo. Include a clear,
detailed description of the changes you've made. Be sure to add and update tests and docs as needed.

We do our best to respond to PRs within a few days. If you've not heard back, feel free to ping on
Discord.

## Other ways to help

Pull requests are awesome, but there are many ways to help.

### Documentation improvements

Modus docs are maintained in a [separate repository](https://github.com/hypermodeinc/docs). Relevant
updates and issues should be opened in that repo.

### Blogging and presenting your work

Share what you're building with Modus in your preferred medium. We'd love to help amplify your work
and/or provide feedback, so get in touch if you'd like some help!

### Join the community

There are lots of people building with Modus who are excited to connect!

- Chat on [Discord](https://discord.hypermode.com)
- Join the conversation on [X](https://x.com/hypermodeinc)
- Read the latest posts on the [Blog](https://hypermode.com/blog)
- Connect with us on [LinkedIn](https://linkedin.com/company/hypermode)
