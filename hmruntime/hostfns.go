package main

import (
	"context"
	"log"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostExecuteDQL(ctx context.Context, mod wasm.Module, pStmt uint32, isMutation uint32) uint32 {
	mem := mod.Memory()
	stmt, err := readString(mem, pStmt)
	if err != nil {
		log.Println("error reading DQL statement from wasm memory:", err)
		return 0
	}

	r, err := executeDQL(ctx, stmt, isMutation != 0)
	if err != nil {
		log.Println("error executing DQL statement:", err)
		return 0
	}

	return writeString(ctx, mod, string(r))
}

func hostExecuteGQL(ctx context.Context, mod wasm.Module, pStmt uint32) uint32 {
	mem := mod.Memory()
	stmt, err := readString(mem, pStmt)
	if err != nil {
		log.Println("error reading GraphQL string from wasm memory:", err)
		return 0
	}

	r, err := executeGQL(ctx, stmt)
	if err != nil {
		log.Println("error executing GraphQL operation:", err)
		return 0
	}

	return writeString(ctx, mod, string(r))
}
