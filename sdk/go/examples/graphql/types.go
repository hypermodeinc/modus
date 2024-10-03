package main

type Person struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type PeopleData struct {
	People []*Person `json:"people"`
}
