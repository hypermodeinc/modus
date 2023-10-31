import { register } from "hypermode-as/assembly"
import { HMType } from "hypermode-as/protos/HMType"
import { JSON } from "json-as/assembly";

// TODO: find a way to automatically supply this schema/type info so the user doesn't have to.
register("Query.add", "add", ["a", "b"], [HMType.INT, HMType.INT], HMType.INT);
register("Query.getFullName", "getFullName", ["firstName", "lastName"], [HMType.STRING, HMType.STRING], HMType.STRING);
register("Person.fullName", "getFullName", ["firstName", "lastName"], [HMType.STRING, HMType.STRING], HMType.STRING);
register("Query.getPeople", "getPeople", [], [], HMType.JSON);

export function add(a: i32, b: i32): i32 {
  return a + b;
}

export function getFullName(firstName: string, lastName: string): string {

  // We can use the Person class here.
  const p = new Person(firstName, lastName);
  return p.fullName;

  // This works too, and is more concise.
  // return `${firstName} ${lastName}`;
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
