package main

import (
	"fmt"

	"hypruntime/hnsw"
)

func main() {
	g := hnsw.NewGraph[int]()
	keys := []int{1, 2, 3}
	values := [][]float32{
		{1, 1, 1},
		{1, -1, 0.999},
		{1, 0, -0.5},
	}
	nodes, err := hnsw.MakeNodes(keys, values)
	if err != nil {
		panic(err)
	}
	err = g.Add(
		nodes...,
	)
	if err != nil {
		panic(err)
	}

	neighbors, err := g.Search(
		[]float32{0.5, 0.5, 0.5},
		1,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("best friend: %v\n", neighbors[0].Value)
}
