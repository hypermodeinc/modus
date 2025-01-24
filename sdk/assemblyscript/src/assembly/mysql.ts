/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
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

const dbType = "mysql";

export function execute(
  hostName: string,
  statement: string,
  params: Params = new Params(),
): Response {
  return db.execute(hostName, dbType, statement, params);
}

export function query<T>(
  hostName: string,
  statement: string,
  params: Params = new Params(),
): QueryResponse<T> {
  return db.query<T>(hostName, dbType, statement, params);
}

export function queryScalar<T>(
  hostName: string,
  statement: string,
  params: Params = new Params(),
): ScalarResponse<T> {
  return db.queryScalar<T>(hostName, dbType, statement, params);
}
