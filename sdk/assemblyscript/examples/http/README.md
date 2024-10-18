# Modus HTTP Example

This example shows how to call an external host via HTTP `fetch`.

See [./assembly/index.ts](./assembly/index.ts) for the implementation.

Note that any host used by the `fetch` function must have a base URL defined
as a connection in the [`modus.json`](./modus.json) manifest file.
This is a security measure that prevents arbitrary hosts from being called
based solely on function input values.
