/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { expect, it, run } from "as-test"
import { vectors } from ".."

it("should add two vectors together", () => {
  const a = [1, 2, 3]
  const b = [4, 5, 6]
  const result = vectors.add(a, b)
  expect(result).toBe([5, 7, 9])
})

it("should add two vectors together in place", () => {
  const a = [1, 2, 3]
  const b = [4, 5, 6]
  vectors.addInPlace(a, b)
  expect(a).toBe([5, 7, 9])
})

it("should subtract two vectors", () => {
  const a = [1, 2, 3]
  const b = [4, 5, 6]
  const result = vectors.subtract(a, b)
  expect(result).toBe([-3, -3, -3])
})

it("should subtract two vectors in place", () => {
  const a = [1, 2, 3]
  const b = [4, 5, 6]
  vectors.subtractInPlace(a, b)
  expect(a).toBe([-3, -3, -3])
})

it("should add a number to a vector", () => {
  const a = [1, 2, 3]
  const b = 4
  const result = vectors.addNumber(a, b)
  expect(result).toBe([5, 6, 7])
})

it("should add a number to a vector in place", () => {
  const a = [1, 2, 3]
  const b = 4
  vectors.addNumberInPlace(a, b)
  expect(a).toBe([5, 6, 7])
})

it("should subtract a number from a vector", () => {
  const a = [1, 2, 3]
  const b = 4
  const result = vectors.subtractNumber(a, b)
  expect(result).toBe([-3, -2, -1])
})

it("should subtract a number from a vector in place", () => {
  const a = [1, 2, 3]
  const b = 4
  vectors.subtractNumberInPlace(a, b)
  expect(a).toBe([-3, -2, -1])
})

it("should multiply a vector by a number", () => {
  const a = [1, 2, 3]
  const b = 4
  const result = vectors.multiplyNumber(a, b)
  expect(result).toBe([4, 8, 12])
})

it("should multiply a vector by a number in place", () => {
  const a = [1, 2, 3]
  const b = 4
  vectors.multiplyNumberInPlace(a, b)
  expect(a).toBe([4, 8, 12])
})

it("should divide a vector by a number", () => {
  const a = [4, 8, 12]
  const b = 4
  const result = vectors.divideNumber(a, b)
  expect(result).toBe([1, 2, 3])
})

it("should divide a vector by a number in place", () => {
  const a = [4, 8, 12]
  const b = 4
  vectors.divideNumberInPlace(a, b)
  expect(a).toBe([1, 2, 3])
})

it("should compute the dot product of two vectors", () => {
  const a = [1, 2, 3]
  const b = [4, 5, 6]
  const result = vectors.dot(a, b)
  expect(result).toBe(32)
})

it("should compute the magnitude of a vector", () => {
  const a = [1, 2, 3]
  const result = vectors.magnitude(a)
  expect(result).toBe(sqrt<f64>(14))
})

it("should normalize a vector", () => {
  const a = [1, 2, 3]
  const result = vectors.normalize(a)
  const magnitude = vectors.magnitude(result)
  expect(magnitude).toBe(1)
})

it("should compute the sum of a vector", () => {
  const a = [1, 2, 3]
  const result = vectors.sum(a)
  expect(result).toBe(6)
})

it("should compute the product of a vector", () => {
  const a = [1, 2, 3]
  const result = vectors.product(a)
  expect(result).toBe(6)
})

it("should compute the mean of a vector", () => {
  const a = [1, 2, 3]
  const result = vectors.mean(a)
  expect(result).toBe(2)
})

it("should compute the min of a vector", () => {
  const a = [1, 2, 3]
  const result = vectors.min(a)
  expect(result).toBe(1)
})

it("should compute the max of a vector", () => {
  const a = [1, 2, 3]
  const result = vectors.max(a)
  expect(result).toBe(3)
})

it("should compute the absolute value of a vector", () => {
  const a = [1, -2, 3]
  const result = vectors.abs(a)
  expect(result).toBe([1, 2, 3])
})

it("should compute the absolute value of a vector in place", () => {
  const a = [1, -2, 3]
  vectors.absInPlace(a)
  expect(a).toBe([1, 2, 3])
})

it("should compute the euclidian distance between two vectors", () => {
  const a = [1, 2, 3]
  const b = [4, 5, 6]
  const result = vectors.euclidianDistance(a, b)
  expect(result).toBe(sqrt<f64>(27))
})

run()
