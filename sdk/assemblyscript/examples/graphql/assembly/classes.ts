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
