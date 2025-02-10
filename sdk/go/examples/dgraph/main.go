/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/dgraph"
)

const connection = "dgraph"

func DropAll() string {
	err := dgraph.DropAll(connection)
	if err != nil {
		return err.Error()
	}

	return "Success"
}

func DropAttr(attr string) string {
	err := dgraph.DropAttr(connection, attr)
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
	err := dgraph.AlterSchema(connection, schema)
	if err != nil {
		return err.Error()
	}

	return "Success"
}

func QueryPeople() ([]*Person, error) {
	query := dgraph.NewQuery(`
		{
			people(func: type(Person)) {
				uid
				firstName
				lastName
				dgraph.type
			}
		}
		`)

	response, err := dgraph.ExecuteQuery(connection, query)
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
	query := dgraph.NewQuery(`
		query queryPerson($firstName: string, $lastName: string) {
			people(func: eq(firstName, $firstName)) @filter(eq(lastName, $lastName)) {
				uid
				firstName
				lastName
				dgraph.type
			}
		}
		`).
		WithVariable("$firstName", firstName).
		WithVariable("$lastName", lastName)

	response, err := dgraph.ExecuteQuery(connection, query)
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
	mutation := dgraph.NewMutation().WithSetNquads(fmt.Sprintf(`
		_:user1 <firstName> "%s" .
		_:user1 <lastName> "%s" .
		_:user1 <dgraph.type> "Person" .
		`, dgraph.EscapeRDF(firstName), dgraph.EscapeRDF(lastName)))

	response, err := dgraph.ExecuteMutations(connection, mutation)
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

	personJson, err := json.Marshal(person)
	if err != nil {
		return nil, err
	}

	mutation := dgraph.NewMutation().WithSetJson(string(personJson))

	response, err := dgraph.ExecuteMutations(connection, mutation)
	if err != nil {
		return nil, err
	}

	return response.Uids, nil
}

func UpdatePerson(nameToChangeFrom, nameToChangeTo string) (map[string]string, error) {
	query := dgraph.NewQuery(`
		query findPerson($name: string) {
			person as var(func: eq(firstName, $name))
		}
		`).WithVariable("$name", nameToChangeFrom)

	mutation := dgraph.NewMutation().WithSetNquads(
		fmt.Sprintf(`uid(person) <firstName> "%s" .`, dgraph.EscapeRDF(nameToChangeTo)),
	)

	response, err := dgraph.ExecuteQuery(connection, query, mutation)
	if err != nil {
		return nil, err
	}

	return response.Uids, nil
}
