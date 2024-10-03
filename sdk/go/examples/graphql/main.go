/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"errors"

	"github.com/hypermodeAI/functions-go/pkg/graphql"
)

const hostName = "dgraphGraphql"

func QueryPeople() ([]*Person, error) {
	query := `
	query {
		people: queryPerson {
		  id
		  firstName
		  lastName
		}
	  }
	`

	response, err := graphql.Execute[PeopleData](hostName, query, nil)
	if err != nil {
		return nil, err
	}

	return response.Data.People, nil
}

func QuerySpecificPerson(firstName, lastName string) (*Person, error) {
	statement := `
	query queryPeople($firstName: String!, $lastName: String!) {
		people: queryPerson(
			first: 1,
			filter: { firstName: { eq: $firstName }, lastName: { eq: $lastName } }
		) {
			id
			firstName
			lastName
		}
	  }
	`

	vars := map[string]any{
		"firstName": firstName,
		"lastName":  lastName,
	}

	response, err := graphql.Execute[PeopleData](hostName, statement, vars)
	if err != nil {
		return nil, err
	}

	if len(response.Data.People) == 0 {
		return nil, nil // Person not found
	}

	return response.Data.People[0], nil
}

func AddPerson(firstName, lastName string) (*Person, error) {
	statement := `
	mutation addPerson($firstName: String!, $lastName: String!) {
		addPerson(input: [{ firstName: $firstName, lastName: $lastName }]) {
			people: person {
				id
				firstName
				lastName
			}
		}
	}
	`

	vars := map[string]any{
		"firstName": firstName,
		"lastName":  lastName,
	}

	response, err := graphql.Execute[PeopleData](hostName, statement, vars)
	if err != nil {
		return nil, err
	}

	if len(response.Data.People) == 0 {
		return nil, errors.New("Failed to add the person.")
	}

	return response.Data.People[0], nil
}
