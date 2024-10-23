/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

type Country struct {
	Code    string `json:"code,omitempty"`
	Name    string `json:"name,omitempty"`
	Capital string `json:"capital,omitempty"`
}

type CountriesResponse struct {
	Countries []*Country `json:"countries"`
}
type CountryResponse struct {
	Country *Country `json:"country"`
}
