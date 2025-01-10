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

  __INITIALIZE(): this {
    return this;
  }

  __SERIALIZE(): string {
    return JSON.stringify(this);
  }

  /* eslint-disable @typescript-eslint/no-unused-vars */
  __DESERIALIZE(
    data: string,
    key_start: i32,
    key_end: i32,
    value_start: i32,
    value_end: i32,
  ): boolean {
    const obj = JSON.parse<EagerResult>(data);
    this.Keys = obj.Keys;
    this.Records = obj.Records;
    return true;
  }
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
      idof<T>() === idof<JSON.Raw[]>() ||
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

  __INITIALIZE(): this {
    return this;
  }

  __SERIALIZE(): string {
    return JSON.stringify(this);
  }

  /* eslint-disable @typescript-eslint/no-unused-vars */
  __DESERIALIZE(
    data: string,
    key_start: i32,
    key_end: i32,
    value_start: i32,
    value_end: i32,
  ): boolean {
    const obj = JSON.parse<Record>(data);
    this.Keys = obj.Keys;
    this.Values = obj.Values;
    return true;
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

  __INITIALIZE(): this {
    return this;
  }

  __SERIALIZE(): string {
    return JSON.stringify(this);
  }

  /* eslint-disable @typescript-eslint/no-unused-vars */
  __DESERIALIZE(
    data: string,
    key_start: i32,
    key_end: i32,
    value_start: i32,
    value_end: i32,
  ): boolean {
    const obj = JSON.parse<Node>(data);
    this.ElementId = obj.ElementId;
    this.Props = obj.Props;
    this.Labels = obj.Labels;
    return true;
  }
}


@json
export class Relationship extends Entity {

  @alias("StartElementId")
  StartElementId!: string;


  @alias("EndElementId")
  EndElementId!: string;


  @alias("Type")
  Type!: string;

  __INITIALIZE(): this {
    return this;
  }

  __SERIALIZE(): string {
    return JSON.stringify(this);
  }

  /* eslint-disable @typescript-eslint/no-unused-vars */
  __DESERIALIZE(
    data: string,
    key_start: i32,
    key_end: i32,
    value_start: i32,
    value_end: i32,
  ): boolean {
    const obj = JSON.parse<Relationship>(data);
    this.ElementId = obj.ElementId;
    this.Props = obj.Props;
    this.StartElementId = obj.StartElementId;
    this.EndElementId = obj.EndElementId;
    this.Type = obj.Type;
    return true;
  }
}


@json
export class Path {

  @alias("Nodes")
  Nodes!: Node[];


  @alias("Relationships")
  Relationships!: Relationship[];

  __INITIALIZE(): this {
    return this;
  }

  __SERIALIZE(): string {
    return JSON.stringify(this);
  }

  /* eslint-disable @typescript-eslint/no-unused-vars */
  __DESERIALIZE(
    data: string,
    key_start: i32,
    key_end: i32,
    value_start: i32,
    value_end: i32,
  ): boolean {
    const obj = JSON.parse<Path>(data);
    this.Nodes = obj.Nodes;
    this.Relationships = obj.Relationships;
    return true;
  }
}


@json
export class Point2D {

  @alias("X")
  X!: f64;


  @alias("Y")
  Y!: f64;


  @alias("SpatialRefId")
  SpatialRefId!: u32;

  String(): string {
    return `Point{SpatialRefId=${this.SpatialRefId}, X=${this.X}, Y=${this.Y}}`;
  }

  __INITIALIZE(): this {
    return this;
  }

  __SERIALIZE(): string {
    return JSON.stringify(this);
  }

  /* eslint-disable @typescript-eslint/no-unused-vars */
  __DESERIALIZE(
    data: string,
    key_start: i32,
    key_end: i32,
    value_start: i32,
    value_end: i32,
  ): boolean {
    const obj = JSON.parse<Point2D>(data);
    this.X = obj.X;
    this.Y = obj.Y;
    this.SpatialRefId = obj.SpatialRefId;
    return true;
  }
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

  String(): string {
    return `Point{SpatialRefId=${this.SpatialRefId}, X=${this.X}, Y=${this.Y}, Z=${this.Z}}`;
  }

  __INITIALIZE(): this {
    return this;
  }

  __SERIALIZE(): string {
    return JSON.stringify(this);
  }

  /* eslint-disable @typescript-eslint/no-unused-vars */
  __DESERIALIZE(
    data: string,
    key_start: i32,
    key_end: i32,
    value_start: i32,
    value_end: i32,
  ): boolean {
    const obj = JSON.parse<Point3D>(data);
    this.X = obj.X;
    this.Y = obj.Y;
    this.Z = obj.Z;
    this.SpatialRefId = obj.SpatialRefId;
    return true;
  }
}
