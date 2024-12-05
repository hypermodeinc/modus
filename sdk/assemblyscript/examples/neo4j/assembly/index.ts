/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { neo4j } from "@hypermode/modus-sdk-as";

// This host name should match one defined in the modus.json manifest file.
const hostName: string = "my-database";

class Person {
  name: string;
  age: number;
  friends: string[] | null;

  constructor(name: string, age: number, friends: string[] | null = null) {
    this.name = name;
    this.age = age;
    this.friends = friends;
  }
}

export function CreatePeopleAndRelationships(): string {
  const person = new neo4j.Variables();

  const people: Person[] = [
    new Person("Alice", 42, ["Bob", "Peter", "Anna"]),
    new Person("Bob", 19),
    new Person("Peter", 50),
    new Person("Anna", 30),
  ];

  for (let i = 0; i < people.length; i++) {
    person.set<string>("name", people[i].name);
    person.set<number>("age", people[i].age);
    if (people[i].friends) {
      person.set<string[]>("friends", people[i].friends!);
    }
    const createPersonQuery = `
      MATCH (p:Person {name: $person.name})
                UNWIND $person.friends AS friend_name
                MATCH (friend:Person {name: friend_name})
                MERGE (p)-[:KNOWS]->(friend)
    `;
    const result = neo4j.executeQuery(hostName, createPersonQuery, person);
    if (!result) {
      throw new Error("Error creating person.");
    }
  }

  return "People and relationships created successfully";
}

export function GetAliceFriendsUnder40(): neo4j.Node[] {
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

  const nodes: neo4j.Node[] = [];

  for (let i = 0; i < result.Records.length; i++) {
    const record = result.Records[i];
    console.log(record.get("friend"));
    const node = neo4j.getRecordValue<neo4j.Node>(record, "friend");
    nodes.push(node);
  }

  return nodes;
}
