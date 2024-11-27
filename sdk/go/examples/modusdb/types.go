/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

type Person struct {
	Uid       string `json:"uid,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	// DType     []string `json:"dgraph.type,omitempty"`
}

type PeopleData struct {
	People []*Person `json:"people"`
}
