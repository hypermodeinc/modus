/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import * as utils from "./utils";

// @ts-expect-error: decorator
@external("modus_secrets", "getSecretValue")
declare function hostGetSecretValue(name: string): string;

/**
 * Retrieves a secret value from the host environment.
 * Throws if the secret is not found or if an error occurs.
 */
export function getSecretValue(name: string): string {
  const result = hostGetSecretValue(name);
  if (utils.resultIsInvalid(result)) {
    throw new Error("Secret not found");
  }
  return result;
}
