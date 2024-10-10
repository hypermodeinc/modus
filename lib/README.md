# Modus Shared Libraries

This directory contains shared libraries that are used by various Modus components, or externally.

Each subdirectory contains a separate Go module that is versioned independently from the rest of the repository,
using a module path and tag scheme that matches the folder path.

For example, updates to the `manifest` subdirectory are tagged as `lib/manifest/vX.Y.Z`, and are installed with

```sh
go get -u github.com/hypermodeinc/modus/lib/manifest
```

## Usage outside Modus

The shared libraries can be used for any purpose, subject to the Modus [LICENSE](../LICENSE).
Primarily they are used by the Modus runtime, and hosting infrastructure such as [Hypermode](https://hypermode.com).
