/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/modusdb"
)

func DropAll() string {
	err := modusdb.DropAll()
	if err != nil {
		return err.Error()
	}

	return "Success"
}

func AlterSchema() string {
	schema := `
	firstName: string @index(term) .
	lastName: string @index(term) .
	modusdb.type: [string] @index(exact) .
  
	type Person {
		firstName
		lastName
	}
	`
	err := modusdb.AlterSchema(schema)
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
		modusdb.type
	  }
	}
	`

	response, err := modusdb.Query(query)
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
	query := fmt.Sprintf(`
	query queryPerson {
	  people(func: eq(firstName, "%v")) @filter(eq(lastName, "%v")) {
		  uid
		  firstName
		  lastName
		  modusdb.type
	  }
  }
	`, firstName, lastName)

	response, err := modusdb.Query(query)
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

func AddPersonAsRDF(firstName, lastName string) (*map[string]uint64, error) {
	mutation := fmt.Sprintf(`
  _:user1 <firstName> "%s" .
  _:user1 <lastName> "%s" .
  _:user1 <modusdb.type> "Person" .
  `, firstName, lastName)

	response, err := modusdb.Mutate(&modusdb.MutationRequest{
		Mutations: []*modusdb.Mutation{
			{
				SetNquads: mutation,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func AddPersonAsJSON(firstName, lastName string) (*map[string]uint64, error) {
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

	response, err := modusdb.Mutate(&modusdb.MutationRequest{
		Mutations: []*modusdb.Mutation{
			{
				SetJson: string(data),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
