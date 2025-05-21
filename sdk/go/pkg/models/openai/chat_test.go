/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package openai_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

func TestMessagesRoundTrip(t *testing.T) {
	msgs := []openai.RequestMessage{
		openai.NewSystemMessage("You are a helpful assistant."),
		openai.NewUserMessage("What is the capital of France?"),
		openai.NewAssistantMessage("The capital of France is Paris."),
	}

	data, err := json.Marshal(msgs)
	if err != nil {
		t.Fatalf("Failed to marshal messages: %v", err)
	}

	parsedMsgs, err := openai.ParseMessages(data)
	if err != nil {
		t.Fatalf("Failed to parse messages: %v", err)
	}
	if len(parsedMsgs) != len(msgs) {
		t.Fatalf("Expected %d messages, but got %d", len(msgs), len(parsedMsgs))
	}

	for i, msg := range msgs {
		if msg.Role() != parsedMsgs[i].Role() {
			t.Errorf("Expected role %s for message %d, but got %s", msg.Role(), i, parsedMsgs[i].Role())
		}
	}

	roundTrip, err := json.Marshal(parsedMsgs)
	if err != nil {
		t.Fatalf("Failed to marshal parsed messages: %v", err)
	}
	if !bytes.Equal(data, roundTrip) {
		t.Fatalf("Expected original and parsed messages to have the same JSON representation, but they differ")
	}
}
