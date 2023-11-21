import { JSON } from "json-as";
import { queryDQL, queryGQL, GQLExtensions } from "hypermode-as";

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

  const data = queryDQL<PeopleData>(query);
  data.people.forEach(p => {
    p.fullName = `${p.firstName} ${p.lastName}`;
  });

  return JSON.stringify(data.people);
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
  
  const results = queryGQL<PeopleGQLResults>(query);
  const people = results.data.people;
  const tracing = results.extensions!.tracing;
  const duration = tracing.duration / 1000000.0;
  console.log(`Start: ${tracing.startTime.toISOString()}`);
  console.log(`End: ${tracing.endTime.toISOString()}`);
  console.log(`Duration: ${duration}ms`);
  return JSON.stringify(people);
}

// @ts-ignore
@json
class PeopleGQLResults {
    data!: PeopleData;
    extensions: GQLExtensions | null = null;
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

  static Create(firstName: string, lastName: string): Person {
    const p = new Person();
    p.firstName = firstName;
    p.lastName = lastName;
    return p;
  }
};
