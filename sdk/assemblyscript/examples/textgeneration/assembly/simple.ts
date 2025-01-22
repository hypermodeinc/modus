/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { models } from "@hypermode/modus-sdk-as";

import {
  OpenAIChatModel,
  SystemMessage,
  UserMessage,
} from "@hypermode/modus-sdk-as/models/openai/chat";

// In this example, we will generate text and other content using OpenAI chat models.
// See https://platform.openai.com/docs/api-reference/chat/create for more details
// about the options available on the model, which you can set on the input object.

/**
 * This function generates some text based on the instruction and prompt provided.
 */
export function generateText(instruction: string, prompt: string): string {
  // The imported OpenAIChatModel interface follows the OpenAI chat completion model input format.
  const model = models.getModel<OpenAIChatModel>("text-generator");

  // We'll start by creating an input object using the instruction and prompt provided.
  const input = model.createInput([
    new SystemMessage(instruction),
    new UserMessage(prompt),
    // ... if we wanted to add more messages, we could do so here.
  ]);

  // This is one of many optional parameters available for the OpenAI chat model.
  input.temperature = 0.7;

  // Here we invoke the model with the input we created.
  const output = model.invoke(input);

  // The output is also specific to the OpenAIChatModel interface.
  // Here we return the trimmed content of the first choice.
  return output.choices[0].message.content.trim();
}
