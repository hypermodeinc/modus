import { postgresql } from "@hypermode/functions-as";

// The name of the PostgreSQL host, as specified in the hypermode.json manifest
const host = "my-database";


@json
class Person {
  id: i32 = 0;
  name!: string;
  age!: i32;
  home!: postgresql.Location | null;
}

export function getAllPeople(): Person[] {
  const query = "select * from people order by id";
  const response = postgresql.query<Person>(host, query);
  return response.rows;
}

export function getPeopleByName(name: string): Person[] {
  const query = "select * from people where name = $1";

  const params = new postgresql.Params();
  params.push(name);

  const response = postgresql.query<Person>(host, query, params);
  return response.rows;
}

export function getPerson(id: i32): Person | null {
  const query = "select * from people where id = $1";

  const params = new postgresql.Params();
  params.push(id);

  const response = postgresql.query<Person>(host, query, params);
  return response.rows.length > 0 ? response.rows[0] : null;
}

export function addPerson(name: string, age: i32): Person {
  const query = `
    insert into people (name, age)
    values ($1, $2) RETURNING id`;

  const params = new postgresql.Params();
  params.push(name);
  params.push(age);

  const response = postgresql.queryScalar<i32>(host, query, params);

  if (response.rowsAffected != 1) {
    throw new Error("Failed to insert person.");
  }

  const id = response.value;
  return <Person>{ id, name, age };
}

export function updatePersonHome(
  id: i32,
  longitude: f64,
  latitude: f64,
): Person | null {
  const query = `update people set home = point($1, $2) where id = $3`;

  const params = new postgresql.Params();
  params.push(longitude);
  params.push(latitude);
  params.push(id);

  const response = postgresql.execute(host, query, params);

  if (response.rowsAffected != 1) {
    console.error(
      `Failed to update person with id ${id}. The record may not exist.`,
    );
    return null;
  }

  return getPerson(id);
}

export function deletePerson(id: i32): string {
  const query = "delete from people where id = $1";

  const params = new postgresql.Params();
  params.push(id);

  const response = postgresql.execute(host, query, params);

  if (response.rowsAffected != 1) {
    console.error(
      `Failed to delete person with id ${id}. The record may not exist.`,
    );
    return "failure";
  }

  return "success";
}
