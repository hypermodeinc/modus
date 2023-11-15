import { JSON } from "json-as/assembly";
import { queryDQL } from "./hypermode";

export function add(a: i32, b: i32): i32 {
  return a + b;
}

export function getFullName(firstName: string, lastName: string): string {
  return `${firstName} ${lastName}`;
}

export function getPeople(): string {

  const people = [
    Person.Create("Bob", "Smith"),
    Person.Create("Alice", "Jones"),
  ];

  // Non-scalar values must be returned as JSON.
  return JSON.stringify(people);
}

export function queryPeople(): string {
  const results = queryDQL(`
    {
      people(func: type(Person)) {
        id: uid
        firstName: Person.firstName
        lastName: Person.lastName
      }
    }
  `);

  const data = JSON.parse<PersonQueryResponse>(results);
  data.people.forEach(p => {
    p.fullName = `${p.firstName} ${p.lastName}`;
  });

  return JSON.stringify(data.people);
}

// @ts-ignore
@json
class PersonQueryResponse {
  people!: Person[];
}

// @ts-ignore
@json
class Person {

  id: string | null = null;
  firstName: string = "";
  lastName: string = "";
  fullName: string | null = null;

  static Create(firstName: string, lastName: string): Person {
    const p = new Person();
    p.firstName = firstName;
    p.lastName = lastName;
    return p;
  }
};
