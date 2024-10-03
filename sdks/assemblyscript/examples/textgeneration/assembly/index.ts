import { JSON } from "json-as";
import { models } from "@hypermode/functions-as";
import { Product, sampleProductJson } from "./product";

import {
  OpenAIChatModel,
  ResponseFormat,
  SystemMessage,
  UserMessage,
} from "@hypermode/models-as/models/openai/chat";

// In this example, we will generate text using the OpenAI Chat model.
// See https://platform.openai.com/docs/api-reference/chat/create for more details
// about the options available on the model, which you can set on the input object.

// This model name should match the one defined in the hypermode.json manifest file.
const modelName: string = "text-generator";

// This function generates some text based on the instruction and prompt provided.
export function generateText(instruction: string, prompt: string): string {
  // The imported ChatModel interface follows the OpenAI Chat completion model input format.
  const model = models.getModel<OpenAIChatModel>(modelName);
  const input = model.createInput([
    new SystemMessage(instruction),
    new UserMessage(prompt),
  ]);

  // This is one of many optional parameters available for the OpenAI Chat model.
  input.temperature = 0.7;

  // Here we invoke the model with the input we created.
  const output = model.invoke(input);

  // The output is also specific to the ChatModel interface.
  // Here we return the trimmed content of the first choice.
  return output.choices[0].message.content.trim();
}

// This function generates a single product.
export function generateProduct(category: string): Product {
  // We can get creative with the instruction and prompt to guide the model
  // in generating the desired output.  Here we provide a sample JSON of the
  // object we want the model to generate.
  const instruction = `Generate a product for the category provided.
Only respond with valid JSON object in this format:
${sampleProductJson}`;
  const prompt = `The category is "${category}".`;

  // Set up the input for the model, creating messages for the instruction and prompt.
  const model = models.getModel<OpenAIChatModel>(modelName);
  const input = model.createInput([
    new SystemMessage(instruction),
    new UserMessage(prompt),
  ]);

  // Let's increase the temperature to get more creative responses.
  // Be careful though, if the temperature is too high, the model may generate invalid JSON.
  input.temperature = 1.2;

  // This model also has a response format parameter that can be set to JSON,
  // Which, along with the instruction, can help guide the model in generating valid JSON output.
  input.responseFormat = ResponseFormat.Json;

  // Here we invoke the model with the input we created.
  const output = model.invoke(input);

  // The output should contain the JSON string we asked for.
  const json = output.choices[0].message.content.trim();
  const product = JSON.parse<Product>(json);
  return product;
}

// This function generates multiple products.
export function generateProducts(category: string, quantity: i32): Product[] {
  // const instruction = `Generate ${quantity} products for the category provided.`;

  // Similar to the previous example above, we can tailor the instruction and prompt
  // to guide the model in generating the desired output.  Note that understanding the behavior
  // of the model is important to get the desired results.  In this case, we need the model
  // to return an _object_ containing an array, not an array of objects directly.  That's because
  // the model will not reliably generate an array of objects directly.
  const instruction = `Generate ${quantity} products for the category provided.
Only respond with valid JSON object containing a valid JSON array named 'list', in this format:
{"list":[${sampleProductJson}]}`;
  const prompt = `The category is "${category}".`;

  // Set up the input for the model, creating messages for the instruction and prompt.
  const model = models.getModel<OpenAIChatModel>(modelName);
  const input = model.createInput([
    new SystemMessage(instruction),
    new UserMessage(prompt),
  ]);

  // Adjust the model inputs, just like in the previous example.
  // Be careful, if the temperature is too high, the model may generate invalid JSON.
  input.temperature = 1.2;
  input.responseFormat = ResponseFormat.Json;

  // Here we invoke the model with the input we created.
  const output = model.invoke(input);

  // The output should contain the JSON string we asked for.
  const json = output.choices[0].message.content.trim();

  // We can parse that JSON to a compatible object, to get the data we're looking for.
  const results = JSON.parse<Map<string, Product[]>>(json);
  return results.get("list");
}
