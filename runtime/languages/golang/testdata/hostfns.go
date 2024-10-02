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

func EncodeStrings1(items []string) string {
	return *hostEncodeStrings1(&items)
}

func EncodeStrings2(items []*string) string {
	return *hostEncodeStrings2(&items)
}
