package main

type Person struct {
	Uid       string   `json:"uid,omitempty"`
	FirstName string   `json:"firstName,omitempty"`
	LastName  string   `json:"lastName,omitempty"`
	DType     []string `json:"dgraph.type,omitempty"`
}

type PeopleData struct {
	People []*Person `json:"people"`
}
