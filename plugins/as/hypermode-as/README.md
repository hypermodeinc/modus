# Hypermode AssemblyScript Library

This library is used by Hypermode plugins written in AssemblyScript.
It provides access to host functions of the Hypermode platform,
as well as additional AssemblyScript classes and helper functions
for use with Hypermode plugins.

The library itself is packaged as code, not as a compiled WASM module,
thus it does not need to be built independently.  The consuming Hypermode
plugin will compile this library as part of its own build process.

# Usage

During initial development, we'll reference this library locally.
For example, from the plugin's folder, specify the path to this library:

```sh
npm install ../hypermode-as
```

Later, we'll publish it to NPM.  Then installation will use the package name:

```sh
npm install hypermode-as
```

# Testing

We use [as-pect](https://as-pect.gitbook.io/as-pect/) for testing.
Tests are located in the `assembly/__tests__` folder.
At the moment, there are only the example tests provided with as-pect.
We will add Hypermode-specific tests soon, and remove the examples.

To run the tests:

```sh
npm run test
```
