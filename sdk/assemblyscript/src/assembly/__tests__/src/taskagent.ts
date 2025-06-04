import { Agent } from "../../";
import { JSON } from "json-as";

export class TaskManagerAgent extends Agent {
  get name(): string {
    return "TaskManager";
  }

  private tasks: Map<string, bool> = new Map();

  getState(): string | null {
    return JSON.stringify(this.tasks);
  }

  setState(data: string | null): void {
    if (!data) return;

    this.tasks = JSON.parse<Map<string, bool>>(data);
  }

  onInitialize(): void {
    console.info("TaskManager agent started");
  }

  onSuspend(): void {
    console.info("TaskManager agent suspended");
  }

  onResume(): void {
    console.info("TaskManager agent resumed");
  }

  onTerminate(): void {
    console.info("TaskManager agent terminated");
  }

  onReceiveMessage(name: string, data: string | null): string | null {
    if (name == "add") {
      if (!data) return "No task provided";
      this.tasks.set(data, false);
      return `Task added: "${data}"`;
    }

    if (name == "complete") {
      if (!data || !this.tasks.has(data)) return "Task not found";
      this.tasks.set(data, true);
      return `Task completed: "${data}"`;
    }

    if (name == "list") {
      const keys = this.tasks.keys();
      let output = "";
      for (let i = 0, len = keys.length; i < len; i++) {
        const key = keys[i];
        const value = this.tasks.get(key);
        output += `${value ? "[x]" : "[ ]"} ${key}\n`;
      }
      return output.trim();
    }

    if (name == "stats") {
      const keys = this.tasks.keys();
      let total = 0;
      let done = 0;
      for (let i = 0, len = keys.length; i < len; i++) {
        const key = keys[i];
        const value = this.tasks.get(key);
        if (value !== null) {
          total++;
          if (value) done++;
        }
      }
      return `Tasks: ${done}/${total} completed`;
    }

    return "Unknown command";
  }
}
