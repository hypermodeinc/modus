//go:build !wasip1

/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

/*
	These tests are fairly straightforward.
	They are here to demonstrate how to write tests for your functions.
*/

import (
	"testing"
	"time"

	"github.com/hypermodeinc/modus/sdk/go/pkg/console"
)

func TestLogMessage(t *testing.T) {
	msg := "Hello, World!"

	LogMessage(msg)

	values := console.LogCallStack.Pop()
	if values[0] != "" {
		t.Errorf(`LogCalls[0] = %s; want ""`, values[0])
	}
	if values[1] != msg {
		t.Errorf(`LogCalls[1] = %s; want "%s"`, values[1], msg)
	}
}

func TestAdd(t *testing.T) {
	result := Add(1, 2)
	if result != 3 {
		t.Errorf("Add(1, 2) = %d; want 3", result)
	}
}

func TestAdd3(t *testing.T) {
	a := 1
	b := 2
	c := 3
	result := Add3(a, b, &c)
	if result != 6 {
		t.Errorf("Add3(%d, %d, %d) = %d; want 6", a, b, c, result)
	}
}

func TestAddN(t *testing.T) {
	// Test case 1: Empty arguments
	result := AddN()
	if result != 0 {
		t.Errorf("AddN() = %d; want 0", result)
	}

	// Test case 2: Single argument
	result = AddN(5)
	if result != 5 {
		t.Errorf("AddN(5) = %d; want 5", result)
	}

	// Test case 3: Multiple arguments
	result = AddN(1, 2, 3, 4, 5)
	if result != 15 {
		t.Errorf("AddN(1, 2, 3, 4, 5) = %d; want 15", result)
	}

	// Test case 4: Negative numbers
	result = AddN(-1, -2, -3)
	if result != -6 {
		t.Errorf("AddN(-1, -2, -3) = %d; want -6", result)
	}
}

func TestGetCurrentTime(t *testing.T) {
	// mock the clock for testing
	now := time.Now()
	nowFunc = func() time.Time { return now }

	expected := now
	result := GetCurrentTime()
	if result != expected {
		t.Errorf("GetCurrentTime() = %v; want %v", result, expected)
	}
}

func TestGetCurrentTimeFormatted(t *testing.T) {
	// mock the clock for testing
	now := time.Now()
	nowFunc = func() time.Time { return now }

	expected := now.Format(time.DateTime)
	result := GetCurrentTimeFormatted()
	if result != expected {
		t.Errorf("GetCurrentTimeFormatted() = %s; want %s", result, expected)
	}
}

func TestGetFullName(t *testing.T) {
	firstName := "John"
	lastName := "Doe"
	expected := "John Doe"
	result := GetFullName(firstName, lastName)
	if result != expected {
		t.Errorf("GetFullName(%s, %s) = %s; want %s", firstName, lastName, result, expected)
	}
}

func TestSayHello(t *testing.T) {
	name := "Alice"
	expected := "Hello, Alice!"
	result := SayHello(&name)
	if result != expected {
		t.Errorf("SayHello(%s) = %s; want %s", name, result, expected)
	}
}

func TestGetPerson(t *testing.T) {
	expected := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       42,
	}
	result := GetPerson()
	if result != expected {
		t.Errorf("GetPerson() = %v; want %v", result, expected)
	}
}

func TestGetRandomPerson(t *testing.T) {
	result := GetRandomPerson()
	if result.FirstName == "" || result.LastName == "" || result.Age == 0 {
		t.Errorf("GetRandomPerson() returned an invalid person: %v", result)
	}
}

func TestGetPeople(t *testing.T) {
	expected := []Person{
		{
			FirstName: "Bob",
			LastName:  "Smith",
			Age:       42,
		},
		{
			FirstName: "Alice",
			LastName:  "Jones",
			Age:       35,
		},
		{
			FirstName: "Charlie",
			LastName:  "Brown",
			Age:       8,
		},
	}
	result := GetPeople()
	if len(result) != len(expected) {
		t.Errorf("GetPeople() returned %d people; want %d", len(result), len(expected))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("GetPeople()[%d] = %v; want %v", i, result[i], expected[i])
		}
	}
}

func TestGetNameAndAge(t *testing.T) {
	expectedName := "John Doe"
	expectedAge := 42
	resultName, resultAge := GetNameAndAge()
	if resultName != expectedName {
		t.Errorf("GetNameAndAge() returned name %s; want %s", resultName, expectedName)
	}
	if resultAge != expectedAge {
		t.Errorf("GetNameAndAge() returned age %d; want %d", resultAge, expectedAge)
	}
}
