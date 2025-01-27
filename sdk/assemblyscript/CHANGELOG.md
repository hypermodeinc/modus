<!-- markdownlint-disable MD013 -->

# Change Log

## _NOTICE_

This change log is retained for historical reference only, and is no longer updated. See the [root-level change log](../../CHANGELOG.md) for current versions.

## 2024-09-16 - Version 0.12.0

- Update to v2 Hypermode metadata format [#176](https://github.com/hypermodeinc/functions-as/pull/176)
- Fix display output of certain types during build [#180](https://github.com/hypermodeinc/functions-as/pull/180)
- Make json in dgraph response not nullable [#181](https://github.com/hypermodeinc/functions-as/pull/181)
- Add getLabel collection sdk function [#183](https://github.com/hypermodeinc/functions-as/pull/183)

## 2024-08-27 - Version 0.11.2

- Redo Dgraph host functions implementation [#171](https://github.com/hypermodeinc/functions-as/pull/171)
- Add get vector and search by vector functions to collections [#172](https://github.com/hypermodeinc/functions-as/pull/172)

## 2024-08-13 - Version 0.11.1

- Add Dgraph support with host functions [#165](https://github.com/hypermodeinc/functions-as/pull/165)

## 2024-08-12 - Version 0.11.0

- Add nearest neighbor classification support for collections [#156](https://github.com/hypermodeinc/functions-as/pull/156)
- Add namespace support for collections [#158](https://github.com/hypermodeinc/functions-as/pull/158)
- Add getNamespaces host function to collections [#160](https://github.com/hypermodeinc/functions-as/pull/160)
- Add cross namespace search support to collections [#161](https://github.com/hypermodeinc/functions-as/pull/161)
- Update collections example [#162](https://github.com/hypermodeinc/functions-as/pull/162)

## 2024-08-01 - Version 0.10.5

- Upgrade linters and include config in npm package [#152](https://github.com/hypermodeinc/functions-as/pull/152)

## 2024-07-23 - Version 0.10.4

- Result object scoring name changed to distance [#143](https://github.com/hypermodeinc/functions-as/pull/143)
- Collect default values for optional parameters instead of rewriting with supplied_params [#144](https://github.com/hypermodeinc/functions-as/pull/144) [#146](https://github.com/hypermodeinc/functions-as/pull/146)

## 2024-07-16 - Version 0.10.3

- Actual fix of compatibility with older QueryVariables API [#140](https://github.com/hypermodeinc/functions-as/pull/140)

## 2024-07-15 - Version 0.10.2

- No changes (re-publish of 0.10.1 to fix npm packaging issue)

## 2024-07-15 - Version 0.10.1

- Attempted fix compatibility with older QueryVariables API [#137](https://github.com/hypermodeinc/functions-as/pull/137)

## 2024-07-15 - Version 0.10.0

- Add functions and example for accessing PostgreSQL [#119](https://github.com/hypermodeinc/functions-as/pull/119)
- Make plugin version optional [#133](https://github.com/hypermodeinc/functions-as/pull/133)
- Fix transform display of optional string literal [#134](https://github.com/hypermodeinc/functions-as/pull/134)
- Fix exporting imported functions or exporting as an alias does not work [#135](https://github.com/hypermodeinc/functions-as/pull/135)

## 2024-07-10 - Version 0.9.4

- Fix transform error when reexporting function from another file [#129](https://github.com/hypermodeinc/functions-as/pull/129)

## 2024-07-09 - Version 0.9.3

- Support optional parameters [#122](https://github.com/hypermodeinc/functions-as/pull/122)
- Fix capture of nullable types [#126](https://github.com/hypermodeinc/functions-as/pull/126)

## 2024-06-26 - Version 0.9.2

- Move/rename GraphQL types to `graphql` namespace [#118](https://github.com/hypermodeinc/functions-as/pull/118)

## 2024-06-24 - Version 0.9.1

- Remove `@embedder` decorator [#112](https://github.com/hypermodeinc/functions-as/pull/112)
- Update examples [#113](https://github.com/hypermodeinc/functions-as/pull/113) [#115](https://github.com/hypermodeinc/functions-as/pull/115)
- Misc updates [#114](https://github.com/hypermodeinc/functions-as/pull/114)
  - `connection.invokeGraphqlApi` is now `graphql.execute`
  - `connection` and `inference` classes are now marked deprecated
  - `graphql` and `classification` examples are updated
- Add batch upserts [#116](https://github.com/hypermodeinc/functions-as/pull/116)

## 2024-06-21 - Version 0.9.0

_Note: Requires Hypermode Runtime v0.9.0 or newer_

- Add collections host functions [#96](https://github.com/hypermodeinc/functions-as/pull/96)
- Fixed bug in plugin build script [#87](https://github.com/gohypermode/functions-as/pull/87)
- Updated all examples to use new manifest format [#94](https://github.com/gohypermode/functions-as/pull/94)
- Support query variables of any serializable type [#99](https://github.com/gohypermode/functions-as/pull/99)
- Route `console.log` through Runtime host function [#102](https://github.com/gohypermode/functions-as/pull/102)
- Support `@embedder` decorator to denote embedded function [#104](https://github.com/hypermodeinc/functions-as/pull/104)
- Implement new models interface and rework examples [#106](https://github.com/hypermodeinc/functions-as/pull/106)

## Version 0.8.0

_skipped to align with runtime version_

## 2024-05-13 - Version 0.7.0

_Note: Requires Hypermode Runtime v0.7.0 or newer_

- Fixed threshold logic bug in `inference.classifyText` [#76](https://github.com/hypermodeinc/functions-as/pull/76)
- Align `hypermode.json` examples to changes in manifest schema [#79](https://github.com/hypermodeinc/functions-as/pull/79) [#80](https://github.com/hypermodeinc/functions-as/pull/80) [#81](https://github.com/hypermodeinc/functions-as/pull/81)
- Add `http.fetch` API [#84](https://github.com/hypermodeinc/functions-as/pull/84)

## 2024-04-25 - Version 0.6.1

- Fixed compilation transform error when there are no host functions used. [#69](https://github.com/hypermodeinc/functions-as/pull/69)

## 2024-04-25 - Version 0.6.0

_Note: Requires Hypermode Runtime v0.6.0 or newer_

- **(BREAKING)** Most APIs for host functions have been renamed. [#44](https://github.com/hypermodeinc/functions-as/pull/44)
- Types used by host functions are captured in the metadata. [#65](https://github.com/hypermodeinc/functions-as/pull/65)
- **(BREAKING)** DQL-based host functions and examples have been removed. [#66](https://github.com/hypermodeinc/functions-as/pull/66)
- **(BREAKING)** Update host functions, and improve error handling. [#67](https://github.com/hypermodeinc/functions-as/pull/67)

## 2024-04-18 - Version 0.5.0

_Note: Requires Hypermode Runtime v0.5.0_

- Internal metadata format has changed. [#39](https://github.com/hypermodeinc/functions-as/pull/39)
  - Metadata includes function signatures for improved Runtime support.
  - Compiling a project now outputs the metadata.
- **(BREAKING)** Support query parameters of different types. [#40](https://github.com/hypermodeinc/functions-as/pull/40)
- Further improvements to compiler output. [#41](https://github.com/hypermodeinc/functions-as/pull/41)
- Example project now uses a local path to the source library. [#42](https://github.com/hypermodeinc/functions-as/pull/42)
- Capture custom data type definitions in the metadata. [#44](https://github.com/hypermodeinc/functions-as/pull/44) [#52](https://github.com/hypermodeinc/functions-as/pull/52) [#53](https://github.com/hypermodeinc/functions-as/pull/53) [#55](https://github.com/hypermodeinc/functions-as/pull/55) [#56](https://github.com/hypermodeinc/functions-as/pull/56)
- Improve build scripts [#46](https://github.com/hypermodeinc/functions-as/pull/46) [#51](https://github.com/hypermodeinc/functions-as/pull/51)
- Add environment variable to debug metadata [#54](https://github.com/hypermodeinc/functions-as/pull/54)

## 2024-03-22 - Version 0.4.0

- Adds `model.generate<TData>` and `model.generateList<TData>` [#30](https://github.com/hypermodeinc/functions-as/pull/30)
- **(BREAKING)** `model.invokeTextGenerator` has been renamed to `model.generateText`

## 2024-03-14 - Version 0.3.0

- Metadata is now included during build. [#27](https://github.com/hypermodeinc/functions-as/pull/27)

  - You must also add `@hypermode/functions-as/transform` to the transforms in the `asconfig.json` file. For example:

    ```json
    "options": {
        "transform": [
            "@hypermode/functions-as/transform",
            "json-as/transform"
        ],
        "exportRuntime": true
    }
    ```

## 2024-03-13 - Version 0.2.2

- **(BREAKING)** Host functions that previously returned vector embeddings as strings now return `f64[]` instead. [#25](https://github.com/hypermodeinc/functions-as/pull/25)

_note: 0.2.1 was published prematurely and has been unpublished. Use 0.2.2._

## 2024-03-07 - Version 0.2.0

- Added `model.InvokeTextGenerator` [#20](https://github.com/hypermodeinc/functions-as/pull/20)

## 2024-02-23 - Version 0.1.0

This is the first published release of the Hypermode Functions library for AssemblyScript.

- Renamed from pre-release `hypermode-as` to `@hypermode/functions-as`
- Published to NPM
