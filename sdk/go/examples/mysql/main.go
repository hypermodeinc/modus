/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/mysql"
)

// The name of the MySQL host, as specified in the modus.json manifest
const host = "my-database"

type Person struct {
	Id   int             `json:"id"`
	Name string          `json:"name"`
	Age  int             `json:"age"`
	Home *mysql.Location `json:"home"`
}

func GetAllPeople() ([]Person, error) {
	const query = "select * from people order by id"
	response, err := mysql.Query[Person](host, query)
	return response.Rows, err
}

func GetPeopleByName(name string) ([]Person, error) {
	const query = "select * from people where name = ?"
	response, err := mysql.Query[Person](host, query, name)
	return response.Rows, err
}

func GetPerson(id int) (*Person, error) {
	const query = "select * from people where id = ?"
	response, err := mysql.Query[Person](host, query, id)
	if err != nil {
		return nil, err
	}

	if len(response.Rows) == 0 {
		return nil, nil // Person not found
	}

	return &response.Rows[0], nil
}

func AddPerson(name string, age int) (*Person, error) {
	const query = "insert into people (name, age) values (?, ?)"

	response, err := mysql.Execute(host, query, name, age)
	if err != nil {
		return nil, fmt.Errorf("Failed to add person to database: %v", err)
	}

	p := Person{Id: int(response.LastInsertId), Name: name, Age: age}
	return &p, nil
}

func UpdatePersonHome(id int, longitude, latitude float64) (*Person, error) {
	const query = "update people set home = point(?,?) where id = ?"

	response, err := mysql.Execute(host, query, longitude, latitude, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to update person in database: %v", err)
	}

	if response.RowsAffected != 1 {
		return nil, fmt.Errorf("Failed to update person with id %d. The record may not exist.", id)
	}

	return GetPerson(id)
}

func DeletePerson(id int) (string, error) {
	const query = "delete from people where id = ?"

	response, err := mysql.Execute(host, query, id)
	if err != nil {
		return "", fmt.Errorf("Failed to delete person from database: %v", err)
	}

	if response.RowsAffected != 1 {
		return "", fmt.Errorf("Failed to delete person with id %d. The record may not exist.", id)
	}

	return "success", nil
}
