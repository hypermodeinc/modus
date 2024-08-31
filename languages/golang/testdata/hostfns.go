package main

func Add(a, b int) int {
	return hostAdd(a, b)
}

func Echo1(message string) string {
	return *hostEcho1(&message)
}

func Echo2(message string) string {
	return *hostEcho2(&message)
}

func Echo3(message string) string {
	return *hostEcho3(&message)
}

func Echo4(message string) string {
	return *hostEcho4(&message)
}
