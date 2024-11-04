/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/graphql"
)

/*
 * This application uses a public GraphQL API to retrieve information about countries.
 * We are using the API available at: https://github.com/trevorblades/countries
 * The `graphql` module from the Modus SDK allows us to interact with the API.
 */

// Define the host name of the API to be used in the requests.
// Must match one of the http connections defined in the modus.json manifest file.

const hostName = "countries-api"

// Function to list all countries returned by the API
func Countries() ([]*Country, error) {
	statement := `
	  query {
		countries{
		  code
		  name
		  capital
		}
	  }
	`
	// Execute the GraphQL query using the host name and query statement
	// The API returns a response of type `CountriesResponse` containing an array of `Country` objects
	response, err := graphql.Execute[CountriesResponse](hostName, statement, nil)

	// If err is returned (i.e., the API call fails), return the error
	if err != nil {
		return nil, err
	}
	// Otherwise, return the list of countries from the response
	return response.Data.Countries, nil

}

// Function to retrieve a specific country by its unique code
func GetCountryByCode(code string) (*Country, error) {
	statement := `
	query ($code: ID!){
      country(code: $code){
        code
        name
        capital
      }
    }
	`
	// Create a vars map to pass in the `code` parameter
	vars := map[string]any{
		"code": code,
	}

	response, err := graphql.Execute[CountryResponse](hostName, statement, vars)
	if err != nil {
		return nil, err
	}
	return response.Data.Country, nil
}
