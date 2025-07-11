/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { expect, it, run } from "as-test";
import { JSON } from "json-as";
import { Point, Location } from "../database";

it("should serialize a Point object", () => {
  const point = new Point(1, 2);
  const json = JSON.stringify(point);
  expect(json).toBe(`"(1.0,2.0)"`);
});

it("should deserialize a Point object", () => {
  const json = `"(1.0,2.0)"`;
  const point = JSON.parse<Point>(json);
  expect(point.x).toBe(1);
  expect(point.y).toBe(2);
});

it("should serialize a Location object", () => {
  const location = new Location(1, 2);
  const json = JSON.stringify(location);
  expect(json).toBe(`"(1.0,2.0)"`);
});

it("should serialize a Location object", () => {
  const location = new Location(1, 2);
  const json = JSON.stringify(location);
  expect(json).toBe(`"(1.0,2.0)"`);
});

it("should deserialize a Location object", () => {
  const json = `"(1.0,2.0)"`;
  const location = JSON.parse<Location>(json);
  expect(location.longitude).toBe(1);
  expect(location.latitude).toBe(2);
});

it("should deserialize a Location point object", () => {
  const json = `"  POINT (1.0 2.0)"`;
  const location = JSON.parse<Location>(json);
  expect(location.longitude).toBe(1);
  expect(location.latitude).toBe(2);
});

it("should deserialize a Location object with whitespace", () => {
  const json = `  "  (  1.0  ,  2.0  )  "  `;
  const location = JSON.parse<Location>(json);
  expect(location.longitude).toBe(1);
  expect(location.latitude).toBe(2);
});

it("should deserialize a Location point object with whitespace", () => {
  const json = `  " POINT (1.0 2.0)  "  `;
  const location = JSON.parse<Location>(json);
  expect(location.longitude).toBe(1);
  expect(location.latitude).toBe(2);
});

it("should deserialize a Location object if an int is returned", () => {
  const json = `"(1,2)"`;
  const location = JSON.parse<Location>(json);
  expect(location.longitude).toBe(1);
  expect(location.latitude).toBe(2);
});

run();
