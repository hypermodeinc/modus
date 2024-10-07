/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

export function checkEqualLength<T extends number>(a: T[], b: T[]): void {
  if (a.length !== b.length) {
    throw new Error("Vectors must be the same length.");
  }
}

/**
 *
 * Add two vectors together, returning a new vector.
 * @param a: The first vector
 * @param b: The second vector
 * @returns: The sum of the two vectors
 */
export function add<T extends number>(a: T[], b: T[]): T[] {
  checkEqualLength(a, b);
  const result = new Array<T>(a.length);
  for (let i = 0; i < a.length; i++) {
    result[i] = (a[i] + b[i]) as T;
  }
  return result;
}

/**
 *
 * Add two vectors together, modifying the first vector.
 * @param a: The first vector
 * @param b: The second vector
 */
export function addInplace<T extends number>(a: T[], b: T[]): void {
  checkEqualLength(a, b);
  for (let i = 0; i < a.length; i++) {
    a[i] = (a[i] + b[i]) as T;
  }
}

/**
 *
 * Subtract two vectors, returning a new vector.
 * @param a: The first vector
 * @param b: The second vector
 * @returns: The difference of the two vectors
 */
export function subtract<T extends number>(a: T[], b: T[]): T[] {
  checkEqualLength(a, b);
  const result = new Array<T>(a.length);
  for (let i = 0; i < a.length; i++) {
    result[i] = (a[i] - b[i]) as T;
  }
  return result;
}

/**
 *
 * Subtract two vectors, modifying the first vector.
 * @param a: The first vector
 * @param b: The second vector
 */
export function subtractInplace<T extends number>(a: T[], b: T[]): void {
  checkEqualLength(a, b);
  for (let i = 0; i < a.length; i++) {
    a[i] = (a[i] - b[i]) as T;
  }
}

/**
 *
 * add a number to a vector, returning a new vector.
 * @param a: The first vector
 * @param b: The number to add
 * @returns: the result vector, with the number added to each element
 */
export function addNumber<T extends number>(a: T[], b: T): T[] {
  const result = new Array<T>(a.length);
  for (let i = 0; i < a.length; i++) {
    result[i] = (a[i] + b) as T;
  }
  return result;
}

/**
 *
 * add a number to a vector, modifying the first vector.
 * @param a: The first vector
 * @param b: The number to add
 */
export function addNumberInplace<T extends number>(a: T[], b: T): void {
  for (let i = 0; i < a.length; i++) {
    a[i] = (a[i] + b) as T;
  }
}

/**
 *
 * Subtract a number from a vector, returning a new vector.
 * @param a: The first vector
 * @param b: The number to subtract
 * @returns: the result vector, with the number subtracted from each element
 */
export function subtractNumber<T extends number>(a: T[], b: T): T[] {
  const result = new Array<T>(a.length);
  for (let i = 0; i < a.length; i++) {
    result[i] = (a[i] - b) as T;
  }
  return result;
}

/**
 *
 * Subtract a number from a vector, modifying the first vector.
 * @param a: The first vector
 * @param b: The number to subtract
 */
export function subtractNumberInplace<T extends number>(a: T[], b: T): void {
  for (let i = 0; i < a.length; i++) {
    a[i] = (a[i] - b) as T;
  }
}

/**
 * Multiple numbers to a vector, returning a new vector.
 * @param a: The first vector
 * @param b: The number to multiply
 * @returns: the result vector, with the number multiplied to each element
 */
export function multiplyNumber<T extends number>(a: T[], b: T): T[] {
  const result = new Array<T>(a.length);
  for (let i = 0; i < a.length; i++) {
    result[i] = (a[i] * b) as T;
  }
  return result;
}

/**
 *
 * Multiply a number to a vector, modifying the first vector.
 * @param a: The first vector
 * @param b: The number to multiply
 */
export function multiplyNumberInplace<T extends number>(a: T[], b: T): void {
  for (let i = 0; i < a.length; i++) {
    a[i] = (a[i] * b) as T;
  }
}

/**
 *
 * Divide a number from a vector, returning a new vector.
 * @param a: The first vector
 * @param b: The number to divide
 * @returns: the result vector, with the number divided from each element
 */
export function divideNumber<T extends number>(a: T[], b: T): T[] {
  const result = new Array<T>(a.length);
  for (let i = 0; i < a.length; i++) {
    result[i] = (a[i] / b) as T;
  }
  return result;
}

/**
 *
 * Divide a number from a vector, modifying the first vector.
 * @param a: The first vector
 * @param b: The number to divide
 */
export function divideNumberInplace<T extends number>(a: T[], b: T): void {
  for (let i = 0; i < a.length; i++) {
    a[i] = (a[i] / b) as T;
  }
}

/**
 * Calculate the dot product of two vectors.
 * @param a: The first vector
 * @param b: The second vector
 * @returns: The dot product of the two vectors
 */
export function dot<T extends number>(a: T[], b: T[]): T {
  checkEqualLength(a, b);
  let result = 0;
  for (let i = 0; i < a.length; i++) {
    result += a[i] * b[i];
  }
  return result as T;
}

/**
 * Calculate the magnitude of a vector.
 * @param a: The vector
 * @returns: The magnitude of the vector
 */
export function magnitude<T extends number>(a: T[]): f64 {
  return sqrt<f64>(dot(a, a));
}

/**
 * Calculate the cross product of two 3D vectors.
 * @param a: The first vector
 * @param b: The second vector
 * @returns: The cross product of the two vectors
 */
export function normalize<T extends number>(a: T[]): T[] {
  return divideNumber(a, magnitude(a)) as T[];
}

/**
 *
 * Normalize a vector, modifying the first vector.
 * @param a: The vector to normalize
 */
export function normalizeInplace<T extends number>(a: T[]): void {
  divideNumberInplace(a, magnitude(a) as T);
}

/**
 *
 * Calculate the sum of a vector.
 * @param a: The vector
 * @returns: The sum of the vector
 */
export function sum<T extends number>(a: T[]): T {
  let result = 0;
  for (let i = 0; i < a.length; i++) {
    result += a[i];
  }
  return result as T;
}

/**
 *
 * Calculate the product of a vector.
 * @param a: The vector
 * @returns: The product of the vector
 */
export function product<T extends number>(a: T[]): T {
  let result = 1;
  for (let i = 0; i < a.length; i++) {
    result *= a[i];
  }
  return result as T;
}

/**
 *
 * Calculate the mean of a vector.
 * @param a: The vector
 * @returns: The mean of the vector
 */
export function mean<T extends number>(a: T[]): f64 {
  return sum(a) / a.length;
}

/**
 *
 * Calculate the median of a vector.
 * @param a: The vector
 * @returns: The median of the vector
 */
export function min<T extends number>(a: T[]): T {
  let result = a[0];
  for (let i = 1; i < a.length; i++) {
    if (a[i] < result) {
      result = a[i];
    }
  }
  return result;
}

/**
 *
 * Calculate the maximum of a vector.
 * @param a: The vector
 * @returns: The maximum of the vector
 */
export function max<T extends number>(a: T[]): T {
  let result = a[0];
  for (let i = 1; i < a.length; i++) {
    if (a[i] > result) {
      result = a[i];
    }
  }
  return result;
}

/**
 *
 * Calculate the absolute value of a vector.
 * @param a: The vector
 * @returns: The absolute value of the vector
 */
export function abs<T extends number>(a: T[]): T[] {
  const result = new Array<T>(a.length);
  for (let i = 0; i < a.length; i++) {
    result[i] = a[i] < 0 ? (-a[i] as T) : a[i];
  }
  return result;
}

/**
 *
 * Calculate the absolute value of a vector, modifying the first vector.
 * @param a: The vector
 */
export function absInplace<T extends number>(a: T[]): void {
  for (let i = 0; i < a.length; i++) {
    a[i] = a[i] < 0 ? (-a[i] as T) : a[i];
  }
}

/**
 *
 * Calculate the euclidian distance between two vectors.
 * @param a: The first vector
 * @param b: The second vector
 * @returns: The euclidian distance between the two vectors
 */
export function euclidianDistance<T extends number>(a: T[], b: T[]): f64 {
  let sum = 0;
  for (let i = 0; i < a.length; i++) {
    sum += (a[i] - b[i]) ** 2;
  }
  return sqrt<f64>(sum);
}
