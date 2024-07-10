/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"
	"sync"

	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/plugins"
	"hmruntime/utils"

	"github.com/hypermodeAI/manifest"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	wasm "github.com/tetratelabs/wazero/api"
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

type hostQueryResponse struct {
	Error        string
	ResultJson   string
	RowsAffected uint32
}

func (r *hostQueryResponse) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "HostQueryResponse",
		Path: "~lib/@hypermode/functions-as/assembly/database/HostQueryResponse",
	}
}

func hostDatabaseQuery(ctx context.Context, mod wasm.Module, pHostName, pDbType, pStatement, pParamsJson uint32) uint32 {
	var hostName, dbType, statement, paramsJson string
	if err := readParams4(ctx, mod,
		pHostName, pDbType, pStatement, pParamsJson,
		&hostName, &dbType, &statement, &paramsJson); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	var params []any
	if err := utils.JsonDeserialize([]byte(paramsJson), &params); err != nil {
		logger.Err(ctx, err).Msg("Error deserializing database query parameters.")
		return 0
	}

	pgResponse := executeQuery(ctx, hostName, dbType, statement, params)

	var resultJson []byte
	if pgResponse.Result != nil {
		var err error
		resultJson, err = utils.JsonSerialize(pgResponse.Result)
		if err != nil {
			logger.Err(ctx, err).Msg("Error serializing result.")
			return 0
		}
	}

	response := hostQueryResponse{
		Error:        pgResponse.Error,
		ResultJson:   string(resultJson),
		RowsAffected: pgResponse.RowsAffected,
	}

	offset, err := writeResult(ctx, mod, response)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}

type pgResponse struct {
	Error        string `json:"error"`
	Result       any    `json:"result"`
	RowsAffected uint32 `json:"rowsAffected"`
}

func executeQuery(ctx context.Context, dsname, dsType, stmt string, params []any) *pgResponse {
	switch dsType {
	case manifest.HostTypePostgresql:
		ds, err := dsr.getPGPool(ctx, dsname)
		if err != nil {
			logger.Err(ctx, err).Msg("Error finding host in manifest.")
			return nil
		}

		result, err := ds.query(ctx, stmt, params)
		if err != nil {
			logger.Err(ctx, err).Msg("Error executing query.")
			return nil
		}

		return result

	default:
		logger.Error(ctx).
			Str("dsname", dsname).
			Str("dsType", dsType).
			Msg("Unsupported host type.")

		return nil
	}
}

func (ds *postgresqlDS) query(ctx context.Context, stmt string, params []any) (*pgResponse, error) {

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
	rows, err := tx.Query(ctx, stmt, params...)
	if err != nil {
		return nil, fmt.Errorf("error running query: %w", err)
	}

	data, err := pgx.CollectRows(rows, pgx.RowToMap)
	if err != nil {
		return nil, fmt.Errorf("error reading result: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	response := &pgResponse{
		Result: data,
		// Error: "",
		// RowsAffected: uint32(len(data)),
	}

	return response, nil
}
