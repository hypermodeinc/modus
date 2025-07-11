/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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
export function now(): string {
  const ts = hostGetTimeInZone(null);
  if (ts === null) {
    throw new Error("Failed to get the local time.");
  }
  return ts;
}

/**
 * @deprecated Use `now` instead (lowercase).
 */
export function Now(): string {
  return now();
}

/**
 * Returns the current local time, in ISO 8601 extended format.
 *
 * @param tz - A valid IANA time zone identifier, such as "America/New_York"
 */
export function nowInZone(tz: string): string {
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
 * @deprecated Use `nowInZone` instead (lowercase).
 */
export function NowInZone(tz: string): string {
  return nowInZone(tz);
}

/**
 * Returns the local time zone identifier, in IANA format.
 */
export function getTimeZone(): string {
  return process.env.get("TZ");
}

/**
 * @deprecated Use `getTimeZone` instead (lowercase).
 */
export function GetTimeZone(): string {
  return getTimeZone();
}

/**
 * Determines whether the specified time zone is valid.
 * @param tz A time zone identifier, in IANA format.
 * @returns `true` if the time zone is valid; otherwise, `false`.
 */
export function isValidTimeZone(tz: string): bool {
  if (tz === "") {
    return false;
  }
  return hostGetTimeInZone(tz) !== null;
}

/**
 * @deprecated Use `isValidTimeZone` instead (lowercase).
 */
export function IsValidTimeZone(tz: string): bool {
  return isValidTimeZone(tz);
}
