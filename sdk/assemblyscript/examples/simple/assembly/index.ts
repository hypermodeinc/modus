/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { Person } from "./person";

// This function adds two 32-bit signed integers together, and returns the result.
export function add(a: i32, b: i32): i32 {
  return a + b;
}

// This function takes a first name and a last name, and concatenates them to returns a full name.
export function getFullName(firstName: string, lastName: string): string {
  return `${firstName} ${lastName}`;
}

// This function makes a list of people, and returns it.
export function getPeople(): Person[] {
  return [
    new Person("Bob", "Smith"),
    new Person("Alice", "Jones"),
    new Person("Charlie", "Brown"),
  ];
}

// This function returns a random person from the list of people.
export function getRandomPerson(): Person {
  const people = getPeople();
  const index = <i32>Math.floor(Math.random() * people.length);
  const person = people[index];
  return person;
}

// This function demonstrates various ways to log messages and errors.
export function testErrors(): void {
  // This is a simple log message. It has no level.
  console.log("This is a simple log message.");

  // These messages are logged at different levels.
  console.debug("This is a debug message.");
  console.info("This is an info message.");
  console.warn("This is a warning message.");

  // This logs an error message, but allows the function to continue.
  console.error("This is an error message.");
  console.error(
    `This is line 1 of a multi-line error message.
This is line 2 of a multi-line error message.
This is line 3 of a multi-line error message.`,
  );

  // This throws an error, which will log the message as "fatal" and exit the function.
  throw new Error(
    `This is a message from a thrown error.
This is a second line from a thrown error.`,
  );
}
