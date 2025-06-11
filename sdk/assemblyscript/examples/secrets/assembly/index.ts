/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { secrets } from "@hypermode/modus-sdk-as";

export function getSecretValue(name: string): string {
  return secrets.getSecretValue(name);
}
