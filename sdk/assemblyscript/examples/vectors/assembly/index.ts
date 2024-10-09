/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { vectors } from "@hypermode/modus-sdk-as";

export function add(a: f64[], b: f64[]): f64[] {
  return vectors.add(a, b);
}

export function addInPlace(a: f64[], b: f64[]): f64[] {
  vectors.addInPlace(a, b);
  return a;
}

export function subtract(a: f64[], b: f64[]): f64[] {
  return vectors.subtract(a, b);
}

export function subtractInPlace(a: f64[], b: f64[]): f64[] {
  vectors.subtractInPlace(a, b);
  return a;
}

export function addNumber(a: f64[], b: f64): f64[] {
  return vectors.addNumber(a, b);
}

export function addNumberInPlace(a: f64[], b: f64): f64[] {
  vectors.addNumberInPlace(a, b);
  return a;
}

export function subtractNumber(a: f64[], b: f64): f64[] {
  return vectors.subtractNumber(a, b);
}

export function subtractNumberInPlace(a: f64[], b: f64): f64[] {
  vectors.subtractNumberInPlace(a, b);
  return a;
}

export function multiplyNumber(a: f64[], b: f64): f64[] {
  return vectors.multiplyNumber(a, b);
}

export function multiplyNumberInPlace(a: f64[], b: f64): f64[] {
  vectors.multiplyNumberInPlace(a, b);
  return a;
}

export function divideNumber(a: f64[], b: f64): f64[] {
  return vectors.divideNumber(a, b);
}

export function divideNumberInPlace(a: f64[], b: f64): f64[] {
  vectors.divideNumberInPlace(a, b);
  return a;
}

export function dot(a: f64[], b: f64[]): f64 {
  return vectors.dot(a, b);
}

export function magnitude(a: f64[]): f64 {
  return vectors.magnitude(a);
}

export function normalize(a: f64[]): f64[] {
  return vectors.normalize(a);
}

export function sum(a: f64[]): f64 {
  return vectors.sum(a);
}

export function product(a: f64[]): f64 {
  return vectors.product(a);
}

export function mean(a: f64[]): f64 {
  return vectors.mean(a);
}

export function min(a: f64[]): f64 {
  return vectors.min(a);
}

export function max(a: f64[]): f64 {
  return vectors.max(a);
}

export function abs(a: f64[]): f64[] {
  return vectors.abs(a);
}

export function absInPlace(a: f64[]): f64[] {
  vectors.absInPlace(a);
  return a;
}

export function euclidianDistance(a: f64[], b: f64[]): f64 {
  return vectors.euclidianDistance(a, b);
}
