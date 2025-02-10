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
 * @deprecated Use `executeQuery` or `executeMutations` instead,
 * which don't require creating a `Request` object explicitly.
 */
export function execute(connection: string, request: Request): Response {
  return _execute(connection, request);
}

/**
 * Executes a DQL query on the Dgraph database, optionally with mutations.
 * @param connection - the name of the connection
 * @param query - the query to execute
 * @param mutations - the mutations to execute, if any
 * @returns The response from the Dgraph server
 * @remarks This function is a convenience wrapper around the `execute` function.
 */
export function executeQuery(
  connection: string,
  query: Query,
  ...mutations: Mutation[]
): Response {
  const request = new Request(query, mutations);
  return _execute(connection, request);
}

/**
 * Executes one or more mutations on the Dgraph database.
 * @param connection - the name of the connection
 * @param mutations - the mutations to execute
 * @returns The response from the Dgraph server
 * @remarks This function is a convenience wrapper around the `execute` function.
 */
export function executeMutations(
  connection: string,
  ...mutations: Mutation[]
): Response {
  const request = new Request(null, mutations);
  return _execute(connection, request);
}

// @ts-expect-error: decorator
@inline
function _execute(connection: string, request: Request): Response {
  const response = hostExecuteQuery(connection, request);
  if (!response) {
    throw new Error("Error executing DQL.");
  }

  return response;
}

/**
 * Alters the schema of the dgraph database
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
 * Drops an attribute from the schema.
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
 * Drops all data from the database.
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
 * Represents a Dgraph request.
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
 * Represents a Dgraph query.
 */
export class Query {
  constructor(query: string = "", variables: Variables = new Variables()) {
    this.query = query;
    this.variables = variables.toMap();
  }
  query: string = "";
  variables: Map<string, string> = new Map<string, string>();

  /**
   * Adds a variable to the query.
   */
  withVariable<T>(name: string, value: T): this {
    if (isString<T>()) {
      this.variables.set(name, value as string);
      return this;
    } else if (isInteger<T>()) {
      this.variables.set(name, JSON.stringify(value));
      return this;
    } else if (isFloat<T>()) {
      this.variables.set(name, JSON.stringify(value));
      return this;
    } else if (isBoolean<T>()) {
      this.variables.set(name, JSON.stringify(value));
      return this;
    } else {
      throw new Error(
        "Unsupported DQL variable type. Must be string, integer, float, or boolean.",
      );
    }
  }
}

/**
 * Represents a Dgraph mutation.
 */
export class Mutation {
  /**
   * Creates a new Mutation object.
   * @param setJson - JSON for setting data
   * @param delJson - JSON for deleting data
   * @param setNquads - RDF N-Quads for setting data
   * @param delNquads - RDF N-Quads for deleting data
   * @param condition - Condition for the mutation, as a DQL `@if` directive
   *
   * @remarks - All parameters are optional.
   * For clarity, prefer passing no arguments here, but use the `with*` methods instead.
   */
  constructor(
    public setJson: string = "",
    public delJson: string = "",
    public setNquads: string = "",
    public delNquads: string = "",
    public condition: string = "",
  ) {}

  /**
   * Adds a JSON string for setting data to the mutation.
   */
  withSetJson(json: string): Mutation {
    this.setJson = json;
    return this;
  }

  /**
   * Adds a JSON string for deleting data to the mutation.
   */
  withDelJson(json: string): Mutation {
    this.delJson = json;
    return this;
  }

  /**
   * Adds RDF N-Quads for setting data to the mutation.
   */
  withSetNquads(nquads: string): Mutation {
    this.setNquads = nquads;
    return this;
  }

  /**
   * Adds RDF N-Quads for deleting data to the mutation.
   */
  withDelNquads(nquads: string): Mutation {
    this.delNquads = nquads;
    return this;
  }

  /**
   * Adds a condition to the mutation, as a DQL `@if` directive.
   */
  withCondition(cond: string): Mutation {
    this.condition = cond;
    return this;
  }
}

/**
 * Represents a Dgraph response.
 */
export class Response {
  Json: string = "";
  Uids: Map<string, string> | null = null;
}

/**
 * Represents a set of variables to be used in a Dgraph query.
 */
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
        "Unsupported DQL variable type. Must be string, integer, float, or boolean.",
      );
    }
  }

  public toMap(): Map<string, string> {
    return this.data;
  }
}

/**
 * Ensures proper escaping of RDF string literals
 */
export function escapeRDF(value: string): string {
  let output = "";
  for (let i = 0; i < value.length; i++) {
    const cc = value.charCodeAt(i);
    switch (cc) {
      case 0x5c:
        output += "\\\\";
        break;
      case 0x22:
        output += '\\"';
        break;
      case 0x0a:
        output += "\\n";
        break;
      case 0x0d:
        output += "\\r";
        break;
      case 0x09:
        output += "\\t";
        break;
      default: {
        // handle control characters
        if (cc < 0x20 || (cc >= 0x7f && cc < 0x9f)) {
          output += "\\u" + cc.toString(16).padStart(4, "0");
        } else {
          output += value.charAt(i);
        }
      }
    }
  }

  return output;
}
