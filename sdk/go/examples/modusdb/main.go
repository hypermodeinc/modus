/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
	"github.com/hypermodeinc/modus/sdk/go/pkg/modusdb"
)

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
	  }
	}
	`

	response, err := modusdb.Query(&query)
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
	  }
  }
	`, firstName, lastName)

	response, err := modusdb.Query(&query)
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
		// DType:     []string{"Person"},
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

func GenerateText(prompt string) (*string, error) {

	// The imported ChatModel type follows the OpenAI Chat completion model input format.
	model, err := models.GetModel[openai.ChatModel]("text-generator")
	if err != nil {
		return nil, err
	}

	input, err := model.CreateInput(
		openai.NewSystemMessage("You are a helpful assistant. Try and be as concise as possible."),
		openai.NewUserMessage(prompt),
	)
	if err != nil {
		return nil, err
	}

	input.Temperature = 0.7

	output, err := model.Invoke(input)
	if err != nil {
		return nil, err
	}

	outputStr := strings.TrimSpace(output.Choices[0].Message.Content)
	return &outputStr, nil
}
