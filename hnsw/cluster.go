package hnsw

import (
	"cmp"
	"fmt"
	"math"
	"sort"

	"github.com/viterin/vek/vek32"
)

// cluster represents a cluster of nodes with a centroid.
type cluster[K cmp.Ordered] struct {
	Centroid Vector
	Nodes    []*layerNode[K]
}

// calculateCentroid calculates the centroid of a given set of vectors.
func calculateCentroid(vectors []Vector) (Vector, error) {
	if len(vectors) == 0 {
		return nil, fmt.Errorf("no vectors provided for centroid calculation")
	}

	dim := len(vectors[0])
	centroid := make(Vector, dim)
	for _, vec := range vectors {
		if len(vec) != dim {
			return nil, fmt.Errorf("dimension mismatch")
		}
		for i, val := range vec {
			centroid[i] += val
		}
	}

	// Compute the average.
	vek32.DivNumber_Inplace(centroid, float32(len(vectors)))

	return centroid, nil
}

func (g *Graph[K]) KMeansClusters(k int, maxIterations int, tolerance float32) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	var err error
	g.clusters, err = g.kMeansClusters(k, maxIterations, tolerance)
	return err
}

// kMeansClusters performs K-Means clustering on the graph using K-Means++ initialization.
func (g *Graph[K]) kMeansClusters(k int, maxIterations int, tolerance float32) ([]cluster[K], error) {
	if g.Len() == 0 {
		return nil, fmt.Errorf("graph is empty, cannot perform clustering")
	}

	if k <= 0 {
		return nil, fmt.Errorf("invalid number of clusters: %d", k)
	}

	if k > len(g.layers[0].nodes) {
		return nil, fmt.Errorf("number of clusters exceeds the number of nodes in the graph")
	}

	// Initialize centroids using K-Means++ method
	centroids, err := g.initializeCentroidsKMeansPP(k)
	if err != nil {
		return nil, fmt.Errorf("error initializing centroids: %v", err)
	}

	clusters := make([]cluster[K], k)
	prevCentroids := make([]Vector, k)

	// Step 2: Perform K-Means iterations
	for it := 0; it < maxIterations; it++ {
		// Clear nodes from each cluster
		for i := range clusters {
			clusters[i].Nodes = nil
		}

		// Step 2.1: Assign nodes to the nearest centroid
		for _, node := range g.layers[0].nodes {
			closestCluster := 0
			closestDistance := float32(math.Inf(1))

			for i, centroid := range centroids {
				dist, err := g.Distance(node.Value, centroid)
				if err != nil {
					return nil, fmt.Errorf("failed to calculate distance: %v", err)
				}

				if dist < closestDistance {
					closestDistance = dist
					closestCluster = i
				}
			}

			clusters[closestCluster].Nodes = append(clusters[closestCluster].Nodes, node)
		}

		// Step 2.2: Update centroids
		for i := range clusters {
			vectors := make([]Vector, len(clusters[i].Nodes))
			for j, node := range clusters[i].Nodes {
				vectors[j] = node.Value
			}

			newCentroid, err := calculateCentroid(vectors)
			if err != nil {
				return nil, err
			}

			clusters[i].Centroid = newCentroid
			prevCentroids[i] = centroids[i]
			centroids[i] = newCentroid
		}

		// Step 2.3: Check for convergence
		converged := true
		for i := range clusters {
			dist, err := g.Distance(centroids[i], prevCentroids[i])
			if err != nil {
				return nil, fmt.Errorf("failed to calculate distance: %v", err)
			}

			if dist > tolerance {
				converged = false
				break
			}
		}

		if converged {
			break
		}

	}

	// Assign centroids to each cluster for reference
	for i := range clusters {
		clusters[i].Centroid = centroids[i]
	}

	return clusters, nil
}

// initializeCentroidsKMeansPP initializes centroids using the K-Means++ method.
func (g *Graph[K]) initializeCentroidsKMeansPP(numClusters int) ([]Vector, error) {
	keys := make([]K, 0, len(g.layers[0].nodes))
	for key := range g.layers[0].nodes {
		keys = append(keys, key)
	}

	firstCentroid := g.layers[0].entry().Value
	centroids := []Vector{firstCentroid}

	// Step 2: Select remaining centroids
	for len(centroids) < numClusters {
		distances := make([]float32, len(keys))
		totalDistance := float32(0)

		// Calculate the distance from each point to the nearest centroid
		for i, key := range keys {
			node := g.layers[0].nodes[key]
			minDistance := float32(math.Inf(1))

			for _, centroid := range centroids {
				distance, err := g.Distance(node.Value, centroid)
				if err != nil {
					return nil, fmt.Errorf("failed to calculate distance: %v", err)
				}

				if distance < minDistance {
					minDistance = distance
				}
			}

			distances[i] = minDistance
			totalDistance += minDistance
		}

		// Step 3: Select the next centroid with a probability proportional to its distance
		rng := defaultRand()
		threshold := rng.Float32() * totalDistance
		cumulativeDistance := float32(0)
		for i, distance := range distances {
			cumulativeDistance += distance
			if cumulativeDistance >= threshold {
				centroids = append(centroids, g.layers[0].nodes[keys[i]].Value)
				break
			}
		}
	}

	return centroids, nil
}

// SubCluster performs sub-clustering on each cluster.
func (g *Graph[K]) SubCluster(subK int, maxIterations int) error {
	var subClusters []cluster[K]

	for _, cluster := range g.clusters {
		subCluster, err := cluster.kMeansSubCluster(subK, maxIterations)
		if err != nil {
			return err
		}
		subClusters = append(subClusters, subCluster...)
	}

	g.clusters = subClusters

	return nil
}

// kMeansSubCluster performs K-means clustering within a given cluster.
func (c *cluster[K]) kMeansSubCluster(subK int, maxIterations int) ([]cluster[K], error) {
	if len(c.Nodes) == 0 {
		return nil, fmt.Errorf("cluster is empty")
	}

	// Use the KMeansClusters function to sub-cluster the nodes.
	graph := &Graph[K]{
		layers: []*layer[K]{
			{
				nodes: make(map[K]*layerNode[K]),
			},
		},
	}

	for _, node := range c.Nodes {
		graph.layers[0].nodes[node.Key] = node
	}

	subClusters, err := graph.kMeansClusters(subK, maxIterations, 0.0001)
	if err != nil {
		return nil, err
	}

	return subClusters, nil
}

// HierarchicalClustering performs hierarchical clustering on the graph.
func (g *Graph[K]) HierarchicalClustering(threshold float32) ([]cluster[K], error) {
	// Initial clustering using K-Means
	initialClusters, err := g.kMeansClusters(g.M, 10, 0.0001)
	if err != nil {
		return nil, err
	}

	// Iteratively merge clusters based on a distance threshold
	for {
		merged := false
		newClusters := []cluster[K]{}

		for i := 0; i < len(initialClusters); i++ {
			if i == len(initialClusters)-1 {
				newClusters = append(newClusters, initialClusters[i])
				break
			}

			clusterA := initialClusters[i]
			clusterB := initialClusters[i+1]

			dist, err := g.Distance(clusterA.Centroid, clusterB.Centroid)
			if err != nil {
				return nil, err
			}

			if dist < threshold {
				// Merge clusters
				mergedNodes := append(clusterA.Nodes, clusterB.Nodes...)
				newCentroid, err := calculateCentroid([]Vector{clusterA.Centroid, clusterB.Centroid})
				if err != nil {
					return nil, err
				}

				newClusters = append(newClusters, cluster[K]{
					Centroid: newCentroid,
					Nodes:    mergedNodes,
				})

				i++ // Skip the next cluster as it's merged
				merged = true
			} else {
				newClusters = append(newClusters, clusterA)
			}
		}

		initialClusters = newClusters

		if !merged {
			break
		}
	}

	return initialClusters, nil
}

// ClusteredSearch performs a search using clusters to narrow down the search space.
func (g *Graph[K]) ClusteredSearch(near Vector, k int) ([]SearchResultNode[K], error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if len(g.clusters) == 0 {
		return nil, fmt.Errorf("no clusters available")
	}

	// Step 2: Find the closest cluster to the target vector 'near'
	closestCluster := g.clusters[0]
	closestDistance := float32(math.Inf(1))

	for _, cluster := range g.clusters {
		dist, err := g.Distance(near, cluster.Centroid)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate distance: %v", err)
		}

		if dist < closestDistance {
			closestDistance = dist
			closestCluster = cluster
		}
	}

	// Step 3: Search within the closest cluster
	var searchResults []SearchResultNode[K]
	for _, node := range closestCluster.Nodes {
		dist, err := g.Distance(near, node.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate distance: %v", err)
		}
		searchResults = append(searchResults, SearchResultNode[K]{Node: node.Node, Distance: dist})
	}

	// Step 4: Sort search results by distance
	sort.Slice(searchResults, func(i, j int) bool {
		return searchResults[i].Distance < searchResults[j].Distance
	})

	// Return top-k results
	if len(searchResults) > k {
		searchResults = searchResults[:k]
	}

	return searchResults, nil
}
