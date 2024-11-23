/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

// This class is used in the examples.
// It is separated just to demonstrated how to use multiple files.
// Note, this is just one way to define a class in AssemblyScript.

/**
 * A simple object representing a person.
 */
export class Person {
  /**
   * The first name of the person.
   */
  firstName: string;

  /**
   * The last name of the person.
   */
  lastName: string;

  /**
   * The full name of the person.
   */
  fullName: string;

  constructor(firstName: string, lastName: string) {
    this.firstName = firstName;
    this.lastName = lastName;
    this.fullName = `${firstName} ${lastName}`;
  }
}
