package main

import (
	"fmt"

	_ "github.com/hypermodeinc/modus/sdk/go"
)

func SayHello(name *string) string {

	var s string
	if name == nil {
		s = "World"
	} else {
		s = *name
	}

	return fmt.Sprintf("Hello, %s!", s)
}
