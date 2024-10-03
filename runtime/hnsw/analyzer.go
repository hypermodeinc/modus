/*
 * The code in this file originates from https://github.com/coder/hnsw
 * and is licensed under the terms of the Creative Commons Zero v1.0 Universal license
 * See the LICENSE file in the "hnsw" directory that accompanied this code for further details.
 * See also: https://github.com/coder/hnsw/blob/main/LICENSE
 *
 * SPDX-FileCopyrightText: Ammar Bandukwala <ammar@ammar.io>
 * SPDX-License-Identifier: CC0-1.0
 */

package hnsw

import "cmp"

// Analyzer is a struct that holds a graph and provides
// methods for analyzing it. It offers no compatibility guarantee
// as the methods of measuring the graph's health with change
// with the implementation.
type Analyzer[K cmp.Ordered] struct {
	Graph *Graph[K]
}

func (a *Analyzer[T]) Height() int {
	return len(a.Graph.layers)
}

// Connectivity returns the average number of edges in the
// graph for each non-empty layer.
func (a *Analyzer[T]) Connectivity() []float64 {
	var layerConnectivity []float64
	for _, layer := range a.Graph.layers {
		if len(layer.nodes) == 0 {
			continue
		}

		var sum float64
		for _, node := range layer.nodes {
			sum += float64(len(node.neighbors))
		}

		layerConnectivity = append(layerConnectivity, sum/float64(len(layer.nodes)))
	}

	return layerConnectivity
}

// Topography returns the number of nodes in each layer of the graph.
func (a *Analyzer[T]) Topography() []int {
	var topography []int
	for _, layer := range a.Graph.layers {
		topography = append(topography, len(layer.nodes))
	}
	return topography
}
