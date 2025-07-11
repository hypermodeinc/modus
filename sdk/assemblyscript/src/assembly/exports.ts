/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// NOTE: This file is special.
// Any function exported here will become a WASM export in the final build.
// It is hooked via afterParse in the transform.

export {
  activateAgent as _modus_agent_activate,
  handleEvent as _modus_agent_handle_event,
  handleMessage as _modus_agent_handle_message,
  getAgentState as _modus_agent_get_state,
  setAgentState as _modus_agent_set_state,
} from "./agent";
