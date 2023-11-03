import { JSON } from "json-as/assembly";

export function add(a: i32, b: i32): i32 {
  return a + b;
}

export function getFullName(firstName: string, lastName: string): string {
  return `${firstName} ${lastName}`;
}

export function getPeople(): string {
  const people = [
    new Person("Bob", "Smith"),
    new Person("Alice", "Jones")
  ];

  // Non-scalar values must be returned as JSON.
  return JSON.stringify(people);
}

// @ts-ignore
@json
class Person {
  constructor(firstName: string, lastName: string) {
    this.firstName = firstName;
    this.lastName = lastName;
    this.fullName = getFullName(firstName, lastName);
  }
  
  firstName: string;
  lastName: string;
  fullName: string;
};
