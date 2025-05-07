/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { Agent } from "@hypermode/modus-sdk-as";

/**
 * This is a very simple agent that is used to demonstrate how Modus Agents work.
 * It keeps a simple counter that can be incremented and queried.
 * A more complex agent would have more state and more complex logic, including
 * interacting with AI models, databases, services, and possibly other agents.
 */
export class CounterAgent extends Agent {
  // Agents are identified by a name.  Each agent in your project must have a unique name.
  // The name is used to register the agent with the host, and to send messages to it.
  // It should be a short, descriptive name that reflects the purpose of the agent.
  get name(): string {
    return "Counter";
  }

  // Agents can keep state in the form of class properties.
  // This state is only visible to the active instance of the agent.
  // In this case, we are just using a simple integer field to hold the count.
  private count: i32 = 0;

  // The agent should be able to save its state and restore it later.
  // This is used for persisting data across soft restarts of the agent,
  // such as when updating the agent code, or when the agent is suspended and resumed.
  // The getState and setState methods below are used for this purpose.

  // This method should return the current state of the agent as a string.
  // Any format is fine, but it should be consistent and easy to parse.
  getState(): string | null {
    return this.count.toString();
  }

  // This method should set the state of the agent from a string.
  // The string should be in the same format as the one returned by getState.
  // Be sure to consider data compatibility when changing the format of the state.
  setState(data: string | null): void {
    if (data == null) {
      return;
    }
    this.count = i32.parse(data);
  }

  // When the agent is started, this method is automatically called.
  // It is optional, but can be used to initialize state, retrieve data, etc.
  // This is a good place to set up any listeners or subscriptions.
  onStart(): void {
    console.info("Counter agent started");
  }

  // When the agent is stopped, this method is automatically called.
  // It is optional, but can be used to clean up any resources, send final data, etc.
  // This is a good place to unsubscribe from any listeners or subscriptions.
  onStop(): void {
    console.info("Counter agent stopped");
  }

  // If the agent is reloaded, this method is automatically called.
  // It is optional, but can be used to keep track of the agent's status.
  onReload(): void {
    console.info("Counter agent reloaded");
  }

  // This method is called when the agent receives a message.
  // This is how agents update their state and share data.
  onReceiveMessage(name: string, data: string | null): string | null {
    // A "count" message just returns the current count.
    if (name == "count") {
      return this.count.toString();
    }

    // An "increment" message increments the count and returns the new count.
    // If the message has data, it is used as the increment value. Otherwise it defaults to 1.
    if (name == "increment") {
      if (data != null) {
        this.count += i32.parse(data);
      } else {
        this.count++;
      }
      return this.count.toString();
    }

    return null;
  }
}
