import { JSON } from "json-as/assembly";

export function add(a: i32, b: i32): i32 {
  return a + b;
}

export function getFullName(firstName: string, lastName: string): string {
  return `${firstName} ${lastName}`;
}


// This would be ideal, as it returns an array of Person objects.
// However, it doesn't work because we don't yet know how to read the objects on the host.
// export function getPeople(): Person[] {
//   return [
//     new Person("Bob", "Smith"),
//     new Person("Alice", "Jones")
//   ];
// }

export function getPeople(): string {
  const people = [
    new Person("Bob", "Smith"),
    new Person("Alice", "Jones")
  ];

  // For now, we have to return complex data as JSON.
  return JSON.stringify(people);
}

// @ts-ignore
@json
class Person {
  constructor(public firstName: string, public lastName: string) {}

  // TODO: This works in AssemblyScript, but it doesn't return in the JSON,
  // and Dgraph doesn't resolve back to a new function call either.
  get fullName(): string {
    return `${this.firstName} ${this.lastName}`;
  }
};
