import {
  Exception as __Exception,
  ExceptionState as __ExceptionState
} from "try-as/assembly/types/exception";
import {
  ErrorState as __ErrorState
} from "try-as/assembly/types/error";
import {
  UnreachableState as __UnreachableState
} from "try-as/assembly/types/unreachable";
import {
  AbortState as __AbortState
} from "try-as/assembly/types/abort";
import {
  http
} from "@hypermode/modus-sdk-as";
import {
  Quote
} from "./classes";
import {
  Exception
} from "try-as";
export function __try_getRandomQuote(): Quote {
  if (__ExceptionState.Failures > 0) {
    if (isBoolean<Quote>()) return false;
else if (isInteger<Quote>() || isFloat<Quote>()) return 0;
else if (isManaged<Quote>() || isReference<Quote>()) return changetype<Quote>(0);
else return;
  }
  const request = new http.Request("https://zenquotes.io/api/random");
  const response = http.fetch(request);
  if (!response.ok) {
    __ErrorState.error(new Error(`Failed to fetch quote. Received: ${response.status} ${response.statusText}`), "assembly/index.ts", 16, 5);
    if (isBoolean<Quote>()) return false;
else if (isInteger<Quote>() || isFloat<Quote>()) return 0;
else if (isManaged<Quote>() || isReference<Quote>()) return changetype<Quote>(0);
else return;
  }
  return response.json<Array<Quote>>()[0];
}
export function getRandomQuote(): Quote {
  const request = new http.Request("https://zenquotes.io/api/random");
  const response = http.fetch(request);
  if (!response.ok) {
    throw new Error(`Failed to fetch quote. Received: ${response.status} ${response.statusText}`);
  }
  return response.json<Array<Quote>>()[0];
}
export function safeFetchQuote(): Quote {
  do {
    return __try_getRandomQuote();
  } while (false);
  if (__ExceptionState.Failures > 0) {
    let e = new __Exception(__ExceptionState.Type);
    __ExceptionState.Failures--;
    return {
      author: "Unknown",
      quote: "Failed to fetch quote"
    }
  }
}
export function testCatching(shouldThrow: boolean): string {
  do {
    if (shouldThrow) {
      __ErrorState.error(new Error("Test error"), "assembly/index.ts", 38, 7);
      break;
    }
    return "Success";
  } while (false);
  if (__ExceptionState.Failures > 0) {
    let e = new __Exception(__ExceptionState.Type);
    __ExceptionState.Failures--;
    const err = e as Exception;
    return "Caught error: " + err.toString();
  }
  return "";
}
export function parseI8OrFallback(s: string, fallback: i8): i8 {
  do {
    let n = i8.parse(s);
    return n;
  } while (false);
  if (__ExceptionState.Failures > 0) {
    let _ = new __Exception(__ExceptionState.Type);
    __ExceptionState.Failures--;
    return fallback;
  }
}
