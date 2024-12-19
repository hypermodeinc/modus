/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { neo4j } from "@hypermode/modus-sdk-as";

// This host name should match one defined in the modus.json manifest file.
const hostName: string = "my-database";


@json
class Person {
  name: string;
  age: i32;
  friends: string[] | null;

  constructor(name: string, age: i32, friends: string[] | null = null) {
    this.name = name;
    this.age = age;
    this.friends = friends;
  }
}

export function CreatePeopleAndRelationships(): string {
  const people: Person[] = [
    new Person("Alice", 42, ["Bob", "Peter", "Anna"]),
    new Person("Bob", 19),
    new Person("Peter", 50),
    new Person("Anna", 30),
  ];

  for (let i = 0; i < people.length; i++) {
    const createPersonQuery = `
      MATCH (p:Person {name: $person.name})
      UNWIND $person.friends AS friend_name
      MATCH (friend:Person {name: friend_name})
      MERGE (p)-[:KNOWS]->(friend)
    `;
    const peopleVars = new neo4j.Variables();
    peopleVars.set("person", people[i]);
    const result = neo4j.executeQuery(hostName, createPersonQuery, peopleVars);
    if (!result) {
      throw new Error("Error creating person.");
    }
  }

  return "People and relationships created successfully";
}

export function GetAliceFriendsUnder40(): Person[] {
  const vars = new neo4j.Variables();
  vars.set("name", "Alice");
  vars.set("age", 40);

  const query = `
    MATCH (p:Person {name: $name})-[:KNOWS]-(friend:Person)
        WHERE friend.age < $age
        RETURN friend
  `;

  const result = neo4j.executeQuery(hostName, query, vars);
  if (!result) {
    throw new Error("Error getting friends.");
  }

  const personNodes: Person[] = [];

  for (let i = 0; i < result.Records.length; i++) {
    const record = result.Records[i];
    const node = record.getValue<neo4j.Node>("friend");
    const person = new Person(
      node.getProperty<string>("name"),
      node.getProperty<i32>("age"),
    );
    personNodes.push(person);
  }

  return personNodes;
}

export function GetAliceFriendsUnder40Ages(): i32[] {
  const vars = new neo4j.Variables();
  vars.set("name", "Alice");
  vars.set("age", 40);

  const query = `
    MATCH (p:Person {name: $name})-[:KNOWS]-(friend:Person)
        WHERE friend.age < $age
        RETURN friend.age AS age
  `;

  const result = neo4j.executeQuery(hostName, query, vars);
  if (!result) {
    throw new Error("Error getting friends.");
  }

  const ages: i32[] = [];

  for (let i = 0; i < result.Records.length; i++) {
    const record = result.Records[i];
    const age = record.getValue<i32>("age");
    ages.push(age);
  }

  return ages;
}

export function DeleteAllNodes(): string {
  const query = `
  MATCH (n)
  DETACH DELETE n`;
  const result = neo4j.executeQuery(hostName, query);
  if (!result) {
    throw new Error("Error deleting nodes.");
  }

  return "All nodes deleted successfully";
}
