# Hypermode HTTP Example Plugin

This example shows how to call an external host via HTTP `fetch`.

See [./assembly/index.ts](./assembly/index.ts) for the implementation.

Note that any host used by the `fetch` function must have a base URL defined
as a host in the [`hypermode.json`](./hypermode.json) manifest file.
This is a security measure that prevents arbitrary hosts from being called
based solely on function input values.
