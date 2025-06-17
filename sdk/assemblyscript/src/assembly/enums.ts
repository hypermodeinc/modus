/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// eslint-disable-next-line @typescript-eslint/no-namespace
export namespace Duration {
  export const zero = 0;
  export const nanosecond = 1;
  export const microsecond = 1000 * Duration.nanosecond;
  export const millisecond = 1000 * Duration.microsecond;
  export const second = 1000 * Duration.millisecond;
  export const minute = 60 * Duration.second;
  export const hour = 60 * Duration.minute;
  export const maxValue = i64.MAX_VALUE;
}
export type Duration = i64;

// eslint-disable-next-line @typescript-eslint/no-namespace
export namespace AgentStatus {
  export const Starting = "starting";
  export const Running = "running";
  export const Suspending = "suspending";
  export const Suspended = "suspended";
  export const Resuming = "resuming";
  export const Stopping = "stopping";
  export const Terminated = "terminated";
}
export type AgentStatus = string;
