/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

const ADJS = ["stealth", "covert", "master", "rogue", "steady", "secret", "bold", "infinite", "true", "ops"];
const NOUNS = ["tactic", "script", "sequence", "nexus", "blueprint", "protocol", "strategy", "path", "circuit", "code"];

export function generateAppName(): string {
  const randomAdjective = getRandomItem(ADJS);
  const randomNoun = getRandomItem(NOUNS);

  return `${randomAdjective}-${randomNoun}`;
}

function getRandomItem<T>(list: T[]): T {
  return list[Math.floor(Math.random() * list.length)];
}
