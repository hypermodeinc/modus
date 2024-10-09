/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */
import { JSON } from "json-as";

export function getJWTClaims<T>(): T {
  const claims = process.env.get("JWT_CLAIMS");
  if (!claims) {
    console.warn("No JWT claims found.");
    return instantiate<T>();
  }
  return JSON.parse<T>(claims);
}
