/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { mysql } from "@hypermode/modus-sdk-as";

// The name of the MySQL host, as specified in the modus.json manifest
const host = "my-database";


@json
class Person {
  id: i32 = 0;
  name!: string;
  age!: i32;
  home!: mysql.Location | null;
}

export function getAllPeople(): Person[] {
  const query = "select * from people order by id";
  const response = mysql.query<Person>(host, query);
  return response.rows;
}

export function getPeopleByName(name: string): Person[] {
  const query = "select * from people where name = ?";

  const params = new mysql.Params();
  params.push(name);

  const response = mysql.query<Person>(host, query, params);
  return response.rows;
}

export function getPerson(id: i32): Person | null {
  const query = "select * from people where id = ?";

  const params = new mysql.Params();
  params.push(id);

  const response = mysql.query<Person>(host, query, params);
  return response.rows.length > 0 ? response.rows[0] : null;
}

export function addPerson(name: string, age: i32): Person {
  const query = "insert into people (name, age) values (?, ?)";

  const params = new mysql.Params();
  params.push(name);
  params.push(age);

  const response = mysql.execute(host, query, params);

  if (response.rowsAffected != 1) {
    throw new Error("Failed to insert person.");
  }

  const id = <i32>response.lastInsertId;
  return <Person>{ id, name, age };
}

export function updatePersonHome(
  id: i32,
  longitude: f64,
  latitude: f64,
): Person | null {
  const query = `update people set home = point(?,?) where id = ?`;

  const params = new mysql.Params();
  params.push(longitude);
  params.push(latitude);
  params.push(id);

  const response = mysql.execute(host, query, params);

  if (response.rowsAffected != 1) {
    console.error(
      `Failed to update person with id ${id}. The record may not exist.`,
    );
    return null;
  }

  return getPerson(id);
}

export function deletePerson(id: i32): string {
  const query = "delete from people where id = ?";

  const params = new mysql.Params();
  params.push(id);

  const response = mysql.execute(host, query, params);

  if (response.rowsAffected != 1) {
    console.error(
      `Failed to delete person with id ${id}. The record may not exist.`,
    );
    return "failure";
  }

  return "success";
}
