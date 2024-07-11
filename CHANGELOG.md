# Change Log

## UNRELEASED

- Update manifest library and usage [#275](https://github.com/hypermodeAI/runtime/pull/275)
- Add support for PostgreSQL host functions [#278](https://github.com/hypermodeAI/runtime/pull/278)

## 2024-07-09 - Version 0.9.5

- Use anonymous wasm modules for better performance [#264](https://github.com/hypermodeAI/runtime/pull/264)
- Assume normalized vectors when calculating cosine similarity [#265](https://github.com/hypermodeAI/runtime/pull/265)
- Support optional parameters [#269](https://github.com/hypermodeAI/runtime/pull/269)
- Handle null parameters [#270](https://github.com/hypermodeAI/runtime/pull/2270)

## 2024-06-26 - Version 0.9.4

- Increase batch size for auto-embedding collection texts [#259](https://github.com/hypermodeAI/runtime/pull/259)
- Fix error with multiline input in GraphQL query [#260](https://github.com/hypermodeAI/runtime/pull/260)

## 2024-06-25 - Version 0.9.3

- Don't panic when the metadata DB is not configured [#256](https://github.com/hypermodeAI/runtime/pull/256)
- Don't panic when collections are renamed or deleted [#257](https://github.com/hypermodeAI/runtime/pull/257)

## 2024-06-24 - Version 0.9.2

- Add auto-embedding for collection based on text checkpoint [#250](https://github.com/hypermodeAI/runtime/pull/250)
- Remove extraneous types in graphql schemas [#251](https://github.com/hypermodeAI/runtime/pull/251)
- Allow arrays as inputs to host functions [#252](https://github.com/hypermodeAI/runtime/pull/252)
- Add batch upserts & batch recompute for collection & on auto-embedding [#253](https://github.com/hypermodeAI/runtime/pull/253)

## 2024-06-22 - Version 0.9.1

- Filter collection embedding functions from GraphQL schema [#245](https://github.com/hypermodeAI/runtime/pull/245)
- Remove collection index from memory when manifest changes [#246](https://github.com/hypermodeAI/runtime/pull/246)
- Fix missing execution id and plugin from logs from previous functions version [#247](https://github.com/hypermodeAI/runtime/pull/247)
- Fix content type header when calling Hypermode-hosted models [#248](https://github.com/hypermodeAI/runtime/pull/248)

## 2024-06-21 - Version 0.9.0

- Add nullable check in ReadString [#228](https://github.com/hypermodeAI/runtime/pull/228)
- Lowercase model name before invoking for hypermode hosted models [#221](https://github.com/hypermodeAI/runtime/pull/221)
- Improve HTTP error messages [#222](https://github.com/hypermodeAI/runtime/pull/222)
- Add host function for direct logging [#224](https://github.com/hypermodeAI/runtime/pull/224)
- Refactoring, and add helpers for calling functions [#226](https://github.com/hypermodeAI/runtime/pull/226)
- Add support for new model interface [#229](https://github.com/hypermodeAI/runtime/pull/229)
- Add sequential vector search [#240](https://github.com/hypermodeAI/runtime/pull/240)
- Update Hypermode-hosted model endpoint URL [#242](https://github.com/hypermodeAI/runtime/pull/242)
- Fix bug caused by #226 [#243](https://github.com/hypermodeAI/runtime/pull/243)

## 2024-06-03 - Version 0.8.2

- Send backend ID with Sentry events [#211](https://github.com/hypermodeAI/runtime/pull/211) [#213](https://github.com/hypermodeAI/runtime/pull/213)
- Add some logging for secrets [#212](https://github.com/hypermodeAI/runtime/pull/212)
- Update logging to include Runtime version [#215](https://github.com/hypermodeAI/runtime/pull/215)

## 2024-05-30 - Version 0.8.1

- Fix compatibility with v1 `authHeader` secret [#208](https://github.com/hypermodeAI/runtime/pull/208)
- Fix double-escaped JSON in OpenAI inference history [#209](https://github.com/hypermodeAI/runtime/pull/209)

## 2024-05-29 - Version 0.8.0

- Add Model Inference History to runtime [#186](https://github.com/hypermodeAI/runtime/pull/186)
- Pass auth headers correctly when invoking a GraphQL API [#196](https://github.com/hypermodeAI/runtime/pull/196)
- Use shared manifest module to read `hypermode.json` [#199](https://github.com/hypermodeAI/runtime/pull/199)
- Pass HTTP auth secrets using v2 manifest format [#203](https://github.com/hypermodeAI/runtime/pull/203) [#205](https://github.com/hypermodeAI/runtime/pull/205)

## 2024-05-13 - Version 0.7.0

- Sentry is no longer used when `HYPERMODE_DEBUG` is enabled [#187](https://github.com/hypermodeAI/runtime/pull/187)
- Only listen on `localhost` when `HYPERMODE_DEBUG` is enabled, to prevent firewall prompt [#188](https://github.com/hypermodeAI/runtime/pull/188)
- Improve support for marshaling classes [#189](https://github.com/hypermodeAI/runtime/pull/189) [#191](https://github.com/hypermodeAI/runtime/pull/191)
- Add support for binary data fields [#190](https://github.com/hypermodeAI/runtime/pull/190)
- Add host function for HTTP fetch [#191](https://github.com/hypermodeAI/runtime/pull/191)

## 2024-05-08 - Version 0.6.6

- Remove `Access-Control-Allow-Credentials`. Add `Access-Control-Request-Headers` [#180](https://github.com/hypermodeAI/runtime/pull/180)
- Restrict incoming http requests methods [#182](https://github.com/hypermodeAI/runtime/pull/182)

## 2024-05-07 - Version 0.6.5

- Add `Access-Control-Allow-Credentials` in CORS preflight [#177](https://github.com/hypermodeAI/runtime/pull/177)

## 2024-05-03 - Version 0.6.4

- Add CORS support to all endpoints [#171](https://github.com/hypermodeAI/runtime/pull/171)
- Replace hyphens with underscores in environment variables [#172](https://github.com/hypermodeAI/runtime/pull/172)
- Allow comments and trailing commas in `hypermode.json` [#173](https://github.com/hypermodeAI/runtime/pull/173)

## 2024-05-02 - Version 0.6.3

- Update metrics collection to remove labels [#163](https://github.com/hypermodeAI/runtime/pull/163)
- Add environment and version to health endpoint [#164](https://github.com/hypermodeAI/runtime/pull/164)
- Capture function execution duration in metrics [#165](https://github.com/hypermodeAI/runtime/pull/165)

## 2024-04-29 - Version 0.6.2

- Traces and non-user errors are now sent to Sentry [#158](https://github.com/hypermodeAI/runtime/issues/158)
- Fix OpenAI text generation [#161](https://github.com/hypermodeAI/runtime/issues/161)

## 2024-04-26 - Version 0.6.1

- Fix GraphQL error when resulting data contains a nested null field [#150](https://github.com/hypermodeAI/runtime/issues/150)
- Fix GraphQL error when resolving `__typename` fields; also add `HYPERMODE_TRACE` debugging flag [#151](https://github.com/hypermodeAI/runtime/issues/151)
- Collect metrics and expose metrics and health endpoints [#152](https://github.com/hypermodeAI/runtime/issues/152)
- Add graceful shutdown for HTTP server  [#153](https://github.com/hypermodeAI/runtime/issues/153)
  - Note: It works correctly for system-generated and user-generated (`ctrl-C`) terminations, but [not when debugging in VS Code](https://github.com/golang/vscode-go/issues/120).
- Add version awareness [#155](https://github.com/hypermodeAI/runtime/issues/155)

## 2024-04-25 - Version 0.6.0

Baseline for the change log.

See git commit history for changes for this version and prior.
