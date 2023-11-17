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

TBD

We previously considered [as-pect](https://as-pect.gitbook.io/as-pect/),
but it appears to be unmaintained presently.
