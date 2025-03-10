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
  connection: string,
  dbName: string,
  query: string,
  parametersJson: string,
): EagerResult;

/**
 * Executes a Cypher query on the Neo4j database.
 * @param connection - the name of the connection
 * @param query - the query to execute
 * @param parameters - the parameters to pass to the query
 * @param dbName - the name of the database
 * @returns The EagerResult from the Neo4j server
 */
export function executeQuery(
  connection: string,
  query: string,
  parameters: Variables = new Variables(),
  dbName: string = "neo4j",
): EagerResult {
  const paramsJson = parameters.toJSON();
  const response = hostExecuteQuery(connection, dbName, query, paramsJson);
  if (!response) {
    throw new Error("Error executing Query.");
  }

  return response;
}

/**
 * The result of a Neo4j query.
 */
@json
export class EagerResult {
  keys: string[] = [];
  records: Record[] = [];

  /**
   * @deprecated use `keys` (lowercase) instead
   */
  get Keys(): string[] {
    return this.keys;
  }

  /**
   * @deprecated use `records` (lowercase) instead
   */
  get Records(): Record[] {
    return this.records;
  }
}

/**
 * A record in a Neo4j query result.
 */
@json
export class Record {
  keys: string[] = [];
  values: string[] = [];

  /**
   * @deprecated use `keys` (lowercase) instead
   */
  get Keys(): string[] {
    return this.keys;
  }
  set Keys(value: string[]) {
    this.keys = value;
  }

  /**
   * @deprecated use `values` (lowercase) instead
   */
  get Values(): string[] {
    return this.values;
  }
  set Values(value: string[]) {
    this.values = value;
  }

  /**
  /* Get a value from a record at a given key as a JSON encoded string.
   */
  get(key: string): string {
    for (let i = 0; i < this.Keys.length; i++) {
      if (this.Keys[i] == key) {
        return this.Values[i];
      }
    }
    throw new Error("Key not found in record.");
  }

  /**
   * Get a value from a record at a given key and cast or decode it to a specific type.
   */
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
      idof<T>() === idof<JSON.Raw[]>() ||
      idof<T>() === idof<Point2D>() ||
      idof<T>() === idof<Point3D>()
    ) {
      for (let i = 0; i < this.keys.length; i++) {
        if (this.keys[i] == key) {
          return JSON.parse<T>(this.Values[i]);
        }
      }
      throw new Error("Key not found in record.");
    }
    throw new Error("Unsupported type.");
  }

  asMap(): Map<string, string> {
    const map = new Map<string, string>();
    for (let i = 0; i < this.keys.length; i++) {
      map.set(this.keys[i], this.values[i]);
    }
    return map;
  }


  @serializer
  serialize(self: Record): string {
    let out = "{";
    const end = self.keys.length - 1;
    for (let i = 0; i < end; i++) {
      const key = JSON.stringify(self.keys[i]);
      out += key + ":" + self.values[i] + ",";
    }

    if (end >= 0) {
      const key = JSON.stringify(self.keys[end]);
      out += key + ":" + self.values[end];
    }

    return out + "}";
  }


  @deserializer
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  deserialize(data: string): Record {
    throw new Error("Not implemented.");
  }
}


@json
abstract class Entity {

  @alias("ElementId")
  elementId!: string;


  @alias("Props")
  props!: DynamicMap;

  /**
   * @deprecated use `elementId` (lowercase) instead
   */
  get ElementId(): string {
    return this.elementId;
  }
  set ElementId(value: string) {
    this.elementId = value;
  }

  /**
   * @deprecated use `props` (lowercase) instead
   */
  get Props(): DynamicMap {
    return this.props;
  }
  set Props(value: DynamicMap) {
    this.props = value;
  }

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
  labels!: string[];

  /**
   * @deprecated use `labels` (lowercase) instead
   */
  get Labels(): string[] {
    return this.labels;
  }
  set Labels(value: string[]) {
    this.labels = value;
  }
}


@json
export class Relationship extends Entity {

  @alias("StartElementId")
  startElementId!: string;


  @alias("EndElementId")
  endElementId!: string;


  @alias("Type")
  type!: string;

  /**
   * @deprecated use `startElementId` (lowercase) instead
   */
  get StartElementId(): string {
    return this.startElementId;
  }
  set StartElementId(value: string) {
    this.startElementId = value;
  }

  /**
   * @deprecated use `endElementId` (lowercase) instead
   */
  get EndElementId(): string {
    return this.endElementId;
  }
  set EndElementId(value: string) {
    this.endElementId = value;
  }

  /**
   * @deprecated use `type` (lowercase) instead
   */
  get Type(): string {
    return this.type;
  }
  set Type(value: string) {
    this.type = value;
  }
}


@json
export class Path {

  @alias("Nodes")
  nodes!: Node[];


  @alias("Relationships")
  relationships!: Relationship[];

  /**
   * @deprecated use `nodes` (lowercase) instead
   */
  get Nodes(): Node[] {
    return this.nodes;
  }
  set Nodes(value: Node[]) {
    this.nodes = value;
  }

  /**
   * @deprecated use `relationships` (lowercase) instead
   */
  get Relationships(): Relationship[] {
    return this.relationships;
  }
  set Relationships(value: Relationship[]) {
    this.relationships = value;
  }
}


@json
export class Point2D {

  @alias("X")
  x!: f64;


  @alias("Y")
  y!: f64;


  @alias("SpatialRefId")
  spatialRefId!: u32;

  /**
   * @deprecated use `x` (lowercase) instead
   */
  get X(): f64 {
    return this.x;
  }
  set X(value: f64) {
    this.x = value;
  }

  /**
   * @deprecated use `y` (lowercase) instead
   */
  get Y(): f64 {
    return this.y;
  }
  set Y(value: f64) {
    this.y = value;
  }

  /**
   * @deprecated use `spatialRefId` (lowercase) instead
   */
  get SpatialRefId(): u32 {
    return this.spatialRefId;
  }
  set SpatialRefId(value: u32) {
    this.spatialRefId = value;
  }

  toString(): string {
    return `Point{SpatialRefId=${this.SpatialRefId}, X=${this.X}, Y=${this.Y}}`;
  }

  /**
   * @deprecated use `toString` instead
   */
  String(): string {
    return this.toString();
  }
}


@json
export class Point3D {

  @alias("X")
  x!: f64;


  @alias("Y")
  y!: f64;


  @alias("Z")
  z!: f64;


  @alias("SpatialRefId")
  spatialRefId!: u32;

  /**
   * @deprecated use `x` (lowercase) instead
   */
  get X(): f64 {
    return this.x;
  }
  set X(value: f64) {
    this.x = value;
  }

  /**
   * @deprecated use `y` (lowercase) instead
   */
  get Y(): f64 {
    return this.y;
  }
  set Y(value: f64) {
    this.y = value;
  }

  /**
   * @deprecated use `z` (lowercase) instead
   */
  get Z(): f64 {
    return this.z;
  }
  set Z(value: f64) {
    this.z = value;
  }

  /**
   * @deprecated use `spatialRefId` (lowercase) instead
   */
  get SpatialRefId(): u32 {
    return this.spatialRefId;
  }
  set SpatialRefId(value: u32) {
    this.spatialRefId = value;
  }

  toString(): string {
    return `Point{SpatialRefId=${this.SpatialRefId}, X=${this.X}, Y=${this.Y}, Z=${this.Z}}`;
  }

  /**
   * @deprecated use `toString` instead
   */
  String(): string {
    return this.toString();
  }
}
