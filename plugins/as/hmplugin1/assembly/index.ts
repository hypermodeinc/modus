import { JSON } from "json-as";
import { dql, graphql } from "hypermode-as";

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

  const response = dql.query<PeopleData>(query);
  const people = response.data.people;
  people.forEach(p => p.updateFullName());
  return JSON.stringify(people);
}

export function queryPeople2(): string {
  const statement = `
    query {
      people: queryPerson {
        id
        firstName
        lastName
        fullName
      }
    }
  `;
  
  const results = graphql.execute<PeopleData>(statement);

  // completely optional, but let's log some tracing info
  const tracing = results.extensions!.tracing;
  const duration = tracing.duration / 1000000.0;
  console.log(`Start: ${tracing.startTime.toISOString()}`);
  console.log(`End: ${tracing.endTime.toISOString()}`);
  console.log(`Duration: ${duration}ms`);

  return JSON.stringify(results.data.people);
}

export function newPerson1(firstName: string, lastName: string): string {
  const statement = `
    {
      set {
        _:x <Person.firstName> "${firstName}" .
        _:x <Person.lastName> "${lastName}" .
        _:x <dgraph.type> "Person" .
      }
    }
  `;

  const response = dql.mutate(statement);
  return response.data.uids.get("x")!;
}

export function newPerson2(firstName: string, lastName: string): string {
  const statement = `
    mutation {
      addPerson(input: [{firstName: "${firstName}", lastName: "${lastName}" }]) {
        people: person {
          id
        }
      }
    }
  `;

  const response = graphql.execute<AddPersonPayload>(statement);
  return response.data.addPerson.people[0].id!; 
}

function getPersonCount(): i32 {
  const statement = `
    query {
      aggregatePerson {
        count
      }
    }
  `;

  const response = graphql.execute<AggregatePersonResult>(statement);
  return response.data.aggregatePerson.count;
}

export function getRandomPerson(): string {
  const count = getPersonCount();
  const offset = <u32> Math.floor(Math.random() * count);
  const statement = `
    query {
      people: queryPerson(first: 1, offset: ${offset}) {
        id
        firstName
        lastName
        fullName
      }
    }
  `;
  
  const results = graphql.execute<PeopleData>(statement);
  return JSON.stringify(results.data.people[0]);
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
}

// @ts-ignore
@json
class PeopleData {
  people!: Person[];
}

// @ts-ignore
@json
class AddPersonPayload {
  addPerson!: PeopleData;
}

// @ts-ignore
@json
class AggregatePersonResult {
  aggregatePerson!: GQLAggregateValues;
}

// @ts-ignore
@json
class GQLAggregateValues {
  count: u32 = 0;
}

export function testError(): void {
  console.log("hello");
  console.log("howdy");
  throw new Error("This is a test error");
}
