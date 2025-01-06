/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { models, localtime } from "@hypermode/modus-sdk-as";
import { JSON } from "json-as";

import {
  OpenAIChatModel,
  SystemMessage,
  Tool,
  ToolMessage,
  UserMessage,
} from "@hypermode/modus-sdk-as/models/openai/chat";

/**
 * This function generates some text from the model, which has already been given some tools to call,
 * and instructions on how to use them.  The tools allow the model to answer time-related questions,
 * such as "What time is it?", or "What time is it in New York?".
 * @param prompt The prompt to provide to the model.
 */
export function generateTextWithTools(prompt: string): string {
  const model = models.getModel<OpenAIChatModel>("text-generator");

  const instruction = `
	You are a helpful assistant that understands time in various parts of the world.
	Answer the user's question as directly as possible. If you need more information, ask for it.
	When stating a time zone to the user, either use a descriptive name such as "Pacific Time" or use the location's name.
	If asked for a time in a location, and the given location has multiple time zones, return the time in each of them.
	Make your answers as brief as possible.
    `;

  const input = model.createInput([
    new SystemMessage(instruction),
    new UserMessage(prompt),
  ]);

  // Let's turn the temperature down a bit, to make the model's responses more predictable.
  input.temperature = 0.2;

  // This is how you make tools available to the model.
  // Each tool is specified by the name of the function to call, and a description of what it does.
  //
  // If the function requires parameters, you can specify simple parameters individually with the withParameter method,
  // or you can pass a JSON Schema string directly to the withParametersSchema method.
  //
  // NOTE: A future release of Modus may further simplify this process.
  input.tools = [
    Tool.forFunction(
      "getUserTimeZone",
      "Returns the user's current IANA time zone.",
    ),
    Tool.forFunction(
      "getCurrentTimeInUserTimeZone",
      "Gets the current date and time in the user's time zone.  Returns the user's IANA time zone, and the time as an ISO 8601 string.",
    ),
    Tool.forFunction(
      "getCurrentTime",
      "Returns the current date and time in the specified IANA time zone as an ISO 8601 string.",
    ).withParameter(
      "tz",
      "string",
      `An IANA time zone identifier such as "America/New_York".`,
    ),
  ];

  // To use the tools, you'll need to have a "conversation loop" like this one.
  // The model may be called multiple times, depending on the complexity of the conversation and the tools used.
  while (true) {
    const output = model.invoke(input);
    const msg = output.choices[0].message;

    // If the model requested tools to be called, we'll need to do that before continuing.
    if (msg.toolCalls.length > 0) {
      // First, add the model's response to the conversation as an "Assistant Message".
      input.messages.push(msg.toAssistantMessage());

      // Then call the tools the model requested.  Add each result to the conversation as a "Tool Message".
      for (let i = 0; i < msg.toolCalls.length; i++) {
        const tc = msg.toolCalls[i];

        // The model will tell us which tool to call, and what arguments to pass to it.
        // We'll need to parse the function name and arguments in order to call the appropriate function.
        //
        // If the function errors (for example, due to invalid input), we can either continue
        // the conversation by providing the model with the error message (as shown here),
        // or we can end the conversation by returning error directly or throwing an exception.
        //
        // NOTE: A future release of Modus may simplify this process.
        let toolMsg: ToolMessage<string>;
        const fnName = tc.function.name;
        if (fnName === "getCurrentTime") {
          const args = JSON.parse<Map<string, string>>(tc.function.arguments);
          const result = getCurrentTime(args.get("tz"));
          toolMsg = new ToolMessage(result, tc.id);
        } else if (fnName === "getUserTimeZone") {
          const timeZone = getUserTimeZone();
          toolMsg = new ToolMessage(timeZone, tc.id);
        } else if (fnName === "getCurrentTimeInUserTimeZone") {
          const result = getCurrentTimeInUserTimeZone();
          toolMsg = new ToolMessage(result, tc.id);
        } else {
          throw new Error(`Unknown tool call: ${tc.function.name}`);
        }

        // Add the tool's response to the conversation.
        input.messages.push(toolMsg);
      }
    } else if (msg.content != "") {
      // return the model's final response to the user.
      return msg.content.trim();
    } else {
      // If the model didn't ask for tools, and didn't return any content, something went wrong.
      throw new Error("Invalid response from model.");
    }
  }
}

// The following functions are made available as "tools"
// for the model to call in the example above.

/**
 * This function will return the current time in a given time zone.
 */
function getCurrentTime(tz: string): string {
  if (!localtime.IsValidTimeZone(tz)) {
    return "error: invalid time zone";
  }
  return localtime.NowInZone(tz);
}

/**
 * This function will return the default time zone for the user.
 */
function getUserTimeZone(): string {
  return localtime.GetTimeZone();
}

/**
 * This function combines the functionality of getCurrentTime and getUserTimeZone,
 * to reduce the number of tool calls required by the model.
 */
function getCurrentTimeInUserTimeZone(): string {
  const tz = localtime.GetTimeZone();
  if (tz === "") {
    return "error: cannot determine user's time zone";
  }

  const time = localtime.NowInZone(tz);
  const result = <TimeAndZone>{ time, zone: tz };
  return JSON.stringify(result);
}


@json
class TimeAndZone {
  time!: string;
  zone!: string;
}
