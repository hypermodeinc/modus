# Change Log

## UNRELEASED

- Send backend ID with Sentry events [#211](https://github.com/gohypermode/runtime/pull/211)

## 2024-05-30 - Version 0.8.1

- Fix compatibility with v1 `authHeader` secret [#208](https://github.com/gohypermode/runtime/pull/208)
- Fix double-escaped JSON in OpenAI inference history [#209](https://github.com/hypermodeAI/runtime/pull/209)

## 2024-05-29 - Version 0.8.0

- Add Model Inference History to runtime [#186](https://github.com/gohypermode/runtime/pull/186)
- Pass auth headers correctly when invoking a GraphQL API [#196](https://github.com/gohypermode/runtime/pull/196)
- Use shared manifest module to read `hypermode.json` [#199](https://github.com/gohypermode/runtime/pull/199)
- Pass HTTP auth secrets using v2 manifest format [#203](https://github.com/gohypermode/runtime/pull/203) [#205](https://github.com/gohypermode/runtime/pull/205)

## 2024-05-13 - Version 0.7.0

- Sentry is no longer used when `HYPERMODE_DEBUG` is enabled [#187](https://github.com/gohypermode/runtime/pull/187)
- Only listen on `localhost` when `HYPERMODE_DEBUG` is enabled, to prevent firewall prompt [#188](https://github.com/gohypermode/runtime/pull/188)
- Improve support for marshaling classes [#189](https://github.com/gohypermode/runtime/pull/189) [#191](https://github.com/gohypermode/runtime/pull/191)
- Add support for binary data fields [#190](https://github.com/gohypermode/runtime/pull/190)
- Add host function for HTTP fetch [#191](https://github.com/gohypermode/runtime/pull/191)

## 2024-05-08 - Version 0.6.6

- Remove `Access-Control-Allow-Credentials`. Add `Access-Control-Request-Headers` [#180](https://github.com/gohypermode/runtime/pull/180)
- Restrict incoming http requests methods [#182](https://github.com/gohypermode/runtime/pull/182)

## 2024-05-07 - Version 0.6.5

- Add `Access-Control-Allow-Credentials` in CORS preflight [#177](https://github.com/gohypermode/runtime/pull/177)

## 2024-05-03 - Version 0.6.4

- Add CORS support to all endpoints [#171](https://github.com/gohypermode/runtime/pull/171)
- Replace hyphens with underscores in environment variables [#172](https://github.com/gohypermode/runtime/pull/172)
- Allow comments and trailing commas in `hypermode.json` [#173](https://github.com/gohypermode/runtime/pull/173)

## 2024-05-02 - Version 0.6.3

- Update metrics collection to remove labels [#163](https://github.com/gohypermode/runtime/pull/163)
- Add environment and version to health endpoint [#164](https://github.com/gohypermode/runtime/pull/164)
- Capture function execution duration in metrics [#165](https://github.com/gohypermode/runtime/pull/165)

## 2024-04-29 - Version 0.6.2

- Traces and non-user errors are now sent to Sentry [#158](https://github.com/gohypermode/runtime/issues/158)
- Fix OpenAI text generation [#161](https://github.com/gohypermode/runtime/issues/161)

## 2024-04-26 - Version 0.6.1

- Fix GraphQL error when resulting data contains a nested null field [#150](https://github.com/gohypermode/runtime/issues/150)
- Fix GraphQL error when resolving `__typename` fields; also add `HYPERMODE_TRACE` debugging flag [#151](https://github.com/gohypermode/runtime/issues/151)
- Collect metrics and expose metrics and health endpoints [#152](https://github.com/gohypermode/runtime/issues/152)
- Add graceful shutdown for HTTP server  [#153](https://github.com/gohypermode/runtime/issues/153)
  - Note: It works correctly for system-generated and user-generated (`ctrl-C`) terminations, but [not when debugging in VS Code](https://github.com/golang/vscode-go/issues/120).
- Add version awareness [#155](https://github.com/gohypermode/runtime/issues/155)

## 2024-04-25 - Version 0.6.0

Baseline for the change log.

See git commit history for changes for this version and prior.
