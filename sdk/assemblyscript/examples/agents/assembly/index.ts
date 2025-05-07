/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { CounterAgent } from "./counter";

const ca = new CounterAgent();

export function updateCount(): i32 {
  const count = ca.sendMessage("increment");
  if (count == null) {
    return 0;
  }
  return i32.parse(count);
}

export function getCount(): i32 {
  const count = ca.sendMessage("count");
  if (count == null) {
    return 0;
  }
  return i32.parse(count);
}
