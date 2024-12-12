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
import { DynamicMap } from "./dynamicmap";
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

  getValue<T>(key: string): T {
    if (
      isInteger<T>() ||
      isFloat<T>() ||
      isBoolean<T>() ||
      isString<T>() ||
      idof<T>() === idof<Node>() ||
      idof<T>() === idof<Relationship>() ||
      idof<T>() === idof<Path>() ||
      idof<T>() === idof<DynamicMap>() ||
      idof<T>() === idof<string[]>() ||
      idof<T>() === idof<i8[]>() ||
      idof<T>() === idof<i16[]>() ||
      idof<T>() === idof<i32[]>() ||
      idof<T>() === idof<i64[]>() ||
      idof<T>() === idof<u8[]>() ||
      idof<T>() === idof<u16[]>() ||
      idof<T>() === idof<u32[]>() ||
      idof<T>() === idof<u64[]>() ||
      idof<T>() === idof<f32[]>() ||
      idof<T>() === idof<f64[]>() ||
      idof<T>() === idof<Date>() ||
      idof<T>() === idof<JSON.Raw>() ||
      idof<T>() === idof<Point2D>() ||
      idof<T>() === idof<Point3D>()
    ) {
      for (let i = 0; i < this.Keys.length; i++) {
        if (this.Keys[i] == key) {
          return JSON.parse<T>(this.Values[i]);
        }
      }
      throw new Error("Key not found in record.");
    }
    throw new Error("Unsupported type.");
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
abstract class Entity {

  @alias("ElementId")
  ElementId!: string;


  @alias("Props")
  Props!: DynamicMap;

  getProperty<T>(key: string): T {
    if (
      isInteger<T>() ||
      isFloat<T>() ||
      isBoolean<T>() ||
      isString<T>() ||
      idof<T>() === idof<Date>() ||
      idof<T>() === idof<JSON.Raw>() ||
      idof<T>() === idof<Point2D>() ||
      idof<T>() === idof<Point3D>()
    ) {
      return this.Props.get<T>(key);
    }
    throw new Error("Unsupported type.");
  }
}


@json
export class Node extends Entity {

  @alias("Labels")
  Labels!: string[];
}


@json
export class Relationship extends Entity {

  @alias("StartElementId")
  StartElementId!: string;


  @alias("EndElementId")
  EndElementId!: string;


  @alias("Type")
  Type!: string;
}


@json
export class Path {

  @alias("Nodes")
  Nodes!: Node[];


  @alias("Relationships")
  Relationships!: Relationship[];
}


@json
export class Point2D {

  @alias("X")
  X!: f64;


  @alias("Y")
  Y!: f64;


  @alias("SpatialRefId")
  SpatialRefId!: u32;
}


@json
export class Point3D {

  @alias("X")
  X!: f64;


  @alias("Y")
  Y!: f64;


  @alias("Z")
  Z!: f64;


  @alias("SpatialRefId")
  SpatialRefId!: u32;
}
