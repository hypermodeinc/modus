/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { graphql } from ".";

// This file retains compatibility with previous function versions.

/**
 * @deprecated Import `graphql` instead.
 */
export abstract class connection {
  /**
   * @deprecated Use `graphql.execute` instead.
   */
  static invokeGraphqlApi<T>(
    hostName: string,
    statement: string,
    variables: QueryVariables = new QueryVariables(),
  ): GQLResponse<T> {
    const r = graphql.execute<T>(hostName, statement, variables);
    return <GQLResponse<T>>{
      errors: r.errors,
      data: r.data,
    };
  }
}

/**
 * @deprecated Import `graphql`, and use `graphql.Variables` instead.
 */
export class QueryVariables extends graphql.Variables {}

/**
 * @deprecated Import `graphql`, and use `graphql.Response` instead.
 */
@json
export class GQLResponse<T> extends graphql.Response<T> {}
