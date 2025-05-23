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

  // When the agent is started, this method is automatically called. Implementing it is optional.
  // If you don't need to do anything special when the agent starts, then you can omit it.
  // It can be used to initialize state, retrieve data, etc.
  // This is a good place to set up any listeners or subscriptions.
  onInitialize(): void {
    console.info("Counter agent started");
  }

  // When the agent is suspended, this method is automatically called.  Implementing it is optional.
  // If you don't need to do anything special when the agent is suspended, then you can omit it.
  // The agent may be suspended for a variety of reasons, including:
  // - The agent code has being updated.
  // - The host is shutting down or restarting.
  // - The agent is being suspended to save resources.
  // - The agent is being relocated to a different host.
  // Note that the agent may be suspended and resumed multiple times during its lifetime,
  // but the Modus Runtime will automatically save and restore the state of the agent,
  // so you don't need to worry about that here.
  onSuspend(): void {
    console.info("Counter agent suspended");
  }

  // When the agent is resumed, this method is automatically called.  Implementing it is optional.
  // If you don't need to do anything special when the agent is resumed, then you can omit it.
  onResume(): void {
    console.info("Counter agent resumed");
  }

  // When the agent is terminated, this method is automatically called.  Implementing it is optional.
  // It can be used to send final data somewhere, such as a database or an API.
  // This is a good place to unsubscribe from any listeners or subscriptions.
  // Note that resources are automatically cleaned up when the agent is terminated,
  // so you don't need to worry about that here.
  // Once an agent is terminated, it cannot be resumed.
  onTerminate(): void {
    console.info("Counter agent terminated");
  }

  // This method is called when the agent receives a message.
  // This is how agents update their state and share data.
  // You should implement this method to handle messages sent to this agent.
  onReceiveMessage(name: string, data: string | null): string | null {
    // You can use the name of the message to determine what to do.
    // You can either handle the message here, or pass it to another method or function.
    // If you don't have any response to send back, return null.

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
