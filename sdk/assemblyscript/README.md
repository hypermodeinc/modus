<!-- This readme will display with the repository on GitHub. -->

# Modus SDK for AssemblyScript

## TODO: Rewrite this file in the context of Modus

![NPM Version](https://img.shields.io/npm/v/%40hypermode%2Ffunctions-as)
![NPM Downloads](https://img.shields.io/npm/dw/%40hypermode%2Ffunctions-as)

This repository contains the SDK used by [Hypermode](https://hypermode.com)
Functions written in [AssemblyScript](https://www.assemblyscript.org/).

The SDK provides a library to access the Hypermode platform, and helper classes
and functions for use with Hypermode Functions, and examples to demonstrate usage.

## Documentation

Please visit [docs.hypermode.com](https://docs.hypermode.com/) for detailed docs
covering the Hypermode platform, including this SDK.

## Getting Started

### Prerequisites

Ensure you have [Node.js](https://nodejs.org/) 22 or higher installed and activated.

You can install Node.js with any supported method, including:

- Downloading the official [prebuilt installer](https://nodejs.org/en/download)
- Using a [package manager](https://nodejs.org/en/download/package-manager)

### Trying out the examples

The [`examples`](./examples/) folder contains example projects that demonstrate the
features of the Hypermode platform.

#### Initial setup

Before compiling any example, you must first do the following:

1. Clone this repository locally:

   ```sh
   git clone https://github.com/hypermodeinc/functions-as
   ```

2. From the repository root, navigate to the `src` folder:

   ```sh
   cd ./src/
   ```

3. Install dependencies with NPM:

   ```sh
   npm install
   ```

The above steps only needs to be performed one time.

#### Compiling an example

Now that we've set up the prerequisites, we can compile an example project with the following steps:

1. From the repository root, navigate to the example's folder:

   ```sh
   cd ./examples/<name>/
   ```

2. Install dependencies with NPM:

   ```sh
   npm install
   ```

3. Build the project with NPM:

   ```sh
   npm run build
   ```

On a successful build, you'll find the output in the example's `build` folder.

> _NOTE:_ For convenience, the examples have their `package.json` configured
> with a local path to the source project, as follows:
>
> ```json
> "dependencies": {
>   "@hypermode/functions-as": "../../src",
> },
> ```
>
> In your own project, you should instead reference a version number of
> [the published package](https://www.npmjs.com/package/@hypermode/functions-as),
> as you would with any other project dependency.

### Creating your own Hypermode plugin

When you are ready to start writing your own Hypermode plugin, we recommend starting
from the [Hypermode Base Template](https://github.com/hypermodeinc/base-template).

It is pre-configured to use the latest version of this library, and has all of the
necessary configuration, scripts, and workflows needed to develop your Hypermode functions.

1. Navigate to the [Hypermode Base Template](https://github.com/hypermodeinc/base-template) repository.
2. Click the "Use this template" button in the upper-right corner to create your own repository.
3. In the `functions` folder, edit the `package.json` file to update the `name`, `version`, `description`,
   `author`, and `license` properties as needed for your purposes.
4. In the `functions/assembly` folder, begin writing your Hypermode functions in `index.ts`.
   You can replace the starter function with your own code. You can also add more files to the `assembly`
   folder and use `import` and `export` statements to connect them (as seen in the `simple` example).
   _Note that only the final function exports from `index.ts` will be registered as Hypermode Functions._
5. In the root folder of your project, edit the `hypermode.json` manifest file as needed to connect the
   models and hosts used by your functions.

### Publishing to Hypermode

Currently, your project must be created from the template project in order to publish it to Hypermode.

The Hypermode GitHub application will need to be added to your repository, and Hypermode backend will
need to be set up. During this early phase, Hypermode staff will coordinate with you to perform these
steps manually. Please contact your Hypermode representative for further details.
