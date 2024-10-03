package main

import "github.com/hypermodeAI/functions-go/pkg/utils"

// The Product struct and the sample product will be used in the some of the examples.

type Product struct {
	Id          string  `json:"id,omitempty"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

var sampleProduct = Product{
	Id:          "123",
	Name:        "Shoes",
	Price:       50.0,
	Description: "Great shoes for walking.",
}

// since we'll be using this sample data often, we'll serialize it once and store it as a string.
var sampleProductJson string = func() string {
	bytes, _ := utils.JsonSerialize(sampleProduct)
	return string(bytes)
}()
