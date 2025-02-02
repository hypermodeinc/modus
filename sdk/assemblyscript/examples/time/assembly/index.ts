/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { localtime } from "@hypermode/modus-sdk-as";

/**
 * Returns the current UTC time.
 */
export function getUtcTime(): Date {
  return new Date(Date.now());
}

/**
 * Returns the current local time.
 */
export function getLocalTime(): string {
  return localtime.now();
}

/**
 * Returns the current time in a specified time zone.
 */
export function getTimeInZone(tz: string): string {
  return localtime.nowInZone(tz);
}

/**
 * Returns the local time zone identifier.
 */
export function getLocalTimeZone(): string {
  return localtime.getTimeZone();

  // Alternatively, you can use the following:
  // return process.env.get("TZ");
}
