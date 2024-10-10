/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSON } from "json-as";

const emptyArrayBuffer = new ArrayBuffer(0);

// @ts-expect-error: decorator
@external("modus_http_client", "httpFetch")
declare function fetchFromHost(request: Request): Response;

/**
 * Performs an HTTP request and returns the response.
 * @param requestOrUrl - Either a `Request` instance or a URL string.
 * @returns The HTTP response.
 */
export function fetch<T>(
  requestOrUrl: T,
  options: RequestOptions = new RequestOptions(),
): Response {
  let request: Request;
  if (isString<T>()) {
    const url = changetype<string>(requestOrUrl);
    request = new Request(url, options);
  } else if (idof<T>() == idof<Request>()) {
    const r = changetype<Request>(requestOrUrl);
    request = Request.clone(r, options);
  } else {
    throw new Error("Unsupported request type.");
  }

  const response = fetchFromHost(request);
  if (!response) {
    throw new Error("HTTP fetch failed. Check the logs for more information.");
  }

  return response;
}

/**
 * Represents an HTTP request.
 */
export class Request {
  /**
   * The URL to send the request to.
   * Example: "https://example.com/whatever"
   */
  readonly url: string;

  /**
   * The HTTP method to use for the request.
   * Example: "GET"
   */
  readonly method: string;

  /**
   * The HTTP headers for the request.
   */
  readonly headers: Headers;

  /**
   * The raw HTTP request body, as an `ArrayBuffer`.
   */
  readonly body: ArrayBuffer;

  /**
   * Creates a new `Request` instance.
   * @param url - The URL to send the request to.
   * @param options - The optional parameters for the request.
   */
  constructor(url: string, options: RequestOptions = new RequestOptions()) {
    this.url = url;
    this.method = options.method ? options.method!.toUpperCase() : "GET";
    this.headers = options.headers;
    this.body = options.body ? options.body!.data : emptyArrayBuffer;
  }

  /**
   * Clones the request, replacing any previous options with the new options provided.
   * @param request - The request to clone.
   * @param options - The options to set on the new request.
   * @returns - A new `Request` instance.
   */
  static clone(
    request: Request,
    options: RequestOptions = new RequestOptions(),
  ): Request {
    const newOptions = {
      method: options.method || request.method,
      body: options.body || Content.from(request.body),
    } as RequestOptions;

    if (options.headers.entries().length > 0) {
      newOptions.headers = options.headers;
    } else {
      newOptions.headers = request.headers;
    }

    return new Request(request.url, newOptions);
  }

  /**
   * Returns the request body as a text string.
   */
  text(): string {
    return Content.from(this.body).text();
  }

  /**
   * Deserializes the request body as JSON, and returns an object of type `T`.
   * @typeParam T - The type of object to return.
   * @returns An object of type `T` represented by the JSON in the request body.
   * @throws An error if the response body cannot be deserialized into the specified type.
   */
  json<T>(): T {
    return Content.from(this.body).json<T>();
  }
}

/**
 * Provides options for an HTTP request.
 */
export class RequestOptions {
  /**
   * The HTTP method to use for the request.
   * Default: "GET"
   */
  method: string | null = null;

  /**
   * The HTTP headers for the request.
   */
  headers: Headers = new Headers();

  /**
   * The HTTP request body.
   */
  body: Content | null = null;
}

/**
 * Represents content used in the body of an HTTP request or response.
 */
export class Content {
  // TODO: This should be a ReadableStream, but we can't support that yet.
  readonly data: ArrayBuffer;

  private constructor(data: ArrayBuffer) {
    this.data = data;
  }

  /**
   * Creates a new `Content` instance from the given value.
   * @param value - The value to create the content from. Either an `ArrayBuffer`, a `string`, or a JSON-serializable object.
   * @returns A new `Content` instance.
   */
  static from<T>(value: T): Content {
    if (idof<T>() == idof<ArrayBuffer>()) {
      return new Content(value as ArrayBuffer);
    } else if (isString<T>()) {
      return new Content(String.UTF8.encode(value as string));
    } else {
      return Content.from(JSON.stringify(value));
    }
  }

  /**
   * Returns the content as a text string.
   */
  text(): string {
    return String.UTF8.decode(this.data);
  }

  /**
   * Deserializes the content as JSON, and returns an object of type `T`.
   * @typeParam T - The type of object to return.
   * @returns An object of type `T` represented by the JSON in the content.
   * @throws An error if the cannot be deserialized into the specified type.
   */
  json<T>(): T {
    return JSON.parse<T>(this.text());
  }
}

/**
 * Represents an HTTP response.
 */
export class Response {
  /**
   * The HTTP response status code.
   * Example: 200
   */
  readonly status: u16 = 0;

  /**
   * The HTTP response status text.
   * Example: "OK"
   */
  readonly statusText: string = "";

  /**
   * The HTTP response headers.
   */
  readonly headers: Headers = new Headers();

  /**
   * The raw HTTP response body, as an `ArrayBuffer`.
   */
  readonly body: ArrayBuffer = emptyArrayBuffer;

  private constructor() {}

  /**
   * Returns whether the response status is OK (status code 2xx).
   */
  get ok(): bool {
    return this.status >= 200 && this.status < 300;
  }

  /**
   * Returns the response body as a text string.
   */
  text(): string {
    return Content.from(this.body).text();
  }

  /**
   * Deserializes the response body as JSON, and returns an object of type `T`.
   * @typeParam T - The type of object to return.
   * @returns An object of type `T` represented by the JSON in the response body.
   * @throws An error if the response body cannot be deserialized into the specified type.
   */
  json<T>(): T {
    return Content.from(this.body).json<T>();
  }
}

/**
 * Represents an HTTP header.
 */
class Header {
  /**
   * The header name.
   */
  name!: string;

  /**
   * The header values.
   */
  values!: string[];
}

/**
 * Represents a collection of HTTP headers.
 */
export class Headers {
  // The data key is lower case, since HTTP header keys are case-insensitive.
  // We preserve the original case in the header object's name field.
  // The first header added determines the original case.
  private data: Map<string, Header> = new Map<string, Header>();

  /**
   * Appends a new header to the collection.
   * If a header with the same name already exists, the value is appended to the existing header.
   * @param name - The header name.
   * @param value - The header value.
   */
  append(name: string, value: string): void {
    const key = name.toLowerCase();
    if (this.data.has(key)) {
      this.data.get(key).values.push(value);
    } else {
      this.data.set(key, { name, values: [value] });
    }
  }

  /**
   * Returns all of the headers.
   */
  entries(): string[][] {
    const entries: string[][] = [];
    const headers = this.data.values();
    for (let i = 0; i < headers.length; i++) {
      const header = headers[i];
      for (let j = 0; j < header.values.length; j++) {
        entries.push([header.name, header.values[j]]);
      }
    }
    return entries;
  }

  /**
   * Returns the value of the header with the given name.
   * If the header has multiple values, it returns them as a comma-separated string.
   * @param name - The header name.
   * @returns The header value, or `null` if the header does not exist.
   */
  get(name: string): string | null {
    const key = name.toLowerCase();
    if (!this.data.has(key)) {
      return null;
    }

    const header = this.data.get(key);
    return header.values.join(",");
  }

  /**
   * Creates a new `Headers` instance from a value of type `T`.
   * @param value - The value to create the headers from.
   * It can be a `string[][]`, a `Map<string, string>`, or a `Map<string, string[]>`.
   * @returns A new `Headers` instance.
   * @throws An error if the value type is not supported.
   */
  public static from<T>(value: T): Headers {
    switch (idof<T>()) {
      case idof<string[][]>(): {
        const headers = new Headers();
        const arr = changetype<string[][]>(value);
        for (let i = 0; i < arr.length; i++) {
          for (let j = 1; j < arr[i].length; j++) {
            headers.append(arr[i][0], arr[i][j]);
          }
        }
        return headers;
      }
      case idof<Map<string, string>>(): {
        const headers = new Headers();
        const map = changetype<Map<string, string>>(value);
        const keys = map.keys();
        for (let i = 0; i < keys.length; i++) {
          headers.append(keys[i], map.get(keys[i]));
        }
        return headers;
      }
      case idof<Map<string, string[]>>(): {
        const headers = new Headers();
        const map = changetype<Map<string, string[]>>(value);
        const keys = map.keys();
        for (let i = 0; i < keys.length; i++) {
          const values = map.get(keys[i]);
          for (let j = 0; j < values.length; j++) {
            headers.append(keys[i], values[j]);
          }
        }
        return headers;
      }
      default:
        throw new Error("Unsupported value type to create headers from.");
    }
  }
}
