//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package models_test

import (
	"reflect"
	"testing"

	"github.com/hypermodeAI/functions-go/pkg/models"
)

// NOTE:
// The type alias isn't strictly required, but it prevents the model struct from getting
// an exported "ModelBase" field, while still promoting the ModelBase methods.
type testModelBase = models.ModelBase[TestModelInput, TestModelOutput]

type TestModel struct {
	testModelBase
}

type TestModelInput struct {
	Prompt string `json:"prompt"`
}

type TestModelOutput struct {
	Response string `json:"response"`
}

func TestGetModel(t *testing.T) {
	modelName := "test"
	model, err := models.GetModel[TestModel](modelName)
	if err != nil {
		t.Fatalf("Expected no error, but received: %s", err)
	} else if model == nil {
		t.Fatal("Expected a model, but received nil")
	} else if model.Debug {
		t.Error("Expected debug to be false, but received true")
	}

	info := model.Info()
	expectedInfo := &models.ModelInfo{
		Name:     modelName,
		FullName: modelName + "-mock-model",
	}
	if !reflect.DeepEqual(expectedInfo, info) {
		t.Errorf("Expected model info: %v, but received: %v", expectedInfo, info)
	}

	values := models.LookupModelCallStack.Pop()
	if values == nil {
		t.Error("Expected a model name, but none was found.")
	} else {
		if !reflect.DeepEqual(values[0], &expectedInfo.Name) {
			t.Errorf("Expected model name: %s, but received: %s", expectedInfo.Name, values[0])
		}
	}
}

func TestInvokeModel(t *testing.T) {
	modelName := "test"
	model, err := models.GetModel[TestModel](modelName)
	if err != nil {
		t.Fatalf("Expected no error, but received: %s", err)
	} else if model == nil {
		t.Fatal("Expected a model, but received nil")
	}

	input := &TestModelInput{Prompt: "Say Hello."}
	output, err := model.Invoke(input)
	if err != nil {
		t.Fatalf("Expected no error, but received: %s", err)
	}

	if output.Response != models.MockResponseText {
		t.Errorf("Expected response: %s, but received: %s", models.MockResponseText, output.Response)
	}

	values := models.InvokeModelCallStack.Pop()
	if values == nil {
		t.Fatal("Expected model name and input, but none were found.")
	} else if len(values) != 2 {
		t.Fatalf("Expected 2 values, but received %d", len(values))
	}

	if !reflect.DeepEqual(values[0], &modelName) {
		t.Errorf("Expected model name: %s, but received: %s", modelName, values[0])
	}

	expectedInputJson := `{"prompt":"Say Hello."}`
	if !reflect.DeepEqual(values[1], &expectedInputJson) {
		t.Errorf("Expected input: %s, but received: %s", expectedInputJson, values[1])
	}
}

func TestInvokeModel_bad_model_instance(t *testing.T) {
	model := &TestModel{} // this should cause an error

	input := &TestModelInput{Prompt: "test"}
	output, err := model.Invoke(input)
	if err == nil {
		t.Error("Expected an error, but received nil")
	}
	if output != nil {
		t.Errorf("Expected output to be nil, but received: %v", output)
	}
}
