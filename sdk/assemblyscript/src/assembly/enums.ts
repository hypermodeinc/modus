/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

export enum Duration {
  zero = 0,
  nanosecond = 1,
  microsecond = 1000 * Duration.nanosecond,
  millisecond = 1000 * Duration.microsecond,
  second = 1000 * Duration.millisecond,
  minute = 60 * Duration.second,
  hour = 60 * Duration.minute,
}

// eslint-disable-next-line @typescript-eslint/no-namespace
export namespace AgentStatus {
  export const Starting = "starting";
  export const Running = "running";
  export const Suspending = "suspending";
  export const Suspended = "suspended";
  export const Restoring = "restoring";
  export const Terminating = "terminating";
  export const Terminated = "terminated";
}
export type AgentStatus = string;
