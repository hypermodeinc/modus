/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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
