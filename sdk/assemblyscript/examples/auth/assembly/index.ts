/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { auth } from "@hypermode/modus-sdk-as";

// This is a simple example of a claims class that can be used to parse the JWT claims.
@json
export class ExampleClaims {
  public exp!: i64;
  public iat!: i64;
  public iss!: string;
  public jti!: string;
  public nbf!: i64;
  public sub!: string;

  // This is an example of a custom claim that can be used to parse the user ID.
  @alias("user-id")
  public userId!: string;
}

// This function can be used to get the JWT claims, and parse them into the Claims class.
export function getJWTClaims(): ExampleClaims {
  return auth.getJWTClaims<ExampleClaims>();
}
