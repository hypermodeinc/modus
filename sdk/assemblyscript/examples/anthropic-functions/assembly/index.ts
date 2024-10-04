import { http, models } from "@hypermode/modus-sdk-as";
import { JSON } from "json-as";

import {
  AnthropicMessagesModel,
  Message,
  UserMessage,
  ToolChoiceTool,
} from "@hypermode/models-as/models/anthropic/messages";


@json
class StockPriceInput {
  symbol!: string;
}

// This model name should match the one defined in the hypermode.json manifest file.
const modelName: string = "text-generator";

export function getStockPrice(company: string, useTools: bool): string {
  const model = models.getModel<AnthropicMessagesModel>(modelName);
  const messages: Message[] = [
    new UserMessage(`what is the stock price of ${company}?`),
  ];
  const input = model.createInput(messages);
  // For Anthropic, system is passed as parameter to the invoke, not as a message
  input.system =
    "You are a helpful assistant. Do not answer if you do not have up-to-date information.";

  // Optional parameters
  input.temperature = 1;
  input.maxTokens = 100;

  if (useTools) {
    input.tools = [
      {
        name: "stock_price",
        inputSchema: `{"type":"object","properties":{"symbol":{"type":"string","description":"The stock symbol"}},"required":["symbol"]}`,
        description: "gets the stock price of a symbol",
      },
    ];
    input.toolChoice = ToolChoiceTool("stock_price");
  }

  // Here we invoke the model with the input we created.
  const output = model.invoke(input);

  if (output.content.length !== 1) {
    throw new Error("Unexpected output content length");
  }
  // If tools are not used, the output will be a text block
  if (output.content[0].type === "text") {
    return output.content[0].text!.trim();
  }
  if (output.content[0].type !== "tool_use") {
    throw new Error(`Unexpected content type: ${output.content[0].type}`);
  }

  const toolUse = output.content[0];
  const inputs = toolUse.input!;

  const parsedInput = JSON.parse<StockPriceInput>(inputs);
  const symbol = parsedInput.symbol;
  const stockPrice = callStockPriceApi(symbol);
  return `The stock price of ${symbol} is $${stockPrice.GlobalQuote.price}`;
}


@json
class StockPriceAPIResponse {

  @alias("Global Quote")
  GlobalQuote!: GlobalQuote;
}


@json
class GlobalQuote {

  @alias("01. symbol")
  symbol!: string;


  @alias("05. price")
  price!: string;
}

function callStockPriceApi(symbol: string): StockPriceAPIResponse {
  const req = new http.Request(
    `https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=${symbol}`,
  );
  const resp = http.fetch(req);
  if (!resp.ok) {
    throw new Error(
      `Failed to fetch stock price. Received: ${resp.status} ${resp.statusText}`,
    );
  }

  return resp.json<StockPriceAPIResponse>();
}
