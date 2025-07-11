/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db_test

import (
	"encoding/json"
	"testing"

	"github.com/hypermodeinc/modus/sdk/go/pkg/db"
)

func TestLocationString(t *testing.T) {
	location := db.NewLocation(12.345678901234567, -56.7890123456789)
	expected := "(12.345678901234567,-56.7890123456789)"
	result := location.String()

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestLocationMarshalJSON(t *testing.T) {
	location := db.NewLocation(12.345678901234567, -56.7890123456789)
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
	expected := db.NewLocation(12.345678901234567, -56.7890123456789)
	location := &db.Location{}

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
	location := db.NewLocation(longitude, latitude)

	if location.Longitude != longitude || location.Latitude != latitude {
		t.Errorf("Expected Longitude: %f, Latitude: %f, but got %+v", longitude, latitude, location)
	}
}

func TestParseLocation(t *testing.T) {
	s := "(12.345678901234567,-56.7890123456789)"
	expected := db.NewLocation(12.345678901234567, -56.7890123456789)
	location, err := db.ParseLocation(s)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if location.Longitude != expected.Longitude || location.Latitude != expected.Latitude {
		t.Errorf("Expected %+v, but got %+v", expected, location)
	}
}

func TestParseLocationError(t *testing.T) {
	s := "invalid-location"
	_, err := db.ParseLocation(s)

	if err == nil {
		t.Error("Expected an error, but received none")
	}
}
