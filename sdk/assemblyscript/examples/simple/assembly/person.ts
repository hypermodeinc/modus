// This class is used in the examples.
// It is separated just to demonstrated how to use multiple files.
// Note, this is just one way to define a class in AssemblyScript.
export class Person {
  firstName: string;
  lastName: string;
  fullName: string;

  constructor(firstName: string, lastName: string) {
    this.firstName = firstName;
    this.lastName = lastName;
    this.fullName = `${firstName} ${lastName}`;
  }
}
