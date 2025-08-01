/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSON } from "json-as";
import * as utils from "./utils";
import { NamedParams as Variables } from "./database";
export { Variables };

// @ts-expect-error: decorator
@external("modus_graphql_client", "executeQuery")
declare function hostExecuteQuery(
  connection: string,
  statement: string,
  variables: string,
): string;

export function execute<TData>(
  connection: string,
  statement: string,
  variables: Variables = new Variables(),
): Response<TData> {
  const varsJson = variables.toJSON();
  const response = hostExecuteQuery(connection, statement, varsJson);
  if (utils.resultIsInvalid(response)) {
    throw new Error("Error invoking GraphQL API.");
  }

  const results = JSON.parse<Response<TData>>(response);
  if (results.errors) {
    console.error("GraphQL API Errors:" + JSON.stringify(results.errors));
  }
  return results;
}


@json
export class Response<T> {
  errors: ErrorResult[] | null = null;
  data: T | null = null;
  // extensions: Map<string, ???> | null = null;
}


@json
class ErrorResult {
  message!: string;
  locations: CodeLocation[] | null = null;
  path: string[] | null = null;
  // extensions: Map<string, ???> | null = null;
}


@json
class CodeLocation {
  line!: u32;
  column!: u32;
}
