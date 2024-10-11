/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { expect, it, log, mockImport, run } from "as-test";
import { Request, Response } from "../http";
import { http } from "..";
import { JSON } from "json-as";

mockImport(
  "modus_system.logMessage",
  (level: string, message: string): void => {
    if (level === "debug") {
      console.debug(message);
    } else if (level === "info") {
      console.info(message);
    } else if (level === "warning") {
      console.warn(message);
    } else if (level === "error") {
      console.error(message);
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
