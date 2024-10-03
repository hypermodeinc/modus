package main

import (
	"fmt"

	"github.com/hypermodeAI/functions-go/pkg/postgresql"
)

// The name of the PostgreSQL host, as specified in the hypermode.json manifest
const host = "my-database"

type Person struct {
	Id   int                  `json:"id"`
	Name string               `json:"name"`
	Age  int                  `json:"age"`
	Home *postgresql.Location `json:"home"`
}

func GetAllPeople() ([]Person, error) {
	const query = "select * from people order by id"
	rows, _, err := postgresql.Query[Person](host, query)
	return rows, err
}

func GetPeopleByName(name string) ([]Person, error) {
	const query = "select * from people where name = $1"
	rows, _, err := postgresql.Query[Person](host, query, name)
	return rows, err
}

func GetPerson(id int) (*Person, error) {
	const query = "select * from people where id = $1"
	rows, _, err := postgresql.Query[Person](host, query, id)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil // Person not found
	}

	return &rows[0], nil
}

func AddPerson(name string, age int) (*Person, error) {
	const query = "insert into people (name, age) values ($1, $2) RETURNING id"

	id, _, err := postgresql.QueryScalar[int](host, query, name, age)
	if err != nil {
		return nil, fmt.Errorf("Failed to add person to database: %v", err)
	}

	p := Person{Id: id, Name: name, Age: age}
	return &p, nil
}

func UpdatePersonHome(id int, longitude, latitude float64) (*Person, error) {
	const query = "update people set home = point($1, $2) where id = $3"

	affected, err := postgresql.Execute(host, query, longitude, latitude, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to update person in database: %v", err)
	}

	if affected != 1 {
		return nil, fmt.Errorf("Failed to update person with id %d. The record may not exist.", id)
	}

	return GetPerson(id)
}

func DeletePerson(id int) (string, error) {
	const query = "delete from people where id = $1"

	affected, err := postgresql.Execute(host, query, id)
	if err != nil {
		return "", fmt.Errorf("Failed to delete person from database: %v", err)
	}

	if affected != 1 {
		return "", fmt.Errorf("Failed to delete person with id %d. The record may not exist.", id)
	}

	return "success", nil
}
