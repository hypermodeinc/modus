<!-- markdownlint-disable MD013 -->

# Change Log

## 2025-04-08 - CLI 0.17.3

- feat: Add update prompt on CLI usage [#811](https://github.com/hypermodeinc/modus/pull/811)

## 2025-04-01 - Runtime 0.17.8

- fix: Update Hypermode api key reading for modus local dev [#805](https://github.com/hypermodeinc/modus/pull/805)

## 2025-04-01 - CLI 0.17.2

- fix: Update Hypermode api key reading for modus local dev [#805](https://github.com/hypermodeinc/modus/pull/805)

## 2025-03-28 - Runtime 0.17.7

- feat: support new Dgraph connection string format [#803](https://github.com/hypermodeinc/modus/pull/803)

## 2025-03-20 - Runtime 0.17.6

- fix: correct json in introspection query results [#798](https://github.com/hypermodeinc/modus/pull/798)

## 2025-03-19 - Runtime 0.17.5

- fix: embed runtime version with correct path [#796](https://github.com/hypermodeinc/modus/pull/796)

## 2025-03-19 - AssemblyScript SDK 0.17.4

- fix: correct json input/output of model invocations [#795](https://github.com/hypermodeinc/modus/pull/795)

## 2025-03-18 - Runtime 0.17.4

- fix: dgraph connection should ignore key for localhost [#793](https://github.com/hypermodeinc/modus/pull/793)

## 2025-03-13 - CLI 0.17.1

- chore: cleanup unused dependencies and update remaining [#790](https://github.com/hypermodeinc/modus/pull/790)

## 2025-03-13 - AssemblyScript SDK 0.17.3

- chore: cleanup unused dependencies and update remaining [#790](https://github.com/hypermodeinc/modus/pull/790)

## 2025-03-11 - AssemblyScript SDK 0.17.2

- deps: update json-as to v1.0.0 and upgrade related code [#779](https://github.com/hypermodeinc/modus/pull/779)
- fix: support float32 vector arrays in DQL query variables [#781](https://github.com/hypermodeinc/modus/pull/781)

## 2025-03-11 - Go SDK 0.17.2

- fix: support float32 vector arrays in DQL query variables [#781](https://github.com/hypermodeinc/modus/pull/781)

## 2025-03-11 - Runtime 0.17.3

- fix: API explorer updates [#766](https://github.com/hypermodeinc/modus/pull/766)

## 2025-02-10 - AssemblyScript SDK 0.17.1

- fix: AssemblyScript ESLint import path [#737](https://github.com/hypermodeinc/modus/pull/737)
- fix: serialization of OpenAI `ToolChoice` [#739](https://github.com/hypermodeinc/modus/pull/739)
- fix: wasm filename should be last part of package name [#754](https://github.com/hypermodeinc/modus/pull/754)
- feat: minor API improvements [#758](https://github.com/hypermodeinc/modus/pull/758)

## 2025-02-10 - Go SDK 0.17.1

- chore: update Go toolchain to 1.23.6 [#757](https://github.com/hypermodeinc/modus/pull/757)
- feat: minor API improvements [#758](https://github.com/hypermodeinc/modus/pull/758)

## 2025-02-10 - Runtime 0.17.2

- fix: decode base64 strings to binary data types [#738](https://github.com/hypermodeinc/modus/pull/738)
- perf: use xsync's `LoadOrTryCompute` [#740](https://github.com/hypermodeinc/modus/pull/740)
- fix: CORS: allow all request headers [#741](https://github.com/hypermodeinc/modus/pull/741)
- fix: accept long base64 strings [#742](https://github.com/hypermodeinc/modus/pull/742)
- fix: improve dgraph auth header passing [#752](https://github.com/hypermodeinc/modus/pull/752)
- fix: jwks endpoint should use key ID if available [#753](https://github.com/hypermodeinc/modus/pull/753)
- chore: update ModusDB [#755](https://github.com/hypermodeinc/modus/pull/755)
- chore: bump tailwind to v4 [#756](https://github.com/hypermodeinc/modus/pull/756)
- chore: update Go toolchain to 1.23.6 [#757](https://github.com/hypermodeinc/modus/pull/757)

## 2025-01-24 - Runtime 0.17.1

- fix: separate modusdb instance by app [#729](https://github.com/hypermodeinc/modus/pull/729)

## 2025-01-24 - Runtime 0.17.0

- feat: add modusdb model tracing for local dev [#697](https://github.com/hypermodeinc/modus/pull/697)
- chore: update api-explorer to react 19 [#700](https://github.com/hypermodeinc/modus/pull/700)
- chore: remove localHypermodeModels list and handle 404s properly instead in local dev [#703](https://github.com/hypermodeinc/modus/pull/703)
- fix: remove fallback to default time zone [#706](https://github.com/hypermodeinc/modus/pull/706)
- feat: support MySQL database connections [#722](https://github.com/hypermodeinc/modus/pull/722)
- chore: refactoring / tests [#723](https://github.com/hypermodeinc/modus/pull/723)
- chore: show inferences error in api explorer [#726](https://github.com/hypermodeinc/modus/pull/726)

## 2025-01-24 - AssemblyScript SDK 0.17.0

- fix: resolve warning in `deserializeRawMap` [#692](https://github.com/hypermodeinc/modus/pull/692)
- fix: add json serialization support for neo4j sdk types [#699](https://github.com/hypermodeinc/modus/pull/699)
- feat: improve OpenAI model APIs and examples to better support audio, images, and tool calling [#707](https://github.com/hypermodeinc/modus/pull/707)
- feat: support MySQL database connections [#722](https://github.com/hypermodeinc/modus/pull/722)

## 2025-01-24 - Go SDK 0.17.0

- fix: add json serialization support for neo4j sdk types [#699](https://github.com/hypermodeinc/modus/pull/699)
- feat: improve OpenAI model APIs and examples to better support audio, images, and tool calling [#707](https://github.com/hypermodeinc/modus/pull/707)
- feat: support type aliases and redefinitions [#721](https://github.com/hypermodeinc/modus/pull/721)
- feat: support MySQL database connections [#722](https://github.com/hypermodeinc/modus/pull/722)

## 2025-01-24 - CLI 0.17.0

- feat: check for Binaryen installation [#725](https://github.com/hypermodeinc/modus/pull/725)

## 2025-01-09 - CLI 0.16.6

- fix: assemblyscript builds failing [#698](https://github.com/hypermodeinc/modus/pull/698)

## 2025-01-07 - CLI 0.16.5

- fix: handle space in user profile path [#696](https://github.com/hypermodeinc/modus/pull/696)

## 2025-01-07 - CLI 0.16.4

No changes. Re-released previous version to fix release issue.

## 2025-01-07 - AssemblyScript SDK 0.16.2

No changes. Re-released previous version to fix release issue.

## 2025-01-07 - CLI 0.16.3

- chore: use `toolchain` to set Go version [#684](https://github.com/hypermodeinc/modus/pull/684)

## 2025-01-07 - Runtime 0.16.2

- chore: add timezone header to api explorer [#672](https://github.com/hypermodeinc/modus/pull/672)
- chore: use `toolchain` to set Go version [#684](https://github.com/hypermodeinc/modus/pull/684)

## 2025-01-07 - Go SDK 0.16.1

- fix: update reference links in Modus schema [#670](https://github.com/hypermodeinc/modus/pull/670)
- chore: use `toolchain` to set Go version [#684](https://github.com/hypermodeinc/modus/pull/684)

## 2025-01-07 - AssemblyScript SDK 0.16.1

- fix: reimplement DynamicMap parsing [#678](https://github.com/hypermodeinc/modus/pull/678) [#681](https://github.com/hypermodeinc/modus/pull/681)

## 2025-01-06 - CLI 0.16.2

- fix: revert invalid change [#680](https://github.com/hypermodeinc/modus/pull/680)

## 2025-01-04 - CLI 0.16.1

- fix: attempt to fix windows issue [#677](https://github.com/hypermodeinc/modus/pull/677)

## 2024-12-23 - Runtime 0.16.1

- fix: resolve time zone issues [#666](https://github.com/hypermodeinc/modus/pull/666)

## 2024-12-23 - Runtime 0.16.0

- fix: unused imports should not be included in metadata [#657](https://github.com/hypermodeinc/modus/pull/657)
- fix: update runtime wasm tests [#661](https://github.com/hypermodeinc/modus/pull/661)
- fix: support TinyGo 0.35.0 [#662](https://github.com/hypermodeinc/modus/pull/662)
- feat: add APIs for local time and time zones [#663](https://github.com/hypermodeinc/modus/pull/663)

## 2024-12-23 - Go SDK 0.16.0

- chore: lint example test [#656](https://github.com/hypermodeinc/modus/pull/656)
- fix: unused imports should not be included in metadata [#657](https://github.com/hypermodeinc/modus/pull/657)
- fix: improve Go version handling [#660](https://github.com/hypermodeinc/modus/pull/660)
- fix: update runtime wasm tests [#661](https://github.com/hypermodeinc/modus/pull/661)
- fix: support TinyGo 0.35.0 [#662](https://github.com/hypermodeinc/modus/pull/662)
- feat: add APIs for local time and time zones [#663](https://github.com/hypermodeinc/modus/pull/663)

## 2024-12-23 - AssemblyScript SDK 0.16.0

- test: add tests for AssemblyScript SDK Transform [#659](https://github.com/hypermodeinc/modus/pull/659)
- feat: add APIs for local time and time zones [#663](https://github.com/hypermodeinc/modus/pull/663)

## 2024-12-23 - CLI 0.16.0

- fix: improve Go version handling [#660](https://github.com/hypermodeinc/modus/pull/660)

## 2024-12-13 - Runtime 0.15.0

- fix: doc comments from object fields should be present in generated GraphQL schema [#630](https://github.com/hypermodeinc/modus/pull/630)
- feat: add neo4j support in modus [#636](https://github.com/hypermodeinc/modus/pull/636)
- perf: improve locking code [#637](https://github.com/hypermodeinc/modus/pull/637)
- fix: make jwk conversion warn to remove sentry overflow [#641](https://github.com/hypermodeinc/modus/pull/641)
- chore: update dockerfile [#642](https://github.com/hypermodeinc/modus/pull/642)

## 2024-12-13 - Go SDK 0.15.0

- fix: vector package should return generic type in computations [#628](https://github.com/hypermodeinc/modus/pull/628)
- chore: Remove unused go package reference [#632](https://github.com/hypermodeinc/modus/pull/632)
- feat: add neo4j support in modus [#636](https://github.com/hypermodeinc/modus/pull/636)
- fix: support point types [#650](https://github.com/hypermodeinc/modus/pull/650)

## 2024-12-13 - AssemblyScript SDK 0.15.0

- fix: vector package should return generic type in computations [#628](https://github.com/hypermodeinc/modus/pull/628)
- chore: Delete extraneous copy of Anthropic model interface [#631](https://github.com/hypermodeinc/modus/pull/631)
- feat: add neo4j support in modus [#636](https://github.com/hypermodeinc/modus/pull/636)
- feat: Add `DynamicMap` type [#638](https://github.com/hypermodeinc/modus/pull/638)
- fix: Dgraph AssemblyScript example query issues [#643](https://github.com/hypermodeinc/modus/pull/643)
- fix: neo4j as type unmarshalling & support point types [#650](https://github.com/hypermodeinc/modus/pull/650)

## 2024-12-11 - CLI 0.15.0

- feat: Support `MODUS_HOME` environment variable [#639](https://github.com/hypermodeinc/modus/pull/639)

## 2024-11-27 - CLI 0.14.0

- fix: Lookup of 'latest' SDK when no Modus SDK is referenced [#625](https://github.com/hypermodeinc/modus/pull/625)

## 2024-11-27 - Runtime 0.14.1

- feat: updates to API explorer [#617](https://github.com/hypermodeinc/modus/pull/617) [#626](https://github.com/hypermodeinc/modus/pull/626)
- various dependency updates

## 2024-11-27 - AssemblyScript SDK 0.14.2

- chore: Export base Message class in OpenAI chat SDK [#616](https://github.com/hypermodeinc/modus/pull/616)

## 2024-11-27 - Go SDK 0.14.3

- chore: Export base Message class in OpenAI chat SDK [#616](https://github.com/hypermodeinc/modus/pull/616)
- fix: Filter imports to modus host functions [#623](https://github.com/hypermodeinc/modus/pull/623)

## 2024-11-25 - Go SDK 0.14.2

- fix: Reduce Go build times [#615](https://github.com/hypermodeinc/modus/pull/615)

## 2024-11-25 - CLI 0.13.10

- fix: modus new rename branch to main if not [#613](https://github.com/hypermodeinc/modus/pull/613)
- fix: use git commit with double quotes for Windows [#612](https://github.com/hypermodeinc/modus/pull/612)

## 2024-11-23 - Runtime 0.14.0

- feat: Apply in-code documentation to generated GraphQL [#519](https://github.com/hypermodeinc/modus/pull/519)
- feat: Reduce logger output during development [#576](https://github.com/hypermodeinc/modus/pull/576)
- chore: Trigger internal release pipeline at the end of the release-runtime workflow [#577](https://github.com/hypermodeinc/modus/pull/577)
- feat: Add API explorer to runtime [#578](https://github.com/hypermodeinc/modus/pull/578)
- feat: Add API explorer component to runtime [#584](https://github.com/hypermodeinc/modus/pull/584)
- fix: logic for jwks endpoint unmarshalling was incorrect [#594](https://github.com/hypermodeinc/modus/pull/594)
- chore: Refactoring schema generation [#605](https://github.com/hypermodeinc/modus/pull/605)

## 2024-11-23 - AssemblyScript SDK 0.14.1

- Update sdk/runtime compatibility file [#610](https://github.com/hypermodeinc/modus/pull/610)

## 2024-11-23 - Go SDK 0.14.1

- Update sdk/runtime compatibility file [#610](https://github.com/hypermodeinc/modus/pull/610)

## 2024-11-23 - AssemblyScript SDK 0.14.0

- Capture jsdoc comments [#519](https://github.com/hypermodeinc/modus/pull/519)

## 2024-11-23 - Go SDK 0.14.0

- Capture Go doc comments [#519](https://github.com/hypermodeinc/modus/pull/519)

## 2024-11-20 - CLI 0.13.9

- feat: use user provided name for package.json and go.mod [#595](https://github.com/hypermodeinc/modus/pull/595)
- feat: Add scarf to track installs and modus new [#596](https://github.com/hypermodeinc/modus/pull/596)

## 2024-11-18 - AssemblyScript SDK 0.13.5

- Update json-as to `0.9.26` to fix stack overflow error [#586](https://github.com/hypermodeinc/modus/pull/586)
- Updated various other dependencies

## 2024-11-08 - CLI 0.13.8

- fix: Make `modus --version` just print modus CLI's version [#563](https://github.com/hypermodeinc/modus/pull/563)
- fix: implement retry and caching for CLI downloads [#571](https://github.com/hypermodeinc/modus/pull/571) [#574](https://github.com/hypermodeinc/modus/pull/574)

## 2024-11-08 - Runtime 0.13.2

- fix: Introspection query should succeed when only mutations exist [#558](https://github.com/hypermodeinc/modus/pull/558)
- ci: Add `secrets: inherit` when calling release-info workflow [#555](https://github.com/hypermodeinc/modus/pull/555)
- chore: Refactor metadata dependencies [#564](https://github.com/hypermodeinc/modus/pull/564)
- chore: Use Go workspace to simplify project dependencies [#565](https://github.com/hypermodeinc/modus/pull/565)
- fix: "WASM Host not found in context" error on shutdown [#566](https://github.com/hypermodeinc/modus/pull/566)
- fix: `__typename` should succeed on root fields [#573](https://github.com/hypermodeinc/modus/pull/573)

## 2024-11-06 - AssemblyScript SDK 0.13.4

- fix: Don't empty previous contents of build dir [#568](https://github.com/hypermodeinc/modus/pull/568)

## 2024-11-04 - CLI 0.13.7

- Automatically generate and push releases info to R2 bucket on every release [#526](https://github.com/hypermodeinc/modus/pull/526)
- Consistent help + print enum options + validate SDK prereq [#542](https://github.com/hypermodeinc/modus/pull/542)
  - Consistent padding in the help section
  - `modus new`: Enum options need to print possible options
  - Validate SDK prereq immediately after choosing SDK
  - `modus sdk remove`: Use select prompt to allow selection
- For CLI to track non-prereleases, pull from releases json to remove rate limiting issues from github releases [#543](https://github.com/hypermodeinc/modus/pull/543)
- Modify isOnline in CLI to query releases.json [#544](https://github.com/hypermodeinc/modus/pull/544)
- remove all github api call references, add prerelease support to CLI [#545](https://github.com/hypermodeinc/modus/pull/545)
- Add CLI build CI [#547](https://github.com/hypermodeinc/modus/pull/547)
- Add CLI Lint CI [#550](https://github.com/hypermodeinc/modus/pull/550)
- Update installer script [#551](https://github.com/hypermodeinc/modus/pull/551)

## 2024-10-30 - CLI 0.13.6

- `modus new`: Initialize git repo on interactive flow [#538](https://github.com/hypermodeinc/modus/pull/538)
- `modus new`: Skip confirmation if all required params are provided [#539](https://github.com/hypermodeinc/modus/pull/539)
- `modus new`: Fix issues with Go/TinyGo version detection [#540](https://github.com/hypermodeinc/modus/pull/540)

## 2024-10-30 - AssemblyScript SDK 0.13.3

- Actually fix issue with git info capture [#537](https://github.com/hypermodeinc/modus/pull/537)

## 2024-10-30 - AssemblyScript SDK 0.13.2

- Fix issue with git info capture [#536](https://github.com/hypermodeinc/modus/pull/536)

## 2024-10-30 - Runtime 0.13.1

- Add env file callback support for auth key reloading [#520](https://github.com/hypermodeinc/modus/pull/520)
- Fix timestamp parsing bug [#527](https://github.com/hypermodeinc/modus/pull/527)

## 2024-10-30 - Go/AS SDKs 0.13.1

- Add env file to default project templates [#530](https://github.com/hypermodeinc/modus/pull/530)

## 2024-10-30 - CLI 0.13.5

- Use `<adj>-<noun>` for default app name. [#528](https://github.com/hypermodeinc/modus/pull/528)

## 2024-10-29 - CLI 0.13.4

- `modus build` should install SDK if not already installed [#524](https://github.com/hypermodeinc/modus/pull/524)

## 2024-10-29 - CLI 0.13.3

- Fix Go not found on first install [#522](https://github.com/hypermodeinc/modus/pull/522)

## 2024-10-28 - CLI 0.13.2

- Fix CLI hang on Linux [#521](https://github.com/hypermodeinc/modus/pull/521)

## 2024-10-28 - CLI 0.13.1

- Fix issues with interactive CLI prompts [#517](https://github.com/hypermodeinc/modus/pull/517)

## 2024-10-25 - Version 0.13.0 (all components)

_NOTE: This is the first fully open-source release, using the name "Modus" for the framework. "Hypermode" still refers to the company and the commercial hosting platform - but not the framework. In previous releases, the name "Hypermode" was used for all three._

- Add Modus CLI [#389](https://github.com/hypermodeinc/modus/pull/389) [#483](https://github.com/hypermodeinc/modus/pull/483) [#484](https://github.com/hypermodeinc/modus/pull/484) [#485](https://github.com/hypermodeinc/modus/pull/485)
- Support user defined jwt auth and sdk functions [#405](https://github.com/hypermodeinc/modus/pull/405)
- Migrate from Hypermode to Modus [#412](https://github.com/hypermodeinc/modus/pull/412)
- Import WasmExtractor code [#415](https://github.com/hypermodeinc/modus/pull/415)
- Import Manifest code [#416](https://github.com/hypermodeinc/modus/pull/416)
- Update the runtime's manifest usage [#417](https://github.com/hypermodeinc/modus/pull/417)
- Add Modus Go SDK [#418](https://github.com/hypermodeinc/modus/pull/418)
- Add Local Model Invocation Support [#421](https://github.com/hypermodeinc/modus/pull/421)
- Remove HTTP Timeout, Add Context Timeout on Collections [#422](https://github.com/hypermodeinc/modus/pull/422)
- Add Modus AssemblyScript SDK [#423](https://github.com/hypermodeinc/modus/pull/423)
- Add models to Modus AssemblyScript SDK [#428](https://github.com/hypermodeinc/modus/pull/428)
- Add Vectors SDK support [#431](https://github.com/hypermodeinc/modus/pull/431)
- Update Readme files [#432](https://github.com/hypermodeinc/modus/pull/432)
- Fix vulnerability in AssemblyScript SDK install script [#435](https://github.com/hypermodeinc/modus/pull/435)
- Fix potential array out of bounds in the runtime [#437](https://github.com/hypermodeinc/modus/pull/437)
- Set minimum Go version to 1.23.0 [#438](https://github.com/hypermodeinc/modus/pull/438)
- Change default for environment setting [#439](https://github.com/hypermodeinc/modus/pull/439)
- Remove compatibility code for previous versions [#441](https://github.com/hypermodeinc/modus/pull/441)
- Target Node 22 [#446](https://github.com/hypermodeinc/modus/pull/446)
- Fix object/map field stitching [#447](https://github.com/hypermodeinc/modus/pull/447)
- Use cli component instead of direct node execution modus-sdk-as [#448](https://github.com/hypermodeinc/modus/pull/448)
- Cleanup Go Modules [#450](https://github.com/hypermodeinc/modus/pull/450)
- Modularize / Rename host functions [#452](https://github.com/hypermodeinc/modus/pull/452)
- Add release pipeline for the runtime [#453](https://github.com/hypermodeinc/modus/pull/453) [#454](https://github.com/hypermodeinc/modus/pull/454)
- Remove `go generate` and fix docker build [#455](https://github.com/hypermodeinc/modus/pull/455)
- Remove AWS Secrets Manager client [#456](https://github.com/hypermodeinc/modus/pull/456)
- Make app path required [#457](https://github.com/hypermodeinc/modus/pull/457)
- Improve `.env` file handling [#458](https://github.com/hypermodeinc/modus/pull/458)
- Update command-line args and env variables [#459](https://github.com/hypermodeinc/modus/pull/459)
- Update Sentry telemetry collection rules [#460](https://github.com/hypermodeinc/modus/pull/460)
- Fix entry alignment issue with AssemblyScript maps [#461](https://github.com/hypermodeinc/modus/pull/461)
- Update to use new Modus manifest [#462](https://github.com/hypermodeinc/modus/pull/462)
- Enable GraphQL endpoints to be defined in the manifest [#464](https://github.com/hypermodeinc/modus/pull/464)
- Publish SDKs and templates via release workflows [#465](https://github.com/hypermodeinc/modus/pull/465)
- Fix AssemblyScript build failure when no Git repo is present [#475](https://github.com/hypermodeinc/modus/pull/475)
- Disable AWS Bedrock support temporarily [#479](https://github.com/hypermodeinc/modus/pull/479)
- Update SDK releases [#480](https://github.com/hypermodeinc/modus/pull/480)
- Add metadata shared library [#482](https://github.com/hypermodeinc/modus/pull/482)
- Add `.gitignore` files to default templates [#486](https://github.com/hypermodeinc/modus/pull/486)
- Fix CLI warnings about Go/TinyGo installation [#487](https://github.com/hypermodeinc/modus/pull/487)
- Remove deprecated model fields [#488](https://github.com/hypermodeinc/modus/pull/488)
- Improve dev first use log messages [#489](https://github.com/hypermodeinc/modus/pull/489)
- Highlight endpoints when running in dev [#490](https://github.com/hypermodeinc/modus/pull/490)
- Fix data race in logging adapter [#491](https://github.com/hypermodeinc/modus/pull/491)
- Add Anthropic model interface to the Go SDK [#493](https://github.com/hypermodeinc/modus/pull/493)
- Simplify and polish `modus new` experience [#494](https://github.com/hypermodeinc/modus/pull/494)
- Move hyp settings for local model invocation to env variables [#495](https://github.com/hypermodeinc/modus/pull/495) [#504](https://github.com/hypermodeinc/modus/pull/504)
- Change GraphQL SDK examples to use a generic public GraphQL API [#501](https://github.com/hypermodeinc/modus/pull/501)
- Improve file watching and fix Windows issues [#505](https://github.com/hypermodeinc/modus/pull/505)
- Improve help messages, add `modus info` and show SDK version in `modus new` [#506](https://github.com/hypermodeinc/modus/pull/506)
- Fix runtime shutdown issues with `modus dev` [#508](https://github.com/hypermodeinc/modus/pull/508)
- Monitored manifest and env files for changes [#509](https://github.com/hypermodeinc/modus/pull/509)
- Log bad GraphQL requests in dev [#510](https://github.com/hypermodeinc/modus/pull/510)
- Add JWKS endpoint key support to auth [#511](https://github.com/hypermodeinc/modus/pull/511)
- Use conventions to support GraphQL mutations and adjust query names [#513](https://github.com/hypermodeinc/modus/pull/513)

## 2024-10-02 - Version 0.12.7

- Use reader-writer lock on AWS secrets cache [#400](https://github.com/hypermodeinc/modus/pull/400)
- Improve bucket layout for FunctionExecutionDurationMilliseconds metric and add function_name label [#401](https://github.com/hypermodeinc/modus/pull/401)
- Improve JSON performance [#402](https://github.com/hypermodeinc/modus/pull/402)
- Misc performance improvements [#403](https://github.com/hypermodeinc/modus/pull/403)
- Fix error on void response [#404](https://github.com/hypermodeinc/modus/pull/404)
- Remove unused admin endpoint [#406](https://github.com/hypermodeinc/modus/pull/406)

## 2024-09-26 - Version 0.12.6

- Revert #393 and #396, then apply correct fix for field alignment issue [#397](https://github.com/hypermodeinc/modus/pull/397)

## 2024-09-26 - Version 0.12.5

- Fix AssemblyScript error unpinning objects from memory [#396](https://github.com/hypermodeinc/modus/pull/396)

## 2024-09-26 - Version 0.12.4

- Fix field alignment issue [#393](https://github.com/hypermodeinc/modus/pull/393)
- Improve error logging and debugging [#394](https://github.com/hypermodeinc/modus/pull/394)

## 2024-09-26 - Version 0.12.3

- Arrays in collections host functions should be non-nil [#384](https://github.com/hypermodeinc/modus/pull/384)
- Update error handling for function calls [#385](https://github.com/hypermodeinc/modus/pull/385)
- Fix array-like types passed via interface wrappers [#386](https://github.com/hypermodeinc/modus/pull/386)
- Cast slice values to handle json.Number and others [#387](https://github.com/hypermodeinc/modus/pull/387)
- Trap JSON unsupported value errors [#388](https://github.com/hypermodeinc/modus/pull/388)
- Adjust Sentry transactions [#390](https://github.com/hypermodeinc/modus/pull/390)

## 2024-09-24 - Version 0.12.2

- Fix missing GraphQL type schema [#376](https://github.com/hypermodeinc/modus/pull/376)
- Add FunctionExecutionDurationMillisecondsSummary metric [#377](https://github.com/hypermodeinc/modus/pull/377)
- Fix field alignment issue [#378](https://github.com/hypermodeinc/modus/pull/378)
- Improve execution plan creation [#379](https://github.com/hypermodeinc/modus/pull/379)
- Fix plan creation / registration bug [#380](https://github.com/hypermodeinc/modus/pull/380)

## 2024-09-18 - Version 0.12.1

- Fix panic from Go maps with primitive types [#370](https://github.com/hypermodeinc/modus/pull/370)
- Fix error when using a map as an input parameter [#373](https://github.com/hypermodeinc/modus/pull/373)
- Fix errors related to nil slices [#374](https://github.com/hypermodeinc/modus/pull/374)

## 2024-09-16 - Version 0.12.0

- Add language support for Hypermode Go SDK [#317](https://github.com/hypermodeinc/modus/pull/317) [#351](https://github.com/hypermodeinc/modus/pull/351) [#352](https://github.com/hypermodeinc/modus/pull/352)
- Major refactoring to support multiple guest languages [#347](https://github.com/hypermodeinc/modus/pull/347)
- Rename `hmruntime` to `hypruntime` [#348](https://github.com/hypermodeinc/modus/pull/348)
- Make empty dgraph responses nil [#355](https://github.com/hypermodeinc/modus/pull/355)
- Support objects as parameters to functions via GraphQL input types [#359](https://github.com/hypermodeinc/modus/pull/359)
- Fix GraphQL schema generation for Go functions [#360](https://github.com/hypermodeinc/modus/pull/360)
- Add getLabel to collection host functions [#361](https://github.com/hypermodeinc/modus/pull/361)
- Fix import registration issue [#364](https://github.com/hypermodeinc/modus/pull/364)
- Fix conversion of empty arrays and slices [#365](https://github.com/hypermodeinc/modus/pull/365)
- Fix wasm host not found in context [#366](https://github.com/hypermodeinc/modus/pull/366)

## 2024-08-27 - Version 0.11.2

- Refactor dgraph host functions to use single execute host function [#342](https://github.com/hypermodeinc/modus/pull/342)
- Add vector retrieval and search by vector to collections [#343](https://github.com/hypermodeinc/modus/pull/343)

## 2024-08-16 - Version 0.11.1

- Improve logger registration code [#335](https://github.com/hypermodeinc/modus/pull/335)
- Add dgraph host functions [#336](https://github.com/hypermodeinc/modus/pull/336)

## 2024-08-12 - Version 0.11.0

- Perf improvements to internal storage of hnsw data in Collections [#299](https://github.com/hypermodeinc/modus/pull/299)
- Fix type resolution issues [#304](https://github.com/hypermodeinc/modus/pull/304) [#306](https://github.com/hypermodeinc/modus/pull/306)
- Implement nearest neighbor classification in Collections [#305](https://github.com/hypermodeinc/modus/pull/305)
- Fix certain errors reporting incorrectly [#307](https://github.com/hypermodeinc/modus/pull/307)
- Improve encoding and decoding of arrays and maps [#308](https://github.com/hypermodeinc/modus/pull/308)
- Listen on both IPv4 and IPv6 for localhost [#309](https://github.com/hypermodeinc/modus/pull/309)
- Warn instead of error on some db connection failures [#310](https://github.com/hypermodeinc/modus/pull/310)
- Modularize language-specific features [#314](https://github.com/hypermodeinc/modus/pull/314) [#325](https://github.com/hypermodeinc/modus/pull/325)
- Fix error reporting with host functions [#318](https://github.com/hypermodeinc/modus/pull/318)
- Log cancellations and host function activity [#320](https://github.com/hypermodeinc/modus/pull/320)
- Add namespaces to collections to isolate storage [#321](https://github.com/hypermodeinc/modus/pull/321)
- Use more direct approach for registering host functions [#322](https://github.com/hypermodeinc/modus/pull/322) [#326](https://github.com/hypermodeinc/modus/pull/326) [#327](https://github.com/hypermodeinc/modus/pull/327)
- Add getNamespaces host function for collections [#324](https://github.com/hypermodeinc/modus/pull/324)
- Add cross namespace search to collections [#330](https://github.com/hypermodeinc/modus/pull/330)

## 2024-07-23 - Version 0.10.1

- Add HNSW indexing for collection [#285](https://github.com/hypermodeinc/modus/pull/285)
- Use default parameter value in metadata if it exists [#286](https://github.com/hypermodeinc/modus/pull/286)
- Fix memory corruption issue with multiple input parameters [#288](https://github.com/hypermodeinc/modus/pull/288)

## 2024-07-15 - Version 0.10.0

- Update manifest library and usage [#275](https://github.com/hypermodeinc/modus/pull/275)
- Support pointers when marshalling objects [#277](https://github.com/hypermodeinc/modus/pull/277)
- Add support for PostgreSQL host functions [#278](https://github.com/hypermodeinc/modus/pull/278)
- Fix dbpool reading after failed initialization [#281](https://github.com/hypermodeinc/modus/pull/281)
- Update for metadata changes [#282](https://github.com/hypermodeinc/modus/pull/282)
- Store function info with inference history [#283](https://github.com/hypermodeinc/modus/pull/283)
- Fix issues with GraphQL block quotes [wundergraph/graphql-go-tools/843](https://github.com/wundergraph/graphql-go-tools/pull/843)

## 2024-07-09 - Version 0.9.5

- Use anonymous wasm modules for better performance [#264](https://github.com/hypermodeinc/modus/pull/264)
- Assume normalized vectors when calculating cosine similarity [#265](https://github.com/hypermodeinc/modus/pull/265)
- Support optional parameters [#269](https://github.com/hypermodeinc/modus/pull/269)
- Handle null parameters [#270](https://github.com/hypermodeinc/modus/pull/2270)

## 2024-06-26 - Version 0.9.4

- Increase batch size for auto-embedding collection texts [#259](https://github.com/hypermodeinc/modus/pull/259)
- Fix error with multiline input in GraphQL query [#260](https://github.com/hypermodeinc/modus/pull/260)

## 2024-06-25 - Version 0.9.3

- Don't panic when the metadata DB is not configured [#256](https://github.com/hypermodeinc/modus/pull/256)
- Don't panic when collections are renamed or deleted [#257](https://github.com/hypermodeinc/modus/pull/257)

## 2024-06-24 - Version 0.9.2

- Add auto-embedding for collection based on text checkpoint [#250](https://github.com/hypermodeinc/modus/pull/250)
- Remove extraneous types in graphql schemas [#251](https://github.com/hypermodeinc/modus/pull/251)
- Allow arrays as inputs to host functions [#252](https://github.com/hypermodeinc/modus/pull/252)
- Add batch upsert & batch recompute for collection & on auto-embedding [#253](https://github.com/hypermodeinc/modus/pull/253)

## 2024-06-22 - Version 0.9.1

- Filter collection embedding functions from GraphQL schema [#245](https://github.com/hypermodeinc/modus/pull/245)
- Remove collection index from memory when manifest changes [#246](https://github.com/hypermodeinc/modus/pull/246)
- Fix missing execution id and plugin from logs from previous functions version [#247](https://github.com/hypermodeinc/modus/pull/247)
- Fix content type header when calling Hypermode-hosted models [#248](https://github.com/hypermodeinc/modus/pull/248)

## 2024-06-21 - Version 0.9.0

- Add nullable check in ReadString [#228](https://github.com/hypermodeinc/modus/pull/228)
- Lowercase model name before invoking for hypermode hosted models [#221](https://github.com/hypermodeinc/modus/pull/221)
- Improve HTTP error messages [#222](https://github.com/hypermodeinc/modus/pull/222)
- Add host function for direct logging [#224](https://github.com/hypermodeinc/modus/pull/224)
- Refactoring, and add helpers for calling functions [#226](https://github.com/hypermodeinc/modus/pull/226)
- Add support for new model interface [#229](https://github.com/hypermodeinc/modus/pull/229)
- Add sequential vector search [#240](https://github.com/hypermodeinc/modus/pull/240)
- Update Hypermode-hosted model endpoint URL [#242](https://github.com/hypermodeinc/modus/pull/242)
- Fix bug caused by #226 [#243](https://github.com/hypermodeinc/modus/pull/243)

## 2024-06-03 - Version 0.8.2

- Send backend ID with Sentry events [#211](https://github.com/hypermodeinc/modus/pull/211) [#213](https://github.com/hypermodeinc/modus/pull/213)
- Add some logging for secrets [#212](https://github.com/hypermodeinc/modus/pull/212)
- Update logging to include Runtime version [#215](https://github.com/hypermodeinc/modus/pull/215)

## 2024-05-30 - Version 0.8.1

- Fix compatibility with v1 `authHeader` secret [#208](https://github.com/hypermodeinc/modus/pull/208)
- Fix double-escaped JSON in OpenAI inference history [#209](https://github.com/hypermodeinc/modus/pull/209)

## 2024-05-29 - Version 0.8.0

- Add Model Inference History to runtime [#186](https://github.com/hypermodeinc/modus/pull/186)
- Pass auth headers correctly when invoking a GraphQL API [#196](https://github.com/hypermodeinc/modus/pull/196)
- Use shared manifest module to read `hypermode.json` [#199](https://github.com/hypermodeinc/modus/pull/199)
- Pass HTTP auth secrets using v2 manifest format [#203](https://github.com/hypermodeinc/modus/pull/203) [#205](https://github.com/hypermodeinc/modus/pull/205)

## 2024-05-13 - Version 0.7.0

- Sentry is no longer used when `HYPERMODE_DEBUG` is enabled [#187](https://github.com/hypermodeinc/modus/pull/187)
- Only listen on `localhost` when `HYPERMODE_DEBUG` is enabled, to prevent firewall prompt [#188](https://github.com/hypermodeinc/modus/pull/188)
- Improve support for marshaling classes [#189](https://github.com/hypermodeinc/modus/pull/189) [#191](https://github.com/hypermodeinc/modus/pull/191)
- Add support for binary data fields [#190](https://github.com/hypermodeinc/modus/pull/190)
- Add host function for HTTP fetch [#191](https://github.com/hypermodeinc/modus/pull/191)

## 2024-05-08 - Version 0.6.6

- Remove `Access-Control-Allow-Credentials`. Add `Access-Control-Request-Headers` [#180](https://github.com/hypermodeinc/modus/pull/180)
- Restrict incoming http requests methods [#182](https://github.com/hypermodeinc/modus/pull/182)

## 2024-05-07 - Version 0.6.5

- Add `Access-Control-Allow-Credentials` in CORS preflight [#177](https://github.com/hypermodeinc/modus/pull/177)

## 2024-05-03 - Version 0.6.4

- Add CORS support to all endpoints [#171](https://github.com/hypermodeinc/modus/pull/171)
- Replace hyphens with underscores in environment variables [#172](https://github.com/hypermodeinc/modus/pull/172)
- Allow comments and trailing commas in `hypermode.json` [#173](https://github.com/hypermodeinc/modus/pull/173)

## 2024-05-02 - Version 0.6.3

- Update metrics collection to remove labels [#163](https://github.com/hypermodeinc/modus/pull/163)
- Add environment and version to health endpoint [#164](https://github.com/hypermodeinc/modus/pull/164)
- Capture function execution duration in metrics [#165](https://github.com/hypermodeinc/modus/pull/165)

## 2024-04-29 - Version 0.6.2

- Traces and non-user errors are now sent to Sentry [#158](https://github.com/hypermodeinc/modus/issues/158)
- Fix OpenAI text generation [#161](https://github.com/hypermodeinc/modus/issues/161)

## 2024-04-26 - Version 0.6.1

- Fix GraphQL error when resulting data contains a nested null field [#150](https://github.com/hypermodeinc/modus/issues/150)
- Fix GraphQL error when resolving `__typename` fields; also add `HYPERMODE_TRACE` debugging flag [#151](https://github.com/hypermodeinc/modus/issues/151)
- Collect metrics and expose metrics and health endpoints [#152](https://github.com/hypermodeinc/modus/issues/152)
- Add graceful shutdown for HTTP server [#153](https://github.com/hypermodeinc/modus/issues/153)
  - Note: It works correctly for system-generated and user-generated (`ctrl-C`) terminations, but [not when debugging in VS Code](https://github.com/golang/vscode-go/issues/120).
- Add version awareness [#155](https://github.com/hypermodeinc/modus/issues/155)

## 2024-04-25 - Version 0.6.0

Baseline for the change log.

See git commit history for changes for this version and prior.
