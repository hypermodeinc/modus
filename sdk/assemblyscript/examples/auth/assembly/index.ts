/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { auth } from "@hypermode/modus-sdk-as";


@json
export class Claims {
  public exp!: i64;
  public iat!: i64;
  public iss!: string;
  public jti!: string;
  public nbf!: i64;
  public sub!: string;


  @alias("user-id")
  public userId!: string;
}
export function getJWTClaims(): Claims {
  return auth.jwt.getClaims<Claims>();
}
