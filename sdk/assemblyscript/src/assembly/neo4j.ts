/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSON } from "json-as";
import { NamedParams as Variables } from "./database";
export { Variables };

// @ts-expect-error: decorator
@external("modus_neo4j_client", "executeQuery")
declare function hostExecuteQuery(
  hostName: string,
  dbName: string,
  query: string,
  parametersJson: string,
): EagerResult;

/**
 *
 * Executes a Cypher query on the Neo4j database.
 *
 * @param hostName - the name of the host
 * @param dbName - the name of the database
 * @param query - the query to execute
 * @param parameters - the parameters to pass to the query
 * @param query - the query to execute
 * @param mutations - the mutations to execute
 * @returns The EagerResult from the Neo4j server
 */
export function executeQuery(
  hostName: string,
  query: string,
  parameters: Variables = new Variables(),
  dbName: string = "neo4j",
): EagerResult {
  const paramsJson = parameters.toJSON();
  const response = hostExecuteQuery(hostName, dbName, query, paramsJson);
  if (!response) {
    throw new Error("Error executing Query.");
  }

  return response;
}

export class EagerResult {
  Keys: string[] = [];
  Records: Record[] = [];
}

export class Record {
  Values: string[] = [];
  Keys: string[] = [];

  get(key: string): string {
    for (let i = 0; i < this.Keys.length; i++) {
      if (this.Keys[i] == key) {
        return this.Values[i];
      }
    }
    throw new Error("Key not found in record.");
  }

  asMap(): Map<string, string> {
    const map = new Map<string, string>();
    for (let i = 0; i < this.Keys.length; i++) {
      map.set(this.Keys[i], this.Values[i]);
    }
    return map;
  }
}


@json
export class Node {

  @alias("ElementId")
  ElementId: string = "";


  @alias("Labels")
  Labels: string[] = [];


  @alias("Props")
  Props: Map<string, string> = new Map<string, string>();
}

export function getRecordValue<T>(record: Record, key: string): T {
  for (let i = 0; i < record.Keys.length; i++) {
    if (record.Keys[i] == key) {
      return JSON.parse<T>(record.Values[i]);
    }
  }
  throw new Error("Key not found in record.");
}
