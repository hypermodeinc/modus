/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package sqlclient

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/go-sql-driver/mysql"
	"github.com/twpayne/go-geom/encoding/wkb"
	"github.com/twpayne/go-geom/encoding/wkt"
	"github.com/xo/dburl"
)

type mysqlDS struct {
	pool *sql.DB
}

func (ds *mysqlDS) Shutdown() {
	ds.pool.Close()
}

func (ds *mysqlDS) query(ctx context.Context, stmt string, params []any, execOnly bool) (*dbResponse, error) {

	tx, err := ds.pool.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error setting up a new tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			logger.Warn(ctx, err).Msg("Error rolling back transaction.")
			return
		}
	}()

	// TODO: what if connection times out and we need to retry

	response := new(dbResponse)
	if execOnly {
		if res, err := tx.ExecContext(ctx, stmt, params...); err != nil {
			return nil, err
		} else if ra, err := res.RowsAffected(); err == nil {
			response.RowsAffected = uint32(ra)
			if id, err := res.LastInsertId(); err == nil {
				response.LastInsertID = uint64(id)
			}
		}
	} else {
		rows, err := tx.QueryContext(ctx, stmt, params...)
		if err != nil {
			return nil, err
		}

		data, err := rowsToMap(rows)
		if err != nil {
			return nil, err
		}

		response.Result = data
		response.RowsAffected = uint32(len(data))
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return response, nil
}

func newMysqlDS(ctx context.Context, dsName string) (*mysqlDS, error) {
	connStr, err := getConnectionString(ctx, dsName, manifest.ConnectionTypeMysql)
	if err != nil {
		return nil, err
	}

	// Note: We use "dburl" to support URL-Like connection strings.
	// See:
	//  - https://github.com/xo/dburl/blob/master/README.md
	//  - https://dev.mysql.com/doc/refman/8.4/en/connecting-using-uri-or-key-value-pairs.html#connecting-using-uri
	url, err := dburl.Parse(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mysql connection URL: %w", err)
	}
	cfg, err := mysql.ParseDSN(url.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mysql connection DSN: %w", err)
	}

	// Set some default options that make sense for Modus apps.
	// See: https://github.com/go-sql-driver/mysql/blob/master/README.md#parameters
	q := url.Query()
	if !q.Has("tls") {
		cfg.TLSConfig = "preferred"
	}
	if !q.Has("clientFoundRows") {
		cfg.ClientFoundRows = true
	}

	conn, err := mysql.NewConnector(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create mysql connector: %w", err)
	}

	// Note: this doesn't actually open a connection to the database.
	// It establishes a connection pool that will be used to create connections.
	db := sql.OpenDB(conn)

	// TODO: We may want to make these configurable.
	// Reference: https://github.com/go-sql-driver/mysql/blob/master/README.md#important-settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Try to connect to the database to ensure the connection is valid.
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect mysql database [%s]: %w", dsName, err)
	}

	return &mysqlDS{db}, nil
}

func rowsToMap(rows *sql.Rows) ([]map[string]any, error) {
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("error getting column types: %w", err)
	}

	rawValues := make([]any, len(colTypes))
	for i, ct := range colTypes {
		switch ct.DatabaseTypeName() {
		case "DATE", "DATETIME", "TIMESTAMP":
			rawValues[i] = new([]byte)
		default:
			rawValues[i] = reflect.New(ct.ScanType()).Interface()
		}
	}

	var data []map[string]any
	for rows.Next() {
		if err := rows.Scan(rawValues...); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		m := make(map[string]any, len(colTypes))
		for i, ct := range colTypes {
			switch ct.DatabaseTypeName() {
			case "DATE":
				b := *rawValues[i].(*[]byte)
				if b == nil {
					m[ct.Name()] = nil
					continue
				}
				m[ct.Name()] = string(b)
			case "DATETIME":
				b := *rawValues[i].(*[]byte)
				if b == nil {
					m[ct.Name()] = nil
					continue
				}
				b[10] = 'T' // replace space with T
				m[ct.Name()] = string(b)
			case "TIMESTAMP":
				b := *rawValues[i].(*[]byte)
				if b == nil {
					m[ct.Name()] = nil
					continue
				}
				b[10] = 'T'        // replace space with T
				b = append(b, 'Z') // add Z to indicate UTC
				m[ct.Name()] = string(b)
			case "GEOMETRY":
				b := *rawValues[i].(*[]byte)
				if b == nil {
					m[ct.Name()] = nil
					continue
				}
				geom, err := wkb.Unmarshal(b[4:]) // skip the 4-byte SRID
				if err != nil {
					return nil, fmt.Errorf("error unmarshalling geometry from WKB: %w", err)
				}
				wkt, err := wkt.Marshal(geom)
				if err != nil {
					return nil, fmt.Errorf("error marshalling geometry to WKT: %w", err)
				}
				m[ct.Name()] = wkt
			default:
				m[ct.Name()] = utils.DereferencePointer(rawValues[i])
			}
		}
		data = append(data, m)
	}

	return data, nil
}
