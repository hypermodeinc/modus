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
@external("modus_dgraph_client", "executeQuery")
declare function hostExecuteQuery(
  connection: string,
  request: Request,
): Response;

// @ts-expect-error: decorator
@external("modus_dgraph_client", "alterSchema")
declare function hostAlterSchema(connection: string, schema: string): string;

// @ts-expect-error: decorator
@external("modus_dgraph_client", "dropAttribute")
declare function hostDropAttribute(connection: string, attr: string): string;

// @ts-expect-error: decorator
@external("modus_dgraph_client", "dropAllData")
declare function hostDropAllData(connection: string): string;

/**
 *
 * Executes a DQL query or mutation on the Dgraph database.
 *
 * @param connection - the name of the connection
 * @param query - the query to execute
 * @param mutations - the mutations to execute
 * @returns The response from the Dgraph server
 */
export function execute(connection: string, request: Request): Response {
  const response = hostExecuteQuery(connection, request);
  if (!response) {
    throw new Error("Error executing DQL.");
  }

  return response;
}

/**
 *
 * Alters the schema of the dgraph database
 *
 * @param connection - the name of the connection
 * @param schema - the schema to alter
 * @returns The response from the Dgraph server
 */
export function alterSchema(connection: string, schema: string): string {
  const response = hostAlterSchema(connection, schema);
  if (utils.resultIsInvalid(response)) {
    throw new Error("Error invoking DQL.");
  }

  return response;
}

/**
 *
 * Drops an attribute from the schema.
 *
 * @param connection - the name of the connection
 * @param attr - the attribute to drop
 * @returns The response from the Dgraph server
 */
export function dropAttr(connection: string, attr: string): string {
  const response = hostDropAttribute(connection, attr);
  if (utils.resultIsInvalid(response)) {
    throw new Error("Error invoking DQL.");
  }

  return response;
}

/**
 *
 * Drops all data from the database.
 *
 * @param connection - the name of the connection
 * @returns The response from the Dgraph server
 */
export function dropAll(connection: string): string {
  const response = hostDropAllData(connection);
  if (utils.resultIsInvalid(response)) {
    throw new Error("Error invoking DQL.");
  }

  return response;
}

/**
 *
 * Represents a Dgraph request.
 *
 */
export class Request {
  constructor(Query: Query | null = null, Mutations: Mutation[] | null = null) {
    if (Query) {
      this.query = Query;
    }
    if (Mutations) {
      this.mutations = Mutations;
    }
  }
  query: Query = new Query();
  mutations: Mutation[] = [];
}

/**
 *
 * Represents a Dgraph query.
 *
 */
export class Query {
  constructor(query: string = "", variables: Variables = new Variables()) {
    this.query = query;
    this.variables = variables.toMap();
  }
  query: string = "";
  variables: Map<string, string> = new Map<string, string>();
}

/**
 *
 * Represents a Dgraph mutation.
 *
 */
export class Mutation {
  constructor(
    public setJson: string = "",
    public delJson: string = "",
    public setNquads: string = "",
    public delNquads: string = "",
    public condition: string = "",
  ) {}
}

/**
 *
 * Represents a Dgraph response.
 *
 */
export class Response {
  Json: string = "";
  Uids: Map<string, string> | null = null;
}

export class Variables {
  private data: Map<string, string> = new Map<string, string>();

  public set<T>(name: string, value: T): void {
    if (isString<T>()) {
      this.data.set(name, value as string);
      return;
    } else if (isInteger<T>()) {
      this.data.set(name, JSON.stringify(value));
      return;
    } else if (isFloat<T>()) {
      this.data.set(name, JSON.stringify(value));
      return;
    } else if (isBoolean<T>()) {
      this.data.set(name, JSON.stringify(value));
      return;
    } else {
      throw new Error(
        "Unsupported variable type in dgraph. Must be string, integer, float, boolean.",
      );
    }
  }

  public toMap(): Map<string, string> {
    return this.data;
  }
}
