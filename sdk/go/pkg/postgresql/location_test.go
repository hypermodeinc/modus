/*
 * Copyright 2024 Hypermode, Inc.
 */

package postgresql_test

import (
	"encoding/json"
	"testing"

	"github.com/hypermodeAI/functions-go/pkg/postgresql"
)

func TestLocationString(t *testing.T) {
	location := postgresql.NewLocation(12.345678901234567, -56.7890123456789)
	expected := "(12.345678901234567,-56.7890123456789)"
	result := location.String()

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestLocationMarshalJSON(t *testing.T) {
	location := postgresql.NewLocation(12.345678901234567, -56.7890123456789)
	expected := `"(12.345678901234567,-56.7890123456789)"`
	result, err := json.Marshal(location)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if string(result) != expected {
		t.Errorf("Expected %s, but got %s", expected, string(result))
	}
}

func TestLocationUnmarshalJSON(t *testing.T) {
	data := []byte(`"(12.345678901234567,-56.7890123456789)"`)
	expected := postgresql.NewLocation(12.345678901234567, -56.7890123456789)
	location := &postgresql.Location{}

	err := json.Unmarshal(data, location)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if location.Longitude != expected.Longitude || location.Latitude != expected.Latitude {
		t.Errorf("Expected %+v, but got %+v", expected, location)
	}
}

func TestNewLocation(t *testing.T) {
	longitude := 12.345678901234567
	latitude := -56.7890123456789
	location := postgresql.NewLocation(longitude, latitude)

	if location.Longitude != longitude || location.Latitude != latitude {
		t.Errorf("Expected Longitude: %f, Latitude: %f, but got %+v", longitude, latitude, location)
	}
}

func TestParseLocation(t *testing.T) {
	s := "(12.345678901234567,-56.7890123456789)"
	expected := postgresql.NewLocation(12.345678901234567, -56.7890123456789)
	location, err := postgresql.ParseLocation(s)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if location.Longitude != expected.Longitude || location.Latitude != expected.Latitude {
		t.Errorf("Expected %+v, but got %+v", expected, location)
	}
}

func TestParseLocationError(t *testing.T) {
	s := "invalid-location"
	_, err := postgresql.ParseLocation(s)

	if err == nil {
		t.Error("Expected an error, but received none")
	}
}
