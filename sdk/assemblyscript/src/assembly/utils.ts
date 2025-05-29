/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

/**
 * Checks if the result is invalid, by determining if it is a null pointer.
 * This will work for any managed type, regardless of whether it is nullable or not.
 */
export function resultIsInvalid<T>(result: T): bool {
  return changetype<usize>(result) == 0;
}

/**
 * Logs an error intended for the user, and exits.
 * The message will be displayed in the API response, and in the console logs.
 */
export function throwUserError(message: string): void {
  console.error(message);
  process.exit(1);
}
