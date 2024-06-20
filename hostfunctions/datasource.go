/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"hmruntime/logger"
	"hmruntime/manifestdata"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	wasm "github.com/tetratelabs/wazero/api"

	"github.com/hypermodeAI/manifest"
)

var (
	dsr = newDSRegistry()
)

type postgresqlDS struct {
	pool *pgxpool.Pool
}

type dsRegistry struct {
	sync.RWMutex
	pgCache map[string]*postgresqlDS
}

func newDSRegistry() *dsRegistry {
	return &dsRegistry{
		pgCache: make(map[string]*postgresqlDS),
	}
}

// ShutdownPools shuts down all the connection pools.
func ShutdownPools() {
	dsr.Lock()
	defer dsr.Unlock()

	for _, ds := range dsr.pgCache {
		ds.pool.Close()
	}

	clear(dsr.pgCache)
}

func (r *dsRegistry) getPGPool(ctx context.Context, dsname string) (*postgresqlDS, error) {
	// fast path
	r.RLock()
	ds, ok := r.pgCache[dsname]
	r.RUnlock()
	if ok {
		return ds, nil
	}

	// slow path
	r.Lock()
	defer r.Unlock()

	// we do another lookup to make sure any other goroutine didn't create it
	if ds, ok := r.pgCache[dsname]; ok {
		return ds, nil
	}

	for name, info := range manifestdata.Manifest.Hosts {
		if name != dsname {
			continue
		}

		if info.HostType() != manifest.HostTypePostgresql {
			return nil, fmt.Errorf("host [%s] is not a postgresql host", dsname)
		}

		conf := info.(manifest.PostgresqlHostInfo)
		if conf.ConnStr == "" {
			return nil, fmt.Errorf("postgresql host [%s] has empty connString", dsname)
		}

		dbpool, err := pgxpool.New(ctx, conf.ConnStr)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database %s: %w", conf.ConnStr, err)
		}

		r.pgCache[dsname] = &postgresqlDS{pool: dbpool}
		return r.pgCache[dsname], nil
	}

	return nil, fmt.Errorf("postgresql host [%s] not found", dsname)
}

func databaseQuery(ctx context.Context, mod wasm.Module, pDSName uint32, pDSType uint32, pQueryData uint32) uint32 {
	var dsname, dsType, queryData string
	if err := readParams3(ctx, mod, pDSName, pDSType, pQueryData, &dsname, &dsType, &queryData); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	response, err := json.Marshal(executeQuery(ctx, dsname, dsType, queryData))
	if err != nil {
		logger.Err(ctx, err).Msg("Error serializing response.")
		return 0
	}

	offset, err := writeResult(ctx, mod, string(response))
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}

type pgResponse struct {
	Result any    `json:"result"`
	Error  string `json:"error"`
}

func executeQuery(ctx context.Context, dsname, dsType, queryData string) *pgResponse {
	switch dsType {
	case manifest.HostTypePostgresql:
		ds, err := dsr.getPGPool(ctx, dsname)
		if err != nil {
			logger.Err(ctx, err).Msg("Error finding host in manifest.")
			return &pgResponse{Error: err.Error()}
		}

		result, err := ds.query(ctx, queryData)
		if err != nil {
			return &pgResponse{Error: err.Error()}
		}
		return &pgResponse{Result: result}

	default:
		logger.Err(ctx, fmt.Errorf("unsupported host type %s", dsType)).
			Str("dsname", dsname).
			Str("dsType", dsType).
			Msg("Unsupported host type.")
		return &pgResponse{Error: "unsupported host type"}
	}
}

type pgInput struct {
	Query  string `json:"query"`
	Params []any  `json:"params"`
}

func (ds *postgresqlDS) query(ctx context.Context, queryData string) (any, error) {
	var input pgInput
	if err := json.Unmarshal([]byte(queryData), &input); err != nil {
		return nil, fmt.Errorf("error parsing query data: %w", err)
	}

	tx, err := ds.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error setting up a new tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			logger.Err(ctx, err).Msg("Error rolling back transaction.")
			return
		}
	}()

	// TODO: what if connection times out and we need to retry
	rows, err := tx.Query(ctx, input.Query, input.Params...)
	if err != nil {
		return nil, fmt.Errorf("error running query: %w", err)
	}

	data, err := pgx.CollectRows(rows, pgx.RowToMap)
	if err != nil {
		return nil, fmt.Errorf("error reading result: %w", err)
	}
	return data, tx.Commit(ctx)
}
