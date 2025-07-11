/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

export function isErrorWithName(err: unknown): err is { name: string } {
  return typeof err === "object" && err !== null && "name" in err;
}
