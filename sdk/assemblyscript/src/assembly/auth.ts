/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */
import { JSON } from "json-as";

export function getJWTClaims<T>(): T {
  const claims = process.env.get("CLAIMS");
  if (!claims) {
    console.warn("No JWT claims found.");
    return instantiate<T>();
  }
  return JSON.parse<T>(claims);
}
