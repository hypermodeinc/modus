/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// NOTE: This file is special.
// Any function exported here will become a WASM export in the final build.
// It is hooked via afterParse in the transform.

export {
  activateAgent as _modus_agent_activate,
  shutdownAgent as _modus_agent_shutdown,
  handleMessage as _modus_agent_handle_message,
  getAgentState as _modus_agent_get_state,
  setAgentState as _modus_agent_set_state,
} from "./agent";
