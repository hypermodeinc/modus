/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// @ts-expect-error: decorator
@external("modus_system", "getTimeInZone")
declare function hostGetTimeInZone(tz: string | null): string | null;

// TODO: Find or create a library for AssemblyScript that can handle time zones,
//       and then expose modus_system.getTimeZoneData for that library.

/**
 * Returns the current local time, in ISO 8601 extended format.
 *
 * The time zone is determined in the following order of precedence:
 * - If the X-Time-Zone header is present in the request, the time zone is set to the value of the header.
 * - If the TZ environment variable is set on the host, the time zone is set to the value of the variable.
 * - Otherwise, the time zone is set to the host's local time zone.
 */
export function Now(): string {
  const ts = hostGetTimeInZone(null);
  if (ts === null) {
    throw new Error("Failed to get the local time.");
  }
  return ts;
}

/**
 * Returns the current local time, in ISO 8601 extended format.
 *
 * @param tz - A valid IANA time zone identifier, such as "America/New_York"
 */
export function NowInZone(tz: string): string {
  if (tz === "") {
    throw new Error("A time zone is required.");
  }

  const ts = hostGetTimeInZone(tz);
  if (ts === null) {
    throw new Error("Failed to get the local time.");
  }

  return ts;
}

/**
 * Returns the local time zone identifier, in IANA format.
 */
export function GetTimeZone(): string {
  return process.env.get("TZ");
}
