import { Agent } from "@hypermode/modus-sdk-as";

export class CounterAgent extends Agent {
  // Agents are identified by a name.  Each agent in your project must have a unique name.
  // The name is used to register the agent with the host, and to send messages to it.
  // It should be a short, descriptive name that reflects the purpose of the agent.
  get name(): string {
    return "Counter";
  }

  // Agents can keep state in the form of class properties.
  // This state is only visible to the active instance of the agent.
  // This example agent keeps a simple counter.
  private count: i32 = 0;

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

  // This method is called when the agent receives a message.
  // This is how agents update their state and share data.
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  onReceiveMessage(name: string, data: string | null): string | null {
    if (data == null) {
      console.debug(`Counter agent ${this.id} received message: ${name}`);
    } else {
      console.debug(
        `Counter agent ${this.id} received message: ${name} with data: ${data}`,
      );
    }
    // A "count" message just returns the current count.
    if (name == "count") {
      return this.count.toString();
    }

    // An "increment" message increments the count and returns the new count.
    if (name == "increment") {
      this.count++;
      return this.count.toString();
    }

    return null;
  }

  // The agent should be able to save its state and restore it later.
  // This is used for persisting data across soft restarts of the agent,
  // such as when updating the agent code, or when the agent is suspended and resumed.

  // This method should return the current state of the agent as a string.
  // Any format is fine, but it should be consistent and easy to parse.
  getState(): string | null {
    return this.count.toString();
  }

  // This method should set the state of the agent from a string.
  // The string should be in the same format as the one returned by getState().
  // Be sure to consider data compatibility when changing the format of the state.
  setState(data: string | null): void {
    if (data == null) {
      return;
    }
    this.count = i32.parse(data);
  }
}
