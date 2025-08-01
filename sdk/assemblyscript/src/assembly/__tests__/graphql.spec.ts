/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { expect, it, log, mockFn, mockImport, run } from "as-test";
import { graphql } from "..";
import { JSON } from "json-as";

let returnData: string = "";
mockImport(
  "modus_graphql_client.executeQuery",
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  (connection: string, statement: string, variables: string): string => {
    return returnData;
  },
);

mockFn(console.log, (data: string): void => {
  log(data);
});

it("should execute graphql query", () => {
  const statement = `
    query {
      people: queryPerson {
        id
        firstName
        lastName
      }
    }
  `;

  returnData = '{"data":{"people":[]}}';

  const response = graphql.execute<PeopleData>("dgraph", statement);
  expect(!response.data).toBe(false);
  expect(!response.data!.people).toBe(false);
});

it("should query people", () => {
  const query = `
    query {
      people: queryPerson {
        id
        firstName
        lastName
      }
    }
  `;

  const _person: Person = {
    id: "0xb8",
    firstName: "Jairus",
    lastName: "Tanaka",
  };

  log("Person: " + JSON.stringify(_person));

  returnData = '{"data":{"people":[' + JSON.stringify(_person) + "]}}";

  const response = graphql.execute<PeopleData>("dgraph", query);
  expect(!response.data).toBe(false);
  expect(!response.data!.people).toBe(false);

  const person = response.data!.people[0];
  expect(person.id).toBe("0xb8");
  expect(person.firstName).toBe("Jairus");
  expect(person.lastName).toBe("Tanaka");
});

it("should query people with variables", () => {
  const query = `
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
  vars.set("firstName", "Jairus");
  vars.set("lastName", "Tanaka");

  const response = graphql.execute<PeopleData>("dgraph", query, vars);
  expect(!response.data).toBe(false);
  expect(!response.data!.people).toBe(false);

  const people = response.data!.people;
  expect(people.length).toBeGreaterThan(0);
});

run();


@json
class Person {
  id: string | null = null;
  firstName: string = "";
  lastName: string = "";
}


@json
class PeopleData {
  people: Person[] = [];
}
