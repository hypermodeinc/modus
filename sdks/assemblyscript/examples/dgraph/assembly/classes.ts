// These classes are used by the example functions in the index.ts file.

@json
export class Person {
  constructor(
    uid: string = "",
    firstName: string = "",
    lastName: string = "",
    dType: string[] = [],
  ) {
    this.uid = uid;
    this.firstName = firstName;
    this.lastName = lastName;
    this.dType = dType;
  }
  uid: string = "";
  firstName: string = "";
  lastName: string = "";


  @alias("dgraph.type")
  dType: string[] = [];
}


@json
export class PeopleData {
  people!: Person[];
}
