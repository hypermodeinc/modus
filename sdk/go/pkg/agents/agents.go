/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package agents

import (
	"errors"
	"fmt"
	"time"

	"github.com/hypermodeinc/modus/sdk/go/pkg/console"
)

type AgentInfo struct {
	Id     string
	Name   string
	Status string
}

type AgentStatus = string

const (
	AgentStatusStarting    AgentStatus = "starting"
	AgentStatusRunning     AgentStatus = "running"
	AgentStatusSuspending  AgentStatus = "suspending"
	AgentStatusSuspended   AgentStatus = "suspended"
	AgentStatusRestoring   AgentStatus = "restoring"
	AgentStatusTerminating AgentStatus = "terminating"
	AgentStatusTerminated  AgentStatus = "terminated"
)

var agents = make(map[string]Agent)
var activeAgent *Agent
var activeAgentId *string

// Registers an agent so it can be used.
// This function should be called from an init function.
func Register(agent Agent) {
	agents[agent.Name()] = agent
}

// Starts an agent with the given name.
// This can be called from any user code, such as function or another agent's methods.
func Start(name string) (AgentInfo, error) {
	if _, ok := agents[name]; !ok {
		return AgentInfo{}, fmt.Errorf("agent %s not found", name)
	}

	info := hostSpawnAgentActor(&name)
	return *info, nil
}

// Terminates an agent with the given ID.
// Once terminated, the agent cannot be restored.
func Terminate(agentId string) error {
	if ok := hostTerminateAgent(&agentId); !ok {
		return fmt.Errorf("failed to terminate agent %s", agentId)
	}
	return nil
}

// These functions are only invoked as wasm exports from the host.
// Assigning them to discard variables avoids "unused" linting errors.
var (
	_ = activateAgent
	_ = shutdownAgent
	_ = getAgentState
	_ = setAgentState
	_ = handleMessage
)

// The Modus Runtime will call this function to activate an agent.
//
//go:export _modus_agent_activate
func activateAgent(name, id string, reloading bool) {
	if activeAgent != nil {
		console.Error("Another agent is already active.")
		return
	}

	if agent, ok := agents[name]; !ok {
		console.Errorf("Agent %s not found.", name)
		return
	} else {
		activeAgent = &agent
		activeAgentId = &id

		if reloading {
			if err := agent.OnRestore(); err != nil {
				console.Errorf("Error reloading agent %s: %v", name, err)
			}
		} else {
			if err := agent.OnStart(); err != nil {
				console.Errorf("Error starting agent %s: %v", name, err)
			}
		}
	}
}

// The Modus Runtime will call this function to shutdown an agent.
//
//go:export _modus_agent_shutdown
func shutdownAgent(suspending bool) {
	if activeAgent == nil {
		console.Error("No active agent to shutdown.")
		return
	}

	if suspending {
		if err := (*activeAgent).OnSuspend(); err != nil {
			console.Errorf("Error suspending agent %s: %v", (*activeAgent).Name(), err)
			return
		}
	} else {
		if err := (*activeAgent).OnTerminate(); err != nil {
			console.Errorf("Error terminating agent %s: %v", (*activeAgent).Name(), err)
			return
		}
	}
}

// The Modus Runtime will call this function to get the state of an agent.
//
//go:export _modus_agent_get_state
func getAgentState() *string {
	if activeAgent == nil {
		console.Error("No active agent to get state for.")
		return nil
	}
	return (*activeAgent).GetState()
}

// The Modus Runtime will call this function to set the state of an agent.
//
//go:export _modus_agent_set_state
func setAgentState(data *string) {
	if activeAgent == nil {
		console.Error("No active agent to set state for.")
		return
	}
	(*activeAgent).SetState(data)
}

// The Modus Runtime will call this function when an agent receives a message.
//
//go:export _modus_agent_handle_message
func handleMessage(msgName string, data *string) *string {
	if activeAgent == nil {
		console.Error("No active agent to handle message for.")
		return nil
	}
	if response, err := (*activeAgent).OnReceiveMessage(msgName, data); err != nil {
		console.Errorf("Error handling message %s: %v", msgName, err)
		return nil
	} else {
		return response
	}
}

// The Agent interface defines the methods that an agent must implement.
// Agents are long-running processes that can receive messages and maintain state.
type Agent interface {

	// Id returns the unique identifier of the agent.
	// This is set by the Modus Runtime when the agent is activated.
	// Custom agents should not implement this method, but rather use the default implementation provided by AgentBase.
	Id() string

	// Name returns the name of the agent.
	// This should be unique across all agents in the module.
	// Custom agents must implement this method.
	Name() string

	// GetState returns the serialized state of the agent.
	// Custom agents must implement this method.
	GetState() *string

	// SetState sets the state of the agent from a serialized string.
	// Custom agents must implement this method.
	SetState(data *string)

	// OnStart is called when the agent is started.
	// Custom agents may implement this method to perform any initialization.
	OnStart() error

	// OnSuspend is called when the agent is suspended.
	// Custom agents may implement this method if for example, to send a notification of the suspension.
	// Note that you do not need to save the internal state of the agent here, as that is handled automatically.
	OnSuspend() error

	// OnRestore is called when the agent is restored from a suspended state.
	// Custom agents may implement this method if for example, to send a notification of the restoration.
	// Note that you do not need to restore the internal state of the agent here, as that is handled automatically.
	OnRestore() error

	// OnTerminate is called when the agent is terminated.
	// Custom agents may implement this method to send or save any final data.
	OnTerminate() error

	// OnReceiveMessage is called when the agent receives a message.
	// Custom agents may implement this method to handle incoming messages.
	OnReceiveMessage(msgName string, data *string) (*string, error)
}

// The AgentBase struct provides default implementations for the Agent interface methods.
// Custom agents can embed this struct to avoid implementing all methods.
type AgentBase struct {
}

func (a *AgentBase) Id() string {
	return *activeAgentId
}

func (a *AgentBase) OnStart() error {
	return nil
}

func (a *AgentBase) OnSuspend() error {
	return nil
}

func (a *AgentBase) OnRestore() error {
	return nil
}

func (a *AgentBase) OnTerminate() error {
	return nil
}

func (a *AgentBase) OnReceiveMessage(msgName string, data *string) (*string, error) {
	return nil, nil
}

type MessageResponse struct {
	data *string
}

type message struct {
	name    string
	data    *string
	timeout time.Duration
}

type messageOption func(*message)

// WithData sets the data for the message.
// This is optional and can be omitted if no data is needed.
func WithData(data string) messageOption {
	return func(m *message) {
		m.data = &data
	}
}

// WithTimeout sets the timeout for the message.
// This is optional and defaults to 10 seconds for synchronous messages, if not set.
// A value of 0 means the message will be sent asynchronously without waiting for a response.
func WithTimeout(timeout time.Duration) messageOption {
	return func(m *message) {
		m.timeout = timeout
	}
}

// SendMessage sends a message to the specified agent and waits for a response.
// The message is sent synchronously and the function will block until a response is received or the timeout is reached.
// The timeout can be set using the WithTimeout option, and defaults to 10 seconds if not set.
// Data can be sent with the message using the WithData option.
func SendMessage(agentId, msgName string, options ...messageOption) (*string, error) {
	m := &message{
		name:    msgName,
		timeout: 10 * time.Second,
	}

	for _, opt := range options {
		opt(m)
	}

	return sendMessage(agentId, *m)
}

// SendMessageAsync sends a message to the specified agent without waiting for a response.
// The message is sent asynchronously and the function will return immediately.
// Data can be sent with the message using the WithData option.
// The timeout is ignored for asynchronous messages and is always set to 0.
func SendMessageAsync(agentId, msgName string, options ...messageOption) error {
	m := &message{
		name: msgName,
	}

	for _, opt := range options {
		opt(m)
	}

	m.timeout = 0 // force zero timeout for async message
	_, err := sendMessage(agentId, *m)
	return err
}

func sendMessage(agentId string, m message) (*string, error) {
	if m.timeout < 0 {
		return nil, fmt.Errorf("timeout must be zero or positive value")
	}
	if agentId == "" {
		return nil, fmt.Errorf("agentId cannot be empty")
	}
	if m.name == "" {
		return nil, fmt.Errorf("message name cannot be empty")
	}

	response := hostSendMessage(&agentId, &m.name, m.data, int64(m.timeout))
	if response == nil {
		return nil, errors.New("failed to send message to agent")
	}

	return response.data, nil
}
