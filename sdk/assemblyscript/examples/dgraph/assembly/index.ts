/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { dgraph } from "@hypermode/modus-sdk-as";
import { PeopleData, Person } from "./classes";
import { JSON } from "json-as";

// This connection name should match one defined in the modus.json manifest file.
const connection: string = "dgraph";

export function dropAll(): string {
  return dgraph.dropAll(connection);
}

export function dropAttr(attr: string): string {
  return dgraph.dropAttr(connection, attr);
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
  return dgraph.alterSchema(connection, schema);
}

// This function returns the results of querying for all people in the database.
export function queryPeople(): Person[] {
  const query = new dgraph.Query(`
    {
      people(func: type(Person)) {
        uid
        firstName
        lastName
      }
    }
    `);

  const response = dgraph.executeQuery(connection, query);

  return JSON.parse<PeopleData>(response.Json).people;
}

// This function returns the results of querying for a specific person in the database.
export function querySpecificPerson(
  firstName: string,
  lastName: string,
): Person | null {
  const query = new dgraph.Query(`
    query queryPerson($firstName: string, $lastName: string) {
      people(func: eq(firstName, $firstName)) @filter(eq(lastName, $lastName)) {
          uid
          firstName
          lastName
      }
    }
    `)
    .withVariable("$firstName", firstName)
    .withVariable("$lastName", lastName);

  const response = dgraph.executeQuery(connection, query);
  const people = JSON.parse<PeopleData>(response.Json).people;

  return people.length == 0 ? null : people[0];
}

// This function adds a new person to the database and returns that person.
export function addPersonAsRDF(
  firstName: string,
  lastName: string,
): Map<string, string> | null {
  const mutation = new dgraph.Mutation().withSetNquads(`
    _:user1 <firstName> "${dgraph.escapeRDF(firstName)}" .
    _:user1 <lastName> "${dgraph.escapeRDF(lastName)}" .
    _:user1 <dgraph.type> "Person" .
    `);

  const response = dgraph.executeMutations(connection, mutation);
  return response.Uids;
}

export function addPersonAsJSON(
  firstName: string,
  lastName: string,
): Map<string, string> | null {
  const person = new Person("_:user1", firstName, lastName, ["Person"]);

  const personJson = JSON.stringify(person);
  const mutation = new dgraph.Mutation().withSetJson(personJson);

  const response = dgraph.executeMutations(connection, mutation);
  return response.Uids;
}

export function updatePerson(
  nameToChangeFrom: string,
  nameToChangeTo: string,
): Map<string, string> | null {
  const query = new dgraph.Query(`
    query findPerson($name: string) {
      person as var(func: eq(firstName, $name))
    }
    `).withVariable("$name", nameToChangeFrom);

  const mutation = new dgraph.Mutation().withSetNquads(
    `uid(person) <firstName> "${dgraph.escapeRDF(nameToChangeTo)}" .`,
  );

  const response = dgraph.executeQuery(connection, query, mutation);
  return response.Uids;
}

export function deletePerson(uid: string): Map<string, string> | null {
  const mutation = new dgraph.Mutation().withDelNquads(`<${uid}> * * .`);
  const response = dgraph.executeMutations(connection, mutation);
  return response.Uids;
}

// This function demonstrates what happens when a bad query is executed.
export function testBadQuery(): Person[] {
  const query = new dgraph.Query("this is a bad query");
  const response = dgraph.executeQuery(connection, query);

  return JSON.parse<PeopleData>(response.Json).people;
}
