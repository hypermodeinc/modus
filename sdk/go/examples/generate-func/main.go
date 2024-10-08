package main

import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/dgraph"
)

func AlterSchema() error {
	const schema = `
	name: string @index(term) .
	description: string @index(term) .
	price: float .
	items: [uid] @reverse .

	type Product {
		name
		description
		price
	}

	type ShoppingCart {
		items
	}
	`
	return dgraph.AlterSchema("dgraph", schema)
}

//hyp:generate
type Product struct {
	Uid         string   `json:"uid,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Price       float64  `json:"price,omitempty"`
	DType       []string `json:"dgraph.type,omitempty"`
}

//hyp:generate
type ShoppingCart struct {
	Uid   string    `json:"uid,omitempty"`
	Items []Product `json:"items,omitempty"`
	DType []string  `json:"dgraph.type,omitempty"`
}
