import { JSON } from "json-as";
import { queryDQL, queryGQL } from "hypermode-as";

export function add(a: i32, b: i32): i32 {
  return a + b;
}

export function getFullName(firstName: string, lastName: string): string {
  return `${firstName} ${lastName}`;
}

export function getPeople(): string {
  const people = [
    <Person> {firstName: "Bob", lastName: "Smith"},
    <Person> {firstName: "Alice", lastName: "Jones"}
  ];

  people.forEach(p => p.updateFullName());
  return JSON.stringify(people);
}

export function queryPeople1(): string {
  const query = `
    {
      people(func: type(Person)) {
        id: uid
        firstName: Person.firstName
        lastName: Person.lastName
      }
    }
  `;

  const people = queryDQL<PeopleData>(query).people
  people.forEach(p => p.updateFullName());
  return JSON.stringify(people);
}

export function queryPeople2(): string {
  const query = `
    {
      people: queryPerson {
        id
        firstName
        lastName
        fullName
      }
    }
  `;
  
  const results = queryGQL<PeopleData>(query);

  // completely optional, but let's log some tracing info
  const tracing = results.extensions!.tracing;
  const duration = tracing.duration / 1000000.0;
  console.log(`Start: ${tracing.startTime.toISOString()}`);
  console.log(`End: ${tracing.endTime.toISOString()}`);
  console.log(`Duration: ${duration}ms`);

  return JSON.stringify(results.data.people);
}

// @ts-ignore
@json
class PeopleData {
  people!: Person[];
}

// @ts-ignore
@json
class Person {
  id: string | null = null;
  firstName: string = "";
  lastName: string = "";
  fullName: string | null = null;

  updateFullName(): void {
    this.fullName = `${this.firstName} ${this.lastName}`;
  }
};
