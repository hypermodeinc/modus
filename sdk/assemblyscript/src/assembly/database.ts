/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSON } from "json-as";
import * as utils from "./utils";

// @ts-expect-error: decorator
@external("modus_sql_client", "executeQuery")
declare function hostExecuteQuery(
  connection: string,
  dbType: string,
  statement: string,
  paramsJson: string,
): HostQueryResponse;

class HostQueryResponse {
  error!: string | null;
  resultJson!: string | null;
  rowsAffected!: u32;
  lastInsertId!: u64;
}

interface Params {
  toJSON(): string;
}

export class PositionalParams implements Params {
  private data: string[] = [];

  public push<T>(val: T): void {
    this.data.push(JSON.stringify(val));
  }

  public toJSON(): string {
    return `[${this.data.join(",")}]`;
  }
}

export class NamedParams implements Params {
  private data: Map<string, string> = new Map<string, string>();

  public set<T>(name: string, value: T): void {
    this.data.set(name, JSON.stringify(value));
  }

  public toJSON(): string {
    const segments: string[] = [];
    const keys = this.data.keys();
    const values = this.data.values();

    for (let i = 0; i < this.data.size; i++) {
      const key = JSON.stringify(keys[i]);
      const value = values[i]; // already in JSON
      segments.push(`${key}:${value}`);
    }

    return `{${segments.join(",")}}`;
  }
}

export class Response {
  error: string | null = null;
  rowsAffected: u32 = 0;
  lastInsertId: u64 = 0;
}

export class QueryResponse<T> extends Response {
  rows!: T[];
}

export class ScalarResponse<T> extends Response {
  value!: T;
}

export function execute(
  connection: string,
  dbType: string,
  statement: string,
  params: Params,
): Response {
  let paramsJson = params.toJSON();

  // This flag instructs the host function not to return rows, but to simply execute the statement.
  paramsJson = "exec:" + paramsJson;

  const response = hostExecuteQuery(
    connection,
    dbType,
    statement.trim(),
    paramsJson,
  );

  if (utils.resultIsInvalid(response)) {
    throw new Error("Error performing database query.");
  }

  if (response.error) {
    console.error("Database Error: " + response.error!);
  }

  const results: Response = {
    error: response.error,
    rowsAffected: response.rowsAffected,
    lastInsertId: response.lastInsertId,
  };

  return results;
}

export function query<T>(
  connection: string,
  dbType: string,
  statement: string,
  params: Params,
): QueryResponse<T> {
  const paramsJson = params.toJSON();
  const response = hostExecuteQuery(
    connection,
    dbType,
    statement.trim(),
    paramsJson,
  );

  if (utils.resultIsInvalid(response)) {
    throw new Error("Error performing database query.");
  }

  if (response.error) {
    console.error("Database Error: " + response.error!);
  }

  const results: QueryResponse<T> = {
    error: response.error,
    rows: response.resultJson ? JSON.parse<T[]>(response.resultJson!) : [],
    rowsAffected: response.rowsAffected,
    lastInsertId: response.rowsAffected,
  };

  return results;
}

export function queryScalar<T>(
  connection: string,
  dbType: string,
  statement: string,
  params: Params,
): ScalarResponse<T> {
  const response = query<Map<string, T>>(connection, dbType, statement, params);

  if (response.rows.length == 0 || response.rows[0].size == 0) {
    throw new Error("No results returned from query.");
  }

  return <ScalarResponse<T>>{
    error: response.error,
    value: response.rows[0].values()[0],
    rowsAffected: response.rowsAffected,
    lastInsertId: response.rowsAffected,
  };
}

/**
 * Represents a point in 2D space, having `x` and `y` coordinates.
 * Correctly serializes to and from a SQL point type, in (x, y) order.
 *
 * Note that this class is identical to the Location class, but uses different field names.
 */
@json
export class Point {
  constructor(
    public x: f64,
    public y: f64,
  ) {}

  public toString(): string {
    return `(${this.x},${this.y})`;
  }

  public static fromString(data: string): Point | null {
    const p = parsePointString(data);
    if (p.length == 0) {
      return null;
    }
    return new Point(p[0], p[1]);
  }


  @serializer
  private serialize(self: Point): string {
    return `"${self}"`;
  }


  @deserializer
  private deserialize(data: string): Point {
    if (
      data.length < 7 ||
      data.charAt(0) != '"' ||
      data.charAt(data.length - 1) != '"'
    ) {
      throw new Error("Invalid Point string");
    }

    const p = parsePointString(data.substring(1, data.length - 1));
    if (p.length == 0) {
      throw new Error("Invalid Point string");
    }

    this.x = p[0];
    this.y = p[1];
    return this;
  }
}

/**
 * Represents a location on Earth, having `longitude` and `latitude` coordinates.
 * Correctly serializes to and from a SQL point type, in (longitude, latitude) order.
 *
 * Note that this class is identical to the `Point` class, but uses different field names.
 */
@json
export class Location {
  constructor(
    public longitude: f64,
    public latitude: f64,
  ) {}

  public toString(): string {
    return `(${this.longitude},${this.latitude})`;
  }

  public static fromString(data: string): Location | null {
    const p = parsePointString(data);
    if (p.length == 0) {
      return null;
    }
    return new Location(p[0], p[1]);
  }


  @serializer
  private serialize(self: Location): string {
    return `"${self}"`;
  }


  @deserializer
  private deserialize(data: string): Location {
    data = data.trim();
    const end = data.length - 1;
    if (data.length < 7 || data.charAt(0) != '"' || data.charAt(end) != '"')
      throw new Error(
        "Failed to parse Location string. Expected quotes but found none.",
      );

    const p = parsePointString(data.slice(1, end).trim());

    if (p.length == 0) {
      throw new Error("Invalid Location string");
    }

    this.longitude = p[0];
    this.latitude = p[1];
    return this;
  }
}

function parsePointString(data: string): f64[] {
  // Convert WKT point to Postgres format
  // "POINT (x y)" -> "(x, y)"
  if (data.startsWith("POINT (")) {
    data = data.slice(6).trim().replace(" ", ",");
  }

  const end = data.length - 1;

  if (data.charAt(0) != "(" || data.charAt(end) != ")") {
    console.error(`Invalid Point string: "${data}"`);
    return [];
  }

  const parts = data.slice(1, end).split(",");
  if (parts.length != 2) {
    console.error(`Invalid Point string: "${data}"`);
    return [];
  }

  const x = f64.parse(parts[0].trim());
  const y = f64.parse(parts[1].trim());
  return [x, y];
}
