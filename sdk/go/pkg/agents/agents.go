/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package agents

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/hypermodeinc/modus/sdk/go/pkg/console"
	"github.com/hypermodeinc/modus/sdk/go/pkg/utils"
)

// AgentInfo contains information about an agent.
type AgentInfo struct {
	Id     string
	Name   string
	Status string
}

// AgentEvent is an interface that represents an event emitted by an agent.
// Custom agent events should implement this interface.
type AgentEvent interface {
	// EventName returns the type of the event.
	EventName() string
}

type AgentStatus string

const (
	AgentStatusStarting   AgentStatus = "starting"
	AgentStatusRunning    AgentStatus = "running"
	AgentStatusSuspending AgentStatus = "suspending"
	AgentStatusSuspended  AgentStatus = "suspended"
	AgentStatusResuming   AgentStatus = "resuming"
	AgentStatusStopping   AgentStatus = "stopping"
	AgentStatusTerminated AgentStatus = "terminated"
)

type agentEventAction string

const (
	agentEventActionInitialize agentEventAction = "initialize"
	agentEventActionSuspend    agentEventAction = "suspend"
	agentEventActionResume     agentEventAction = "resume"
	agentEventActionTerminate  agentEventAction = "terminate"
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

	info := hostStartAgent(&name)
	return *info, nil
}

// Stops an agent with the given ID and returns its status info.
// This will terminate the agent, and it cannot be resumed or restarted.
func Stop(agentId string) (AgentInfo, error) {
	if len(agentId) == 0 {
		return AgentInfo{}, errors.New("invalid agent ID")
	}

	if info := hostGetAgentInfo(&agentId); info == nil {
		return AgentInfo{}, fmt.Errorf("agent %s not found", agentId)
	}

	if info := hostStopAgent(&agentId); info != nil {
		return *info, nil
	}

	return AgentInfo{}, fmt.Errorf("failed to stop agent %s", agentId)
}

// Gets information about an agent with the given ID.
func GetInfo(agentId string) (AgentInfo, error) {
	if len(agentId) == 0 {
		return AgentInfo{}, errors.New("invalid agent ID")
	}

	info := hostGetAgentInfo(&agentId)
	if info == nil {
		return AgentInfo{}, fmt.Errorf("agent %s not found", agentId)
	}

	return *info, nil
}

// Returns a list of all agents, except those that have been fully terminated.
func ListAll() ([]AgentInfo, error) {
	agents := hostListAgents()
	if agents == nil {
		return nil, errors.New("failed to list agents")
	}
	return *agents, nil
}

// These functions are only invoked as wasm exports from the host.
// Assigning them to discard variables avoids "unused" linting errors.
var (
	_ = activateAgent
	_ = handleEvent
	_ = handleMessage
	_ = getAgentState
	_ = setAgentState
)

// The Modus Runtime will call this function to set the active agent.
// This is called after the wasm module instance is loaded, regardless of whether
// the agent is newly started or resumed from a suspended state.
//
//go:export _modus_agent_activate
func activateAgent(name, id string) {
	if activeAgent != nil {
		console.Error("Another agent is already active in this module instance.")
		return
	}

	if agent, ok := agents[name]; !ok {
		console.Errorf("Agent %s not found.", name)
	} else {
		activeAgent = &agent
		activeAgentId = &id
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

// The Modus Runtime will call this function when an agent event occurs.
// Events are part of the agent's lifecycle, managed by the Modus Runtime.
//
//go:export _modus_agent_handle_event
func handleEvent(action string) {
	if activeAgent == nil {
		console.Errorf("No active agent to process %s event action.", action)
		return
	}
	agent := *activeAgent

	switch agentEventAction(action) {
	case agentEventActionInitialize:
		if err := (agent).OnInitialize(); err != nil {
			console.Errorf("Error initializing agent %s: %v", agent.Name(), err)
		}
	case agentEventActionSuspend:
		if err := agent.OnSuspend(); err != nil {
			console.Errorf("Error suspending agent %s: %v", agent.Name(), err)
		}
	case agentEventActionResume:
		if err := agent.OnResume(); err != nil {
			console.Errorf("Error resuming agent %s: %v", agent.Name(), err)
		}
	case agentEventActionTerminate:
		if err := agent.OnTerminate(); err != nil {
			console.Errorf("Error terminating agent %s: %v", agent.Name(), err)
		}
	default:
		console.Errorf("Unknown agent event action: %s", action)
	}
}

// The Modus Runtime will call this function when an agent receives a message.
// Messages are user-defined and can be sent from API functions or other agents.
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

	// OnInitialize is called when the agent is started.
	// Custom agents may implement this method to perform any initialization.
	OnInitialize() error

	// OnSuspend is called when the agent is suspended.
	// Custom agents may implement this method if for example, to send a notification.
	// Note that you do not need to save the internal state of the agent here, as that is handled automatically.
	OnSuspend() error

	// OnResume is called when the agent is resumed from a suspended state.
	// Custom agents may implement this method if for example, to send a notification.
	// Note that you do not need to restore the internal state of the agent here, as that is handled automatically.
	OnResume() error

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

func (a *AgentBase) OnInitialize() error {
	return nil
}

func (a *AgentBase) OnSuspend() error {
	return nil
}

func (a *AgentBase) OnResume() error {
	return nil
}

func (a *AgentBase) OnTerminate() error {
	return nil
}

func (a *AgentBase) OnReceiveMessage(msgName string, data *string) (*string, error) {
	return nil, nil
}

// Publishes an event from this agent to any subscribers.
func (a *AgentBase) PublishEvent(event AgentEvent) error {
	bytes, err := utils.JsonSerialize(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event data: %w", err)
	}
	data := string(bytes)
	name := event.EventName()
	hostPublishEvent(activeAgentId, &name, &data)
	return nil
}

type MessageResponse struct {
	data  *string
	error *string
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
// Data can be sent with the message using the WithData option.
// The message is sent synchronously and the function will block until a response is received or the timeout is reached.
// The timeout can be set using the WithTimeout option. If there is no timeout set, it will wait indefinitely until a
// response is received, or until calling function is cancelled.
func SendMessage(agentId, msgName string, options ...messageOption) (*string, error) {
	m := &message{
		name:    msgName,
		timeout: math.MaxInt64, // default to no timeout (wait indefinitely)
	}

	for _, opt := range options {
		opt(m)
	}

	return sendMessage(agentId, *m)
}

// SendMessageAsync sends a message to the specified agent without waiting for a response.
// Data can be sent with the message using the WithData option.
// The message is sent asynchronously and the function will return immediately.
// The WithTimeout option is ignored for asynchronous messages, as they do not wait for a response.
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
	if response.error != nil {
		return nil, errors.New(*response.error)
	}

	return response.data, nil
}
