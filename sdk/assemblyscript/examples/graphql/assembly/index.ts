/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { graphql } from "@hypermode/modus-sdk-as";
import { Country, CountriesResponse, CountryResponse } from "./classes";

/*
 * This application uses a public GraphQL API to retrieve information about countries, continents, and languages.
 * We are using the API available at: https://github.com/trevorblades/countries
 * The `graphql` module from the Modus SDK allows us to interact with the API.
 */

// Define the host name of the API to be used in the requests. Should match one defined in the modus.json manifest file.

const hostName: string = "countries-api";

// Function to list all countries returned by the API
export function countries(): Country[] | null {
  const statement = `
    query {
      countries{
        code
        name
        capital
      }
    }
  `;
  // Execute the GraphQL query using the host name and query statement
  // The API returns a response of type `CountriesResponse` containing an array of `Country` objects
  const response = graphql.execute<CountriesResponse>(hostName, statement);

  // If no data is returned (i.e., the API call fails), return an empty array
  if (!response.data) return [];

  // Otherwise, return the list of countries from the response
  return response.data!.countries;
}

// Function to retrieve a specific country by its unique code
export function getCountryByCode(code: string): Country | null {
  // Construct the GraphQL query with a parameter for the country code
  const statement = `
    query ($code: ID!){
      country(code: $code){
        code
        name
        capital
      }
    }
  `;

  // Create a new instance of Variables to pass in the `code` parameter
  const vars = new graphql.Variables();
  vars.set("code", code);
  // Execute the query, passing the variable along with the statement and host name
  const response = graphql.execute<CountryResponse>(hostName, statement, vars);

  // If no data is returned (i.e., the API call fails), return null
  if (!response.data) return null;

  // Otherwise, return the country information from the response
  return response.data!.country;
}
