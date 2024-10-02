/*
 * Copyright 2024 Hypermode, Inc.
 */

package sqlclient

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/secrets"

	"github.com/hypermodeAI/manifest"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

	for name, info := range manifestdata.GetManifest().Hosts {
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
