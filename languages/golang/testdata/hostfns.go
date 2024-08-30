package main

func Add(a, b int) int {
	return hostAdd(a, b)
}

func Echo(message string) string {
	return *hostEcho(&message)
}
