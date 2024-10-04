/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

// These classes are used by the example functions in the index.ts file.

@json
export class Person {
  id: string | null = null;
  firstName: string = "";
  lastName: string = "";
}


@json
export class PeopleData {
  people!: Person[];
}


@json
export class AddPersonPayload {
  addPerson!: PeopleData;
}


@json
export class AggregatePersonResult {
  aggregatePerson!: GQLAggregateValues;
}


@json
class GQLAggregateValues {
  count: u32 = 0;
}
