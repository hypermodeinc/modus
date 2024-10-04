/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { graphql } from "@hypermode/modus-sdk-as";
import {
  AddPersonPayload,
  AggregatePersonResult,
  PeopleData,
  Person,
} from "./classes";

// This host name should match one defined in the hypermode.json manifest file.
const hostName: string = "dgraph";

// This function returns the results of querying for all people in the database.
export function queryPeople(): Person[] | null {
  const statement = `
    query {
      people: queryPerson {
        id
        firstName
        lastName
      }
    }
  `;

  const response = graphql.execute<PeopleData>(hostName, statement);
  if (!response.data) return [];
  return response.data!.people;
}

// This function returns the results of querying for a specific person in the database.
export function querySpecificPerson(
  firstName: string,
  lastName: string,
): Person | null {
  const statement = `
    query queryPeople($firstName: String!, $lastName: String!) {
      people: queryPerson(
          first: 1,
          filter: { firstName: { eq: $firstName }, lastName: { eq: $lastName } }
      ) {
          id
          firstName
          lastName
      }
    }
  `;

  const vars = new graphql.Variables();
  vars.set("firstName", firstName);
  vars.set("lastName", lastName);

  const response = graphql.execute<PeopleData>(hostName, statement, vars);

  if (!response.data) return null;
  const people = response.data!.people;
  if (people.length === 0) return null;
  return people[0];
}

// This function returns the results of querying for the count of people in the database.
function getPersonCount(): i32 {
  const statement = `
    query {
      aggregatePerson {
        count
      }
    }
  `;

  const response = graphql.execute<AggregatePersonResult>(hostName, statement);
  return response.data!.aggregatePerson.count;
}

// This function returns the results of querying for a random person in the database.
export function getRandomPerson(): Person | null {
  const count = getPersonCount();
  const offset = <u32>Math.floor(Math.random() * count);
  const statement = `
    query {
      people: queryPerson(first: 1, offset: ${offset}) {
        id
        firstName
        lastName
      }
    }
  `;

  const response = graphql.execute<PeopleData>(hostName, statement);
  if (!response.data) return null;
  const people = response.data!.people;
  if (people.length === 0) return null;
  return people[0];
}

// This function adds a new person to the database and returns that person.
export function addPerson(firstName: string, lastName: string): Person {
  const statement = `
    mutation {
      addPerson(input: [{firstName: "${firstName}", lastName: "${lastName}" }]) {
        people: person {
          id
          firstName
          lastName
        }
      }
    }
  `;

  const response = graphql.execute<AddPersonPayload>(hostName, statement);
  return response.data!.addPerson.people[0];
}

// This function demonstrates what happens when a bad query is executed.
export function testBadQuery(): Person[] {
  const statement = "this is a bad query";
  const results = graphql.execute<PeopleData>(hostName, statement);
  if (!results.data) return [];
  return results.data!.people;
}
