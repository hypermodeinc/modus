/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { modusdb } from "@hypermode/modus-sdk-as";
import { PeopleData, Person, Plugin, PluginData } from "./classes";
import { JSON } from "json-as";
import { MutationRequest } from "@hypermode/modus-sdk-as/assembly/modusdb";

export function dropAll(): string {
  return modusdb.dropAll();
}

export function dropData(): string {
  return modusdb.dropData();
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
  return modusdb.alterSchema(schema);
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

  const resp = modusdb.query(query);

  return JSON.parse<PeopleData>(resp.Json).people;
}

export function queryPlugins(): Plugin[] {
  const query = `
  {
	  plugins(func: type(Plugin)) {
		uid
		id
		name
		version
		language
		sdk_version
		build_id
		build_time
		git_repo
		git_commit
		dgraph.type
	  }
	}
  `;

  const resp = modusdb.query(query);

  return JSON.parse<PluginData>(resp.Json).plugins;
}

// This function returns the results of querying for a specific person in the database.
export function querySpecificPerson(
  firstName: string,
  lastName: string,
): Person | null {
  const query = `
  query queryPerson {
    people(func: eq(firstName, "${firstName}")) @filter(eq(lastName, "${lastName}")) {
        uid
        firstName
        lastName
    }
}
  `;

  const resp = modusdb.query(query);

  const people = JSON.parse<PeopleData>(resp.Json).people;

  if (people.length === 0) return null;
  return people[0];
}

// This function adds a new person to the database and returns that person.
export function addPersonAsRDF(
  firstName: string,
  lastName: string,
): Map<string, u64> | null {
  const mutation = `
  _:user1 <firstName> "${firstName}" .
  _:user1 <lastName> "${lastName}" .
  _:user1 <dgraph.type> "Person" .
  `;

  const mutations: modusdb.Mutation[] = [
    new modusdb.Mutation("", "", mutation),
  ];

  return modusdb.mutate(new MutationRequest(mutations));
}

export function addPersonAsJSON(
  firstName: string,
  lastName: string,
): Map<string, u64> | null {
  const person = new Person("_:user1", firstName, lastName, ["Person"]);

  const mutation = JSON.stringify(person);

  const mutations: modusdb.Mutation[] = [new modusdb.Mutation(mutation)];

  return modusdb.mutate(new MutationRequest(mutations));
}

export function deletePerson(uid: string): Map<string, u64> | null {
  const mutation = `<${uid}> * * .`;

  const mutations: modusdb.Mutation[] = [
    new modusdb.Mutation("", "", mutation),
  ];

  return modusdb.mutate(new MutationRequest(mutations));
}

// This function demonstrates what happens when a bad query is executed.
export function testBadQuery(): Person[] {
  const query = "this is a bad query";

  const resp = modusdb.query(query);

  return JSON.parse<PeopleData>(resp.Json).people;
}
