/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// This file is special.
// Any function exported here will become a WASM export in the final build.

export {
  registerAgents as _modus_register_agents,
  startAgent as _modus_start_agent,
  stopAgent as _modus_stop_agent,
  getAgentState as _modus_get_agent_state,
  setAgentState as _modus_set_agent_state,
  handleMessage as _modus_agent_handle_message,
} from "./agent";
