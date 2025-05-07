/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"fmt"
	"strconv"

	"github.com/hypermodeinc/modus/sdk/go/pkg/agents"
)

/*
 * This is a very simple agent that is used to demonstrate how Modus Agents work.
 * It keeps a simple counter that can be incremented and queried.
 * A more complex agent would have more state and more complex logic, including
 * interacting with AI models, databases, services, and possibly other agents.
 */
type CounterAgent struct {

	// All agents must embed the AgentBase struct.
	agents.AgentBase

	// Additional fields can be added to the agent to hold state.
	// This is state is only visible to the active instance of the agent.
	// In this case, we are just using a simple integer field to hold the count.
	count int
}

// Agents are identified by a name.  Each agent in your project must have a unique name.
// The name is used to register the agent with the host, and to send messages to it.
// It should be a short, descriptive name that reflects the purpose of the agent.
func (c *CounterAgent) Name() string {
	return "Counter"
}

// The agent should be able to save its state and restore it later.
// This is used for persisting data across soft restarts of the agent,
// such as when updating the agent code, or when the agent is suspended and resumed.
// The GetState and SetState methods below are used for this purpose.

// This method should return the current state of the agent as a string.
// Any format is fine, but it should be consistent and easy to parse.
func (c *CounterAgent) GetState() *string {
	s := strconv.Itoa(c.count)
	return &s
}

// This method should set the state of the agent from a string.
// The string should be in the same format as the one returned by GetState.
// Be sure to consider data compatibility when changing the format of the state.
func (c *CounterAgent) SetState(data *string) {
	if data == nil {
		return
	}
	if n, err := strconv.Atoi(*data); err == nil {
		c.count = n
	}
}

// When the agent is started, this method is automatically called.
// It is optional, but can be used to initialize state, retrieve data, etc.
// This is a good place to set up any listeners or subscriptions.
func (c *CounterAgent) OnStart() error {
	fmt.Println("Counter agent started")
	return nil
}

// When the agent is stopped, this method is automatically called.
// It is optional, but can be used to clean up any resources, send final data, etc.
// This is a good place to unsubscribe from any listeners or subscriptions.
func (c *CounterAgent) OnStop() error {
	fmt.Println("Counter agent stopped")
	return nil
}

// If the agent is reloaded, this method is automatically called.
// It is optional, but can be used to keep track of the agent's status.
func (c *CounterAgent) OnReload() error {
	fmt.Println("Counter agent reloaded")
	return nil
}

// This method is called when the agent receives a message.
// This is how agents update their state and share data.
func (c *CounterAgent) OnReceiveMessage(msgName string, data *string) (*string, error) {
	switch msgName {
	case "count":
		s := strconv.Itoa(c.count)
		return &s, nil
	case "increment":
		if data == nil {
			c.count++
		} else if n, err := strconv.Atoi(*data); err == nil {
			c.count += n
		}
		s := strconv.Itoa(c.count)
		return &s, nil
	}

	return nil, nil
}
