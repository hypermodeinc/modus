/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/pkg/http"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/anthropic"
)

type stockPriceInput struct {
	Symbol string `json:"symbol"`
}

// This model name should match the one defined in the modus.json manifest file.
const modelName = "text-generator"

func GetStockPrice(company string, useTools bool) (string, error) {
	model, err := models.GetModel[anthropic.MessagesModel](modelName)
	if err != nil {
		return "", err
	}

	input, err := model.CreateInput(
		anthropic.NewUserMessage(
			anthropic.StringContent(fmt.Sprintf("what is the stock price of %s?", company)),
		),
	)
	if err != nil {
		return "", err
	}
	// For Anthropic, system is passed as parameter to the invoke, not as a message
	input.System = "You are a helpful assistant. Do not answer if you do not have up-to-date information."

	// Optional parameters
	input.Temperature = 1
	input.MaxTokens = 100

	if useTools {
		input.Tools = []*anthropic.Tool{
			{
				Name:        "stock_price",
				InputSchema: `{"type":"object","properties":{"symbol":{"type":"string","description":"The stock symbol"}},"required":["symbol"]}`,
				Description: "gets the stock price of a symbol",
			},
		}
		input.ToolChoice = anthropic.ToolChoiceTool("stock_price")
	}

	// Here we invoke the model with the input we created.
	output, err := model.Invoke(input)
	if err != nil {
		return "", err
	}

	if len(output.Content) != 1 {
		return "", fmt.Errorf("unexpected output content length")
	}
	// If tools are not used, the output will be a text block
	if output.Content[0].Type == "text" {
		return strings.TrimSpace(*output.Content[0].Text), nil
	}
	if output.Content[0].Type != "tool_use" {
		return "", fmt.Errorf("unexpected content type: %s", output.Content[0].Type)
	}

	toolUse := output.Content[0]
	inputs := toolUse.Input

	parsedInput := &stockPriceInput{}
	if err := json.Unmarshal([]byte(*inputs), parsedInput); err != nil {
		return "", err
	}
	symbol := parsedInput.Symbol
	stockPrice, err := callStockPriceAPI(symbol)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("The stock price of %s is %s", symbol, stockPrice), nil
}

type stockPriceAPIResponse struct {
	GlobalQuote struct {
		Symbol string `json:"01. symbol"`
		Price  string `json:"05. price"`
	} `json:"Global Quote"`
}

func callStockPriceAPI(symbol string) (string, error) {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s", symbol)
	req := http.NewRequest(url)
	resp, err := http.Fetch(req)
	if err != nil {
		return "", err
	}
	if !resp.Ok() {
		return "", fmt.Errorf("HTTP request failed with status code %d", resp.Status)
	}

	var data stockPriceAPIResponse
	if err := resp.JSON(&data); err != nil {
		return "", err
	}

	return data.GlobalQuote.Price, nil
}
