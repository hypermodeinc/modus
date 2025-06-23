/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { http } from "@hypermode/modus-sdk-as";
import { Quote } from "./classes";
import { Exception } from "try-as";

export function getRandomQuote(): Quote {
  const request = new http.Request("https://zenquotes.io/api/random");

  const response = http.fetch(request);
  if (!response.ok) {
    throw new Error(
      `Failed to fetch quote. Received: ${response.status} ${response.statusText}`,
    );
  }

  return response.json<Quote[]>()[0];
}

export function safeFetchQuote(): Quote {
  try {
    return getRandomQuote();
  } catch (e) {
    return {
      author: "Unknown",
      quote: "Failed to fetch quote",
    };
  }
}

export function testCatching(shouldThrow: boolean): string {
  try {
    if (shouldThrow) {
      throw new Error("Test error");
    }
    return "Success";
  } catch (e) {
    const err = e as Exception;
    return "Caught error: " + err.toString();
  }
  return "";
}

export function parseI8OrFallback(s: string, fallback: i8): i8 {
  try {
    let n = i8.parse(s)
    return n;
  } catch (_) {
    return fallback;
  }
}