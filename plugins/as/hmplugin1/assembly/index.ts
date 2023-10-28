import { register } from "hypermode-as/assembly"
import { HMType } from "hypermode-as/protos/HMType"

// TODO: find a way to automatically supply this schema/type info so the user doesn't have to.
register("Query.add", "add", ["a", "b"], [HMType.INT, HMType.INT], HMType.INT);
register("Query.getFullName", "getFullName", ["firstName", "lastName"], [HMType.STRING, HMType.STRING], HMType.STRING);
register("Person.fullName", "getFullName", ["firstName", "lastName"], [HMType.STRING, HMType.STRING], HMType.STRING);

export function add(a: i32, b: i32): i32 {
  return a + b;
}

export function getFullName(firstName: string, lastName: string): string {
  return `${firstName} ${lastName}`;
}
