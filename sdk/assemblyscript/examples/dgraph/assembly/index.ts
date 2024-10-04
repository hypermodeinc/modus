import { dgraph } from "@hypermode/modus-sdk-as";
import { PeopleData, Person } from "./classes";
import { JSON } from "json-as";

// This host name should match one defined in the hypermode.json manifest file.
const hostName: string = "dgraph";

export function dropAll(): string {
  return dgraph.dropAll(hostName);
}

export function dropAttr(attr: string): string {
  return dgraph.dropAttr(hostName, attr);
}

export function alterSchema(): string {
  const schema = `
  firstName: string @index(term) .
  lastName: string @index(term) .
  dgraph.type: [string] @index(exact) .

  type Person {
      firstName
      lastName
  }
  `;
  return dgraph.alterSchema(hostName, schema);
}

// This function returns the results of querying for all people in the database.
export function queryPeople(): Person[] {
  const query = `
  {
    people(func: type(Person)) {
      uid
      firstName
      lastName
    }
  }
  `;

  const resp = dgraph.execute(
    hostName,
    new dgraph.Request(new dgraph.Query(query)),
  );

  return JSON.parse<PeopleData>(resp.Json).people;
}

// This function returns the results of querying for a specific person in the database.
export function querySpecificPerson(
  firstName: string,
  lastName: string,
): Person | null {
  const statement = `
  query queryPerson($firstName: string, $lastName: string) {
    people(func: eq(firstName, $firstName)) @filter(eq(lastName, $lastName)) {
        uid
        firstName
        lastName
    }
}
  `;

  const vars = new dgraph.Variables();
  vars.set("$firstName", firstName);
  vars.set("$lastName", lastName);

  const resp = dgraph.execute(
    hostName,
    new dgraph.Request(new dgraph.Query(statement, vars)),
  );

  const people = JSON.parse<PeopleData>(resp.Json).people;

  if (people.length === 0) return null;
  return people[0];
}

// This function adds a new person to the database and returns that person.
export function addPersonAsRDF(
  firstName: string,
  lastName: string,
): Map<string, string> | null {
  const mutation = `
  _:user1 <firstName> "${firstName}" .
  _:user1 <lastName> "${lastName}" .
  _:user1 <dgraph.type> "Person" .
  `;

  const mutations: dgraph.Mutation[] = [new dgraph.Mutation("", "", mutation)];

  return dgraph.execute(hostName, new dgraph.Request(null, mutations)).Uids;
}

export function addPersonAsJSON(
  firstName: string,
  lastName: string,
): Map<string, string> | null {
  const person = new Person("_:user1", firstName, lastName, ["Person"]);

  const mutation = JSON.stringify(person);

  const mutations: dgraph.Mutation[] = [new dgraph.Mutation(mutation)];

  return dgraph.execute(hostName, new dgraph.Request(null, mutations)).Uids;
}

export function upsertPerson(
  nameToChangeFrom: string,
  nameToChangeTo: string,
): Map<string, string> | null {
  const query = `
  query {
    person as var(func: eq(firstName, "${nameToChangeFrom}"))
  `;
  const mutation = `
    uid(person) <firstName> "${nameToChangeTo}" .`;

  const dgraphQuery = new dgraph.Query(query);

  const mutationList: dgraph.Mutation[] = [
    new dgraph.Mutation("", "", mutation),
  ];

  const dgraphRequest = new dgraph.Request(dgraphQuery, mutationList);

  const response = dgraph.execute(hostName, dgraphRequest);

  return response.Uids;
}

export function deletePerson(uid: string): Map<string, string> | null {
  const mutation = `<${uid}> * * .`;

  const mutations: dgraph.Mutation[] = [new dgraph.Mutation("", "", mutation)];

  return dgraph.execute(hostName, new dgraph.Request(null, mutations)).Uids;
}

// This function demonstrates what happens when a bad query is executed.
export function testBadQuery(): Person[] {
  const query = "this is a bad query";

  const resp = dgraph.execute(
    hostName,
    new dgraph.Request(new dgraph.Query(query)),
  );

  return JSON.parse<PeopleData>(resp.Json).people;
}
