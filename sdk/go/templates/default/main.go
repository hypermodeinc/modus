package main

import "fmt"

func SayHello(name *string) string {

	var s string
	if name == nil {
		s = "World"
	} else {
		s = *name
	}

	return fmt.Sprintf("Hello, %s!", s)
}
