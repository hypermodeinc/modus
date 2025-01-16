/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { JSON } from "json-as";
import { models } from "@hypermode/modus-sdk-as";

import {
  OpenAIChatModel,
  ResponseFormat,
  SystemMessage,
  UserMessage,
} from "@hypermode/modus-sdk-as/models/openai/chat";

// In this example, we'll demonstrate how to generated structured data using OpenAI chat models.
// We'll be asking the model to populate Product objects with data that fits a particular category.
// The data the model generates will need to conform to the schema of the Product object.

@json
class Product {
  id: string | null = null;
  name: string = "";
  price: f64 = 0.0;
  description: string = "";
}

const sampleProductJson = JSON.stringify(<Product>{
  id: "123",
  name: "Shoes",
  price: 50.0,
  description: "Great shoes for walking.",
});

/**
 * This function generates a single product.
 */
export function generateProduct(category: string): Product {
  // We can get creative with the instruction and prompt to guide the model
  // in generating the desired output.  Here we provide a sample JSON of the
  // object we want the model to generate.
  const instruction = `Generate a product for the category provided.
Only respond with valid JSON object in this format:
${sampleProductJson}`;
  const prompt = `The category is "${category}".`;

  // Set up the input for the model, creating messages for the instruction and prompt.
  const model = models.getModel<OpenAIChatModel>("text-generator");
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

  // We can now parse the JSON string as a Product object.
  const product = JSON.parse<Product>(json);

  return product;
}

/**
 * This function generates multiple products.
 */
export function generateProducts(category: string, quantity: i32): Product[] {
  // We can tailor the instruction and prompt to guide the model in generating the desired output.
  // Note that understanding the behavior of the model is important to get the desired results.
  // In this case, we need the model to return an _object_ containing an array, not an array of objects directly.
  // That's because the model will not reliably generate an array of objects directly.
  const instruction = `Generate ${quantity} products for the category provided.
Only respond with valid JSON object containing a valid JSON array named 'list', in this format:
{"list":[${sampleProductJson}]}`;
  const prompt = `The category is "${category}".`;

  // Set up the input for the model, creating messages for the instruction and prompt.
  const model = models.getModel<OpenAIChatModel>("text-generator");
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

  // Now we can extract the list of products from the data.
  const products = results.get("list");

  return products;
}
