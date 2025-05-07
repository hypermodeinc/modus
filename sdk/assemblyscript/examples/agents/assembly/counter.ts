import { Agent } from "@hypermode/modus-sdk-as";

export class CounterAgent extends Agent {
  // Agents have a name, used for debugging and logging.
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
  onReceiveMessage(name: string, data: string | null): string | null {
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

  // This part is work in progress.
  // The agent should be able to save its state and restore it later.
  // TODO: use this in the runtime and see if we can implement this automatically in user code.

  getState(): string | null {
    return this.count.toString();
  }

  setState(data: string | null): void {
    if (data == null) {
      return;
    }
    this.count = i32.parse(data);
  }
}
