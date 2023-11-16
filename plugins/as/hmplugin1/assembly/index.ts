import { JSON } from "json-as/assembly";
import { JSON as DynJSON } from "assemblyscript-json/assembly";
import { queryDQL, queryGQL, GQLExtensions } from "./hypermode";

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
  const results = queryDQL(`
    {
      people(func: type(Person)) {
        id: uid
        firstName: Person.firstName
        lastName: Person.lastName
      }
    }
  `);

  const data = JSON.parse<PeopleData>(results);
  data.people.forEach(p => {
    p.fullName = `${p.firstName} ${p.lastName}`;
  });

  return JSON.stringify(data.people);
}

export function queryPeople2(): string {
  const results = queryGQL(`
    {
      people: queryPerson {
        id
        firstName
        lastName
        fullName
      }
    }
  `);

  // This works, but we get no concrete objects:
  // const resp = <DynJSON.Obj> DynJSON.parse(results);
  // return resp.getObj("data")!.getArr("people")!.stringify();

  // This is better, but still not great:
  const resp = JSON.parse<PeopleGQLResponse>(results);

  // Ideally, we'd like to do this:
  // const resp = JSON.parse<GQLResponse<PeopleData>>(results);
  // but we're blocked by https://github.com/JairusSW/as-json/issues/53
  
  const people = resp.data.people;
  const duration = resp.extensions!.tracing.duration / 1000000.0;
  console.log(`Duration: ${duration}ms`);
  return JSON.stringify(people);
}

// @ts-ignore
@json
class PeopleGQLResponse {
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
