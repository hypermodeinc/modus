/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSON } from "json-as";
import * as utils from "./utils";

// @ts-expect-error: decorator
@external("hypermode", "databaseQuery")
declare function databaseQuery(
  hostName: string,
  dbType: string,
  statement: string,
  paramsJson: string,
): HostQueryResponse;

class HostQueryResponse {
  error!: string | null;
  resultJson!: string | null;
  rowsAffected!: u32;
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
}

export class QueryResponse<T> extends Response {
  rows!: T[];
}

export class ScalarResponse<T> extends Response {
  value!: T;
}

export function execute(
  hostName: string,
  dbType: string,
  statement: string,
  params: Params,
): Response {
  const paramsJson = params.toJSON();
  const response = databaseQuery(
    hostName,
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
  };

  return results;
}

export function query<T>(
  hostName: string,
  dbType: string,
  statement: string,
  params: Params,
): QueryResponse<T> {
  const paramsJson = params.toJSON();
  const response = databaseQuery(
    hostName,
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
  };

  return results;
}

export function queryScalar<T>(
  hostName: string,
  dbType: string,
  statement: string,
  params: Params,
): ScalarResponse<T> {
  const response = query<Map<string, T>>(hostName, dbType, statement, params);

  if (response.rows.length == 0 || response.rows[0].size == 0) {
    throw new Error("No results returned from query.");
  }

  return <ScalarResponse<T>>{
    error: response.error,
    value: response.rows[0].values()[0],
    rowsAffected: response.rowsAffected,
  };
}
