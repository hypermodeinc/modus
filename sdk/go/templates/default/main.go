package main

// this is struct docs
type Field struct {
	// this is struct name docs
	Name string `json:"name"`
	// this is struct type docs
	Type string `json:"type"`
}

// this is some documentation
func SayHello(name *string) Field {
	return Field{
		Name: "foo",
		Type: "bar",
	}
}
