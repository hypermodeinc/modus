/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import * as db from "./database";
import {
  PositionalParams as Params,
  Response,
  QueryResponse,
  ScalarResponse,
  Point,
  Location,
} from "./database";

export { Params, Response, QueryResponse, ScalarResponse, Point, Location };

const dbType = "postgresql";

export function execute(
  connection: string,
  statement: string,
  params: Params = new Params(),
): Response {
  return db.execute(connection, dbType, statement, params);
}

export function query<T>(
  connection: string,
  statement: string,
  params: Params = new Params(),
): QueryResponse<T> {
  return db.query<T>(connection, dbType, statement, params);
}

export function queryScalar<T>(
  connection: string,
  statement: string,
  params: Params = new Params(),
): ScalarResponse<T> {
  return db.queryScalar<T>(connection, dbType, statement, params);
}
