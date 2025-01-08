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

import (
	"bufio"
	"cmp"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/google/renameio/v2"
)

// errorEncoder is a helper type to encode multiple values

var byteOrder = binary.LittleEndian

func binaryRead(r io.Reader, data interface{}) (int, error) {
	switch v := data.(type) {
	case *int:
		br, ok := r.(io.ByteReader)
		if !ok {
			return 0, fmt.Errorf("reader does not implement io.ByteReader")
		}

		i, err := binary.ReadVarint(br)
		if err != nil {
			return 0, err
		}

		*v = int(i)
		// TODO: this will usually overshoot size.
		return binary.MaxVarintLen64, nil

	case *string:
		var ln int
		_, err := binaryRead(r, &ln)
		if err != nil {
			return 0, err
		}

		s := make([]byte, ln)
		_, err = binaryRead(r, &s)
		*v = string(s)
		return len(s), err

	case *[]float32:
		var ln int
		_, err := binaryRead(r, &ln)
		if err != nil {
			return 0, err
		}

		*v = make([]float32, ln)
		return binary.Size(*v), binary.Read(r, byteOrder, *v)

	case io.ReaderFrom:
		n, err := v.ReadFrom(r)
		return int(n), err

	default:
		return binary.Size(data), binary.Read(r, byteOrder, data)
	}
}

func binaryWrite(w io.Writer, data any) (int, error) {
	switch v := data.(type) {
	case int:
		var buf [binary.MaxVarintLen64]byte
		n := binary.PutVarint(buf[:], int64(v))
		n, err := w.Write(buf[:n])
		return n, err
	case io.WriterTo:
		n, err := v.WriteTo(w)
		return int(n), err
	case string:
		n, err := binaryWrite(w, len(v))
		if err != nil {
			return n, err
		}
		n2, err := io.WriteString(w, v)
		if err != nil {
			return n + n2, err
		}

		return n + n2, nil
	case []float32:
		n, err := binaryWrite(w, len(v))
		if err != nil {
			return n, err
		}
		return n + binary.Size(v), binary.Write(w, byteOrder, v)

	default:
		sz := binary.Size(data)
		err := binary.Write(w, byteOrder, data)
		if err != nil {
			return 0, fmt.Errorf("encoding %T: %w", data, err)
		}
		return sz, err
	}
}

func multiBinaryWrite(w io.Writer, data ...any) (int, error) {
	var written int
	for _, d := range data {
		n, err := binaryWrite(w, d)
		written += n
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

func multiBinaryRead(r io.Reader, data ...any) (int, error) {
	var read int
	for i, d := range data {
		n, err := binaryRead(r, d)
		read += n
		if err != nil {
			return read, fmt.Errorf("reading %T at index %v: %w", d, i, err)
		}
	}
	return read, nil
}

const encodingVersion = 1

// Export writes the graph to a writer.
//
// T must implement io.WriterTo.
func (h *Graph[K]) Export(w io.Writer) error {
	distFuncName, ok := distanceFuncToName(h.Distance)
	if !ok {
		return fmt.Errorf("distance function %v must be registered with RegisterDistanceFunc", h.Distance)
	}
	_, err := multiBinaryWrite(
		w,
		encodingVersion,
		h.M,
		h.Ml,
		h.EfSearch,
		distFuncName,
	)
	if err != nil {
		return fmt.Errorf("encode parameters: %w", err)
	}
	_, err = binaryWrite(w, len(h.layers))
	if err != nil {
		return fmt.Errorf("encode number of layers: %w", err)
	}
	for _, layer := range h.layers {
		_, err = binaryWrite(w, len(layer.nodes))
		if err != nil {
			return fmt.Errorf("encode number of nodes: %w", err)
		}
		for _, node := range layer.nodes {
			_, err = multiBinaryWrite(w, node.Key, node.Value, len(node.neighbors))
			if err != nil {
				return fmt.Errorf("encode node data: %w", err)
			}

			for neighbor := range node.neighbors {
				_, err = binaryWrite(w, neighbor)
				if err != nil {
					return fmt.Errorf("encode neighbor %v: %w", neighbor, err)
				}
			}
		}
	}

	return nil
}

// Import reads the graph from a reader.
// T must implement io.ReaderFrom.
// The imported graph does not have to match the exported graph's parameters (except for
// dimensionality). The graph will converge onto the new parameters.
func (h *Graph[K]) Import(r io.Reader) error {
	var (
		version int
		dist    string
	)
	_, err := multiBinaryRead(r, &version, &h.M, &h.Ml, &h.EfSearch,
		&dist,
	)
	if err != nil {
		return err
	}

	var ok bool
	h.Distance, ok = distanceFuncs[dist]
	if !ok {
		return fmt.Errorf("unknown distance function %q", dist)
	}
	if h.Rng == nil {
		h.Rng = defaultRand()
	}

	if version != encodingVersion {
		return fmt.Errorf("incompatible encoding version: %d", version)
	}

	var nLayers int
	_, err = binaryRead(r, &nLayers)
	if err != nil {
		return err
	}

	h.layers = make([]*layer[K], nLayers)
	for i := 0; i < nLayers; i++ {
		var nNodes int
		_, err = binaryRead(r, &nNodes)
		if err != nil {
			return err
		}

		nodes := make(map[K]*layerNode[K], nNodes)
		for j := 0; j < nNodes; j++ {
			var key K
			var vec Vector
			var nNeighbors int
			_, err = multiBinaryRead(r, &key, &vec, &nNeighbors)
			if err != nil {
				return fmt.Errorf("decoding node %d: %w", j, err)
			}

			neighbors := make([]K, nNeighbors)
			for k := 0; k < nNeighbors; k++ {
				var neighbor K
				_, err = binaryRead(r, &neighbor)
				if err != nil {
					return fmt.Errorf("decoding neighbor %d for node %d: %w", k, j, err)
				}
				neighbors[k] = neighbor
			}

			node := &layerNode[K]{
				Node: Node[K]{
					Key:   key,
					Value: vec,
				},
				neighbors: make(map[K]*layerNeighborNode[K]),
			}

			nodes[key] = node
			for _, neighbor := range neighbors {
				node.neighbors[neighbor] = nil
			}
		}
		// Fill in neighbor pointers
		for _, node := range nodes {
			for key := range node.neighbors {
				dist, err := h.Distance(node.Value, nodes[key].Value)
				if err != nil {
					return fmt.Errorf("computing distance between %v and %v: %w", node.Key, key, err)
				}
				node.neighbors[key] = &layerNeighborNode[K]{node: nodes[key], distance: dist}
			}
		}
		h.layers[i] = &layer[K]{nodes: nodes}
	}

	return nil
}

// SavedGraph is a wrapper around a graph that persists
// changes to a file upon calls to Save. It is more convenient
// but less powerful than calling Graph.Export and Graph.Import
// directly.
type SavedGraph[K cmp.Ordered] struct {
	*Graph[K]
	Path string
}

// LoadSavedGraph opens a graph from a file, reads it, and returns it.
//
// If the file does not exist (i.e. this is a new graph),
// the equivalent of NewGraph is returned.
//
// It does not hold open a file descriptor, so SavedGraph can be forgotten
// without ever calling Save.
func LoadSavedGraph[K cmp.Ordered](path string) (*SavedGraph[K], error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	g := NewGraph[K]()
	if info.Size() > 0 {
		err = g.Import(bufio.NewReader(f))
		if err != nil {
			return nil, fmt.Errorf("import: %w", err)
		}
	}

	return &SavedGraph[K]{Graph: g, Path: path}, nil
}

// Save writes the graph to the file.
func (g *SavedGraph[K]) Save() error {
	tmp, err := renameio.TempFile("", g.Path)
	if err != nil {
		return err
	}
	defer func() {
		if err := tmp.Cleanup(); err != nil {
			fmt.Printf("Failed to clean up temp file: %v", err)
		}
	}()

	wr := bufio.NewWriter(tmp)
	err = g.Export(wr)
	if err != nil {
		return fmt.Errorf("exporting: %w", err)
	}

	err = wr.Flush()
	if err != nil {
		return fmt.Errorf("flushing: %w", err)
	}

	err = tmp.CloseAtomicallyReplace()
	if err != nil {
		return fmt.Errorf("closing atomically: %w", err)
	}

	return nil
}
