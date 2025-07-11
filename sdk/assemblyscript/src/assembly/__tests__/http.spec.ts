/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { expect, it, log, mockImport, run } from "as-test";
import { Request, Response } from "../http";
import { http } from "..";
import { JSON } from "json-as";

mockImport(
  "modus_system.logMessage",
  (level: string, message: string): void => {
    switch (level) {
      case "debug":
        console.debug(message);
        break;
      case "info":
        console.info(message);
        break;
      case "warning":
        console.warn(message);
        break;
      case "error":
        console.error(message);
        break;
    }
  },
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
mockImport("modus_http_client.fetch", (req: Request): Response => {
  const res = instantiate<Response>();
  const txt = '{"x":1,"y":2,"z":3}';
  const len = String.UTF8.byteLength(txt);
  store<ArrayBuffer>(
    changetype<usize>(res),
    new ArrayBuffer(len),
    offsetof<Response>("body"),
  );
  String.UTF8.encodeUnsafe(
    changetype<usize>(txt),
    len,
    changetype<usize>(res.body),
  );
  store<u16>(changetype<usize>(res), 200, offsetof<Response>("status"));
  store<string>(changetype<usize>(res), "OK", offsetof<Response>("statusText"));
  return res;
});

it("can run a http request", () => {
  expect(http.fetch("").status).toBe(200);
});

it("can parse a json payload", () => {
  const res = http.fetch("");
  log(res.text());
  expect(JSON.stringify(res.json<Vec3>())).toBe(res.text());
  expect(JSON.stringify(res.json<Vec3>())).toBe('{"x":1,"y":2,"z":3}');
});

it("can receive the status", () => {
  const res = http.fetch("");
  expect(res.statusText).toBe("OK");
});

run();


@json
class Vec3 {
  x: i32 = 0;
  y: i32 = 0;
  z: i32 = 0;
}
