/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/hypermodeinc/modus/sdk/go/pkg/console"
)

func LogMessage(message string) {
	console.Log(message)
}

func Add(x, y int) int {
	return x + y
}

func Add3(a, b int, c *int) int {
	if c != nil {
		return a + b + *c
	}
	return a + b
}

func AddN(args ...int) int {
	sum := 0
	for _, arg := range args {
		sum += arg
	}
	return sum
}

// this indirection is so we can mock time.Now in tests
var nowFunc = time.Now

func GetCurrentTime() time.Time {
	return nowFunc()
}

func GetCurrentTimeFormatted() string {
	return nowFunc().Format(time.DateTime)
}

func GetFullName(firstName, lastName string) string {
	return firstName + " " + lastName
}

func SayHello(name *string) string {
	if name == nil {
		return "Hello!"
	} else {
		return "Hello, " + *name + "!"
	}
}

type Person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

func GetPerson() Person {
	return Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       42,
	}
}

func GetRandomPerson() Person {
	people := GetPeople()
	i := rand.Intn(len(people))
	return people[i]
}

func GetPeople() []Person {
	return []Person{
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
}

func GetNameAndAge() (name string, age int) {
	p := GetPerson()
	return GetFullName(p.FirstName, p.LastName), p.Age
}

// This is the preferred way to handle errors in functions.
// Simply declare an error interface as the last return value.
// You can use any object that implements the Go error interface.
// For example, you can create a new error with errors.New("message"),
// or with fmt.Errorf("message with %s", "parameters").
func TestNormalError(input string) (string, error) {
	if input == "" {
		return "", errors.New("input is empty")
	}
	output := "You said: " + input
	return output, nil
}

// This is an alternative way to handle errors in functions.
// It is identical in behavior to TestNormalError, but is not Go idiomatic.
func TestAlternativeError(input string) string {
	if input == "" {
		console.Error("input is empty")
		return ""
	}
	output := "You said: " + input
	return output
}

// This panics, will log the message as "fatal" and exits the function.
// Generally, you should not panic.
func TestPanic() {
	panic("This is a message from a panic.\nThis is a second line from a panic.\n")
}

// If you need to exit prematurely without panicking, you can use os.Exit.
// However, you cannot return any values from the function, so if you want
// to log an error message, you should do so before calling os.Exit.
// The exit code should be 0 for success, and non-zero for failure.
func TestExit() {
	console.Error("This is an error message.")
	os.Exit(1)
	println("This line will not be executed.")
}

func TestLogging() {
	// This is a simple log message. It has no level.
	console.Log("This is a simple log message.")

	// These messages are logged at different levels.
	console.Debug("This is a debug message.")
	console.Info("This is an info message.")
	console.Warn("This is a warning message.")

	// This logs an error message, but allows the function to continue.
	console.Error("This is an error message.")
	console.Error(
		`This is line 1 of a multi-line error message.
  This is line 2 of a multi-line error message.
  This is line 3 of a multi-line error message.`,
	)

	// You can also use Go's built-in printing commands.
	println("This is a println message.")
	fmt.Println("This is a fmt.Println message.")
	fmt.Printf("This is a fmt.Printf message (%s).\n", "with a parameter")

	// You can even just use stdout/stderr and the log level will be "info" or "error" respectively.
	fmt.Fprintln(os.Stdout, "This is an info message printed to stdout.")
	fmt.Fprintln(os.Stderr, "This is an error message printed to stderr.")

	// NOTE: The main difference between using console functions and Go's built-in functions, is in how newlines are handled.
	// - Using console functions, newlines are preserved and the entire message is logged as a single message.
	// - Using Go's built-in functions, each line is logged as a separate message.
	//
	// Thus, if you are logging data for debugging, we highly recommend using the console functions
	// to keep the data together in a single log message.
	//
	// The console functions also allow you to better control the reported logging level.
}
