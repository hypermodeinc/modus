/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hypermodeAI/functions-go/pkg/dgraph"
)

const hostName = "dgraph"

func DropAll() string {
	err := dgraph.DropAll(hostName)
	if err != nil {
		return err.Error()
	}

	return "Success"
}

func DropAttr(attr string) string {
	err := dgraph.DropAttr(hostName, attr)
	if err != nil {
		return err.Error()
	}

	return "Success"
}

func AlterSchema() string {
	schema := `
	firstName: string @index(term) .
	lastName: string @index(term) .
	dgraph.type: [string] @index(exact) .
  
	type Person {
		firstName
		lastName
	}
	`
	err := dgraph.AlterSchema(hostName, schema)
	if err != nil {
		return err.Error()
	}

	return "Success"
}

func QueryPeople() ([]*Person, error) {
	query := `
	{
	  people(func: type(Person)) {
		uid
		firstName
		lastName
		dgraph.type
	  }
	}
	`

	response, err := dgraph.Execute(hostName, &dgraph.Request{
		Query: &dgraph.Query{
			Query: query,
		},
	})
	if err != nil {
		return nil, err
	}

	var peopleData PeopleData
	if err := json.Unmarshal([]byte(response.Json), &peopleData); err != nil {
		return nil, err
	}

	return peopleData.People, nil
}

func QuerySpecificPerson(firstName, lastName string) (*Person, error) {
	statement := `
	query queryPerson($firstName: string, $lastName: string) {
	  people(func: eq(firstName, $firstName)) @filter(eq(lastName, $lastName)) {
		  uid
		  firstName
		  lastName
		  dgraph.type
	  }
  }
	`

	variables := map[string]string{
		"$firstName": firstName,
		"$lastName":  lastName,
	}

	response, err := dgraph.Execute(hostName, &dgraph.Request{
		Query: &dgraph.Query{
			Query:     statement,
			Variables: variables,
		},
	})
	if err != nil {
		return nil, err
	}

	var peopleData PeopleData
	if err := json.Unmarshal([]byte(response.Json), &peopleData); err != nil {
		return nil, err
	}

	if len(peopleData.People) == 0 {
		return nil, nil // Person not found
	}

	return peopleData.People[0], nil
}

func AddPersonAsRDF(firstName, lastName string) (map[string]string, error) {
	mutation := fmt.Sprintf(`
  _:user1 <firstName> "%s" .
  _:user1 <lastName> "%s" .
  _:user1 <dgraph.type> "Person" .
  `, firstName, lastName)

	response, err := dgraph.Execute(hostName, &dgraph.Request{
		Mutations: []*dgraph.Mutation{
			{
				SetNquads: mutation,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return response.Uids, nil
}

func AddPersonAsJSON(firstName, lastName string) (map[string]string, error) {
	person := Person{
		Uid:       "_:user1",
		FirstName: firstName,
		LastName:  lastName,
		DType:     []string{"Person"},
	}

	data, err := json.Marshal(person)
	if err != nil {
		return nil, err
	}

	response, err := dgraph.Execute(hostName, &dgraph.Request{
		Mutations: []*dgraph.Mutation{
			{
				SetJson: string(data),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return response.Uids, nil
}

func UpsertPerson(nameToChangeFrom, nameToChangeTo string) (map[string]string, error) {
	query := fmt.Sprintf(`
	query {
	  person as var(func: eq(firstName, "%s"))
	`, nameToChangeFrom)

	mutation := fmt.Sprintf(`
	uid(person) <firstName> "%s" .
	`, nameToChangeTo)

	response, err := dgraph.Execute(hostName, &dgraph.Request{
		Query: &dgraph.Query{
			Query: query,
		},
		Mutations: []*dgraph.Mutation{
			{
				SetNquads: mutation,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return response.Uids, nil
}
