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
	"hmruntime/secrets"
	"hmruntime/utils"

	"github.com/hypermodeAI/manifest"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	wasm "github.com/tetratelabs/wazero/api"
)

var dsr = newDSRegistry()

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

// ShutdownPGPools shuts down all the PostgreSQL connection pools.
func ShutdownPGPools() {
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
			return nil, fmt.Errorf("host %s is not a postgresql host", dsname)
		}

		conf := info.(manifest.PostgresqlHostInfo)
		if conf.ConnStr == "" {
			return nil, fmt.Errorf("postgresql host %s has empty connString", dsname)
		}

		fullConnStr, err := secrets.ApplyHostSecretsToString(ctx, info, conf.ConnStr)
		if err != nil {
			return nil, fmt.Errorf("failed to apply secrets to connection string for host %s: %w", dsname, err)
		}

		dbpool, err := pgxpool.New(ctx, fullConnStr)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgres host %s: %w", dsname, err)
		}

		r.pgCache[dsname] = &postgresqlDS{pool: dbpool}
		return r.pgCache[dsname], nil
	}

	return nil, fmt.Errorf("postgresql host %s not found", dsname)
}

type hostQueryResponse struct {
	Error        *string
	ResultJson   *string
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

	pgResponse, err := executeQuery(ctx, hostName, dbType, statement, params)
	if err != nil {
		logger.Err(ctx, err).Msg("Error executing database query.")
		return 0
	}

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
		RowsAffected: pgResponse.RowsAffected,
	}

	if len(resultJson) > 0 {
		s := string(resultJson)
		response.ResultJson = &s
	}

	offset, err := writeResult(ctx, mod, response)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}

type pgResponse struct {
	Error        *string
	Result       any
	RowsAffected uint32
}

func executeQuery(ctx context.Context, dsname, dsType, stmt string, params []any) (*pgResponse, error) {
	switch dsType {
	case manifest.HostTypePostgresql:
		ds, err := dsr.getPGPool(ctx, dsname)
		if err != nil {
			return nil, err
		}

		return ds.query(ctx, stmt, params)

	default:
		return nil, fmt.Errorf("host %s has an unsupported type: %s", dsname, dsType)
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
		return nil, err
	}

	data, err := pgx.CollectRows(rows, pgx.RowToMap)
	if err != nil {
		return nil, err
	}

	rowsAffected := uint32(rows.CommandTag().RowsAffected())

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	response := &pgResponse{
		// Error: "",
		Result:       data,
		RowsAffected: rowsAffected,
	}

	return response, nil
}
