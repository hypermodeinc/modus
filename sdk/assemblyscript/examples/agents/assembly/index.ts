/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { CounterAgent } from "./counter";

const a1 = new CounterAgent();
const a2 = new CounterAgent();

export function updateCount1(): i32 {
  const count = a1.sendMessage("increment");
  if (count == null) {
    return 0;
  }
  return i32.parse(count);
}

export function updateCount2(): i32 {
  const count = a2.sendMessage("increment");
  if (count == null) {
    return 0;
  }
  return i32.parse(count);
}

export function getCount1(): i32 {
  const count = a1.sendMessage("count");
  if (count == null) {
    return 0;
  }
  return i32.parse(count);
}

export function getCount2(): i32 {
  const count = a2.sendMessage("count");
  if (count == null) {
    return 0;
  }
  return i32.parse(count);
}
