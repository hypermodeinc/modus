# hnsw

Package `hnsw` implements Hierarchical Navigable Small World graphs in Go. You
can read up about how they work [here](https://www.pinecone.io/learn/series/faiss/hnsw/). In essence,
they allow for fast approximate nearest neighbor searches with high-dimensional
vector data.

This package can be thought of as an in-memory alternative to your favorite
vector database (e.g. Pinecone, Weaviate). It implements just the essential
operations:

| Operation | Complexity            | Description                                  |
| --------- | --------------------- | -------------------------------------------- |
| Insert    | $O(log(n))$           | Insert a vector into the graph               |
| Delete    | $O(M^2 \cdot log(n))$ | Delete a vector from the graph               |
| Search    | $O(log(n))$           | Search for the nearest neighbors of a vector |
| Lookup    | $O(1)$                | Retrieve a vector by ID                      |

> [!NOTE]
> Complexities are approximate where $n$ is the number of vectors in the graph
> and $M$ is the maximum number of neighbors each node can have. This [paper](https://arxiv.org/pdf/1603.09320) is a good resource for understanding the effect of
> the various construction parameters.

## Usage

```
go get github.com/coder/hnsw@main
```

```go
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
g.Add(
	nodes...,
)

neighbors, err := g.Search(
	[]float32{0.5, 0.5, 0.5},
	1,
)
if err != nil {
	panic(err)
}
fmt.Printf("best friend: %v\n", neighbors[0].Value)
```

## Persistence

While all graph operations are in-memory, `hnsw` provides facilities for loading/saving from persistent storage.

For an `io.Reader`/`io.Writer` interface, use `Graph.Export` and `Graph.Import`.

If you're using a single file as the backend, hnsw provides a convenient `SavedGraph` type instead:

```go
path := "some.graph"
g1, err := LoadSavedGraph[int](path)
if err != nil {
    panic(err)
}
// Insert some vectors
for i := 0; i < 128; i++ {
    g1.Add(hnsw.MakeNode(i, []float32{float32(i)}))
}

// Save to disk
err = g1.Save()
if err != nil {
    panic(err)
}

// Later...
// g2 is a copy of g1
g2, err := LoadSavedGraph[int](path)
if err != nil {
    panic(err)
}
```

See more:

- [Export](https://pkg.go.dev/github.com/coder/hnsw#Graph.Export)
- [Import](https://pkg.go.dev/github.com/coder/hnsw#Graph.Import)
- [SavedGraph](https://pkg.go.dev/github.com/coder/hnsw#SavedGraph)

We use a fast binary encoding for the graph, so you can expect to save/load
nearly at disk speed. On my M3 Macbook I get these benchmark results:

```
goos: darwin
goarch: arm64
pkg: github.com/coder/hnsw
BenchmarkGraph_Import-16            4029            259927 ns/op         796.85 MB/s      496022 B/op       3212 allocs/op
BenchmarkGraph_Export-16            7042            168028 ns/op        1232.49 MB/s      239886 B/op       2388 allocs/op
PASS
ok      github.com/coder/hnsw   2.624s
```

when saving/loading a graph of 100 vectors with 256 dimensions.

## Performance

By and large the greatest effect you can have on the performance of the graph
is reducing the dimensionality of your data. At 1536 dimensions (OpenAI default),
70% of the query process under default parameters is spent in the distance function.

If you're struggling with slowness / latency, consider:

- Reducing dimensionality
- Increasing $M$

And, if you're struggling with excess memory usage, consider:

- Reducing $M$ a.k.a `Graph.M` (the maximum number of neighbors each node can have)
- Reducing $m_L$ a.k.a `Graph.Ml` (the level generation parameter)

## Memory Overhead

The memory overhead of a graph looks like:

$$
\displaylines{
mem_{graph} = n \cdot \log(n) \cdot \text{size(id)} \cdot M \\
mem_{base} = n \cdot d \cdot 4 \\
mem_{total} = mem_{graph} + mem_{base}
}
$$

where:

- $n$ is the number of vectors in the graph
- $\text{size(key)}$ is the average size of the key in bytes
- $M$ is the maximum number of neighbors each node can have
- $d$ is the dimensionality of the vectors
- $mem_{graph}$ is the memory used by the graph structure across all layers
- $mem_{base}$ is the memory used by the vectors themselves in the base or 0th layer

You can infer that:

- Connectivity ($M$) is very expensive if keys are large
- If $d \cdot 4$ is far larger than $M \cdot \text{size(key)}$, you should expect linear memory usage spent on representing vector data
- If $d \cdot 4$ is far smaller than $M \cdot \text{size(key)}$, you should expect $n \cdot \log(n)$ memory usage spent on representing graph structure

In the example of a graph with 256 dimensions, and $M = 16$, with 8 byte keys, you would see that each vector takes:

- $256 \cdot 4 = 1024$ data bytes
- $16 \cdot 8 = 128$ metadata bytes

and memory growth is mostly linear.
