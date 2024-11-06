/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import * as utils from "./utils";

// @ts-expect-error: decorator
@external("modus_modusdb_client", "dropAll")
declare function hostDropAll(): string;

// @ts-expect-error: decorator
@external("modus_modusdb_client", "dropData")
declare function hostDropData(): string;

// @ts-expect-error: decorator
@external("modus_modusdb_client", "alterSchema")
declare function hostAlterSchema(schema: string): string;

// @ts-expect-error: decorator
@external("modus_modusdb_client", "mutate")
declare function hostMutate(
  mutationReq: MutationRequest,
): Map<string, u64> | null;

// @ts-expect-error: decorator
@external("modus_modusdb_client", "query")
declare function hostQuery(query: string): Response;

/**
 *
 * Drops all from the database
 *
 * @returns The response from the ModusDB server
 */
export function dropAll(): string {
  const response = hostDropAll();
  if (utils.resultIsInvalid(response)) {
    throw new Error("Error dropping all from modusDB.");
  }

  return response;
}

/**
 *
 * Drops all data from the database
 *
 * @returns The response from the ModusDB server
 */
export function dropData(): string {
  const response = hostDropData();
  if (utils.resultIsInvalid(response)) {
    throw new Error("Error dropping data from modusDB.");
  }

  return response;
}

/**
 *
 * Alters the schema of the ModusDB database
 *
 * @param schema - the schema to alter
 * @returns The response from the ModusDB server
 */
export function alterSchema(schema: string): string {
  const response = hostAlterSchema(schema);
  if (utils.resultIsInvalid(response)) {
    throw new Error("Error altering schema in modusDB.");
  }

  return response;
}

/**
 *
 * Mutates the ModusDB database
 *
 * @param mutations - the mutations to execute
 * @returns The response from the ModusDB server
 */
export function mutate(mutationReq: MutationRequest): Map<string, u64> | null {
  const response = hostMutate(mutationReq);
  if (utils.resultIsInvalid(response)) {
    throw new Error("Error mutating data in modusDB.");
  }

  return response;
}

/**
 *
 * Queries the ModusDB database
 *
 * @param query - the query to execute
 * @returns The response from the ModusDB server
 */
export function query(query: string): Response {
  const response = hostQuery(query);
  if (utils.resultIsInvalid(response)) {
    throw new Error("Error querying data in modusDB.");
  }

  return response;
}

/**
 *
 * Represents a mutation request object for the ModusDB server
 */
export class MutationRequest {
  constructor(public Mutations: Mutation[] | null = null) {}
}

/**
 *
 * Represents a mutation in the ModusDB database
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
 * Represents a response from the ModusDB server
 *
 */
export class Response {
  Json: string = "";
  Txn: TxnContext | null = null;
  Latency: Latency | null = null;
  Metrics: Metrics | null = null;
  Uids: Map<string, string> | null = null;
  Rdf: u8[] | null = null;
  Hdrs: Map<string, ListOfString> | null = null;
}

/**
 *
 * Represents a transaction context in the ModusDB server
 *
 */
export class TxnContext {
  StartTs: u64 = 0;
  CommitTs: u64 = 0;
  Aborted: bool = false;
  Keys: string[] | null = null;
  Preds: string[] | null = null;
  Hash: string = "";
}

/**
 *
 * Represents a latency in the ModusDB server
 *
 */
export class Latency {
  ParsingNs: u64 = 0;
  ProcessingNs: u64 = 0;
  EncodingNs: u64 = 0;
  AssignTimestampNs: u64 = 0;
  TotalNs: u64 = 0;
}

/**
 *
 * Represents metrics in the ModusDB server
 *
 */
export class Metrics {
  NumUids: Map<string, u64> = new Map<string, u64>();
}

/**
 *
 * Represents a list of strings
 *
 */
export class ListOfString {
  Value: string[] | null = null;
}
