package main

import (
	"context"
	"log"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostQueryDQL(ctx context.Context, mod wasm.Module, pq uint32) uint32 {
	mem := mod.Memory()
	q, err := readString(mem, pq)
	if err != nil {
		log.Println("error reading query string from wasm memory:", err)
		return 0
	}

	r, err := queryDQL(ctx, q)
	if err != nil {
		log.Println("error querying Dgraph:", err)
		return 0
	}

	return writeString(ctx, mod, string(r))
}
