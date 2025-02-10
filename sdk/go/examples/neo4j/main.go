/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/neo4j"
)

// The name of the PostgreSQL connection, as specified in the modus.json manifest
const connection = "my-database"

func CreatePeopleAndRelationships() (string, error) {
	people := []map[string]any{
		{"name": "Alice", "age": 42, "friends": []string{"Bob", "Peter", "Anna"}},
		{"name": "Bob", "age": 19},
		{"name": "Peter", "age": 50},
		{"name": "Anna", "age": 30},
	}

	for _, person := range people {
		_, err := neo4j.ExecuteQuery(connection,
			"MERGE (p:Person {name: $person.name, age: $person.age})",
			map[string]any{"person": person})
		if err != nil {
			return "", err
		}
	}

	for _, person := range people {
		if person["friends"] != "" {
			_, err := neo4j.ExecuteQuery(connection, `
				MATCH (p:Person {name: $person.name})
                UNWIND $person.friends AS friend_name
                MATCH (friend:Person {name: friend_name})
                MERGE (p)-[:KNOWS]->(friend)
			`, map[string]any{
				"person": person,
			})
			if err != nil {
				return "", err
			}
		}
	}

	return "People and relationships created successfully", nil
}

type Person struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

func GetAliceFriendsUnder40() ([]Person, error) {
	response, err := neo4j.ExecuteQuery(connection, `
        MATCH (p:Person {name: $name})-[:KNOWS]-(friend:Person)
        WHERE friend.age < $age
        RETURN friend
        `,
		map[string]any{
			"name": "Alice",
			"age":  40,
		},
		neo4j.WithDbName("neo4j"),
	)
	if err != nil {
		return nil, err
	}

	nodeRecords := make([]Person, len(response.Records))

	for i, record := range response.Records {
		node, _ := neo4j.GetRecordValue[neo4j.Node](record, "friend")
		name, err := neo4j.GetProperty[string](&node, "name")
		if err != nil {
			return nil, err
		}
		age, err := neo4j.GetProperty[int64](&node, "age")
		if err != nil {
			return nil, err
		}
		nodeRecords[i] = Person{
			Name: name,
			Age:  age,
		}
	}

	return nodeRecords, nil
}

func GetAliceFriendsUnder40Ages() ([]int64, error) {
	response, err := neo4j.ExecuteQuery(connection, `
        MATCH (p:Person {name: $name})-[:KNOWS]-(friend:Person)
        WHERE friend.age < $age
        RETURN friend.age AS age
        `, map[string]any{
		"name": "Alice",
		"age":  40,
	},
		neo4j.WithDbName("neo4j"),
	)
	if err != nil {
		return nil, err
	}

	ageRecords := make([]int64, len(response.Records))

	for i, record := range response.Records {
		age, _ := neo4j.GetRecordValue[int64](record, "age")
		ageRecords[i] = age
	}

	return ageRecords, nil
}

func DeleteAllNodes() (string, error) {
	_, err := neo4j.ExecuteQuery(connection, `
		MATCH (n)
		DETACH DELETE n
	`, nil)
	if err != nil {
		return "", err
	}

	return "All nodes deleted", nil
}
