/*
* Copyright 2024 Hypermode, Inc.
 */

package dqlclient

import (
	"context"
	"fmt"
	"hmruntime/manifestdata"
	"hmruntime/secrets"
	"sync"

	"github.com/hypermodeAI/manifest"

	"github.com/dgraph-io/dgo/v230"
	"github.com/dgraph-io/dgo/v230/protos/api"
)

var dgr = newDgraphRegistry()

type dgraphRegistry struct {
	sync.RWMutex
	dgraphConnectorCache map[string]*dgraphConnector
}

func newDgraphRegistry() *dgraphRegistry {
	return &dgraphRegistry{
		dgraphConnectorCache: make(map[string]*dgraphConnector),
	}
}

func ShutdownConns() {
	dgr.Lock()
	defer dgr.Unlock()
	for _, ds := range dgr.dgraphConnectorCache {
		ds.conn.Close()
	}
	clear(dgr.dgraphConnectorCache)
}

func (dr *dgraphRegistry) getDgraphConnector(ctx context.Context, dgName string) (*dgraphConnector, error) {
	dr.RLock()
	ds, ok := dr.dgraphConnectorCache[dgName]
	dr.RUnlock()
	if ok {
		return ds, nil
	}

	dr.Lock()
	defer dr.Unlock()

	if ds, ok := dr.dgraphConnectorCache[dgName]; ok {
		return ds, nil
	}

	for name, info := range manifestdata.GetManifest().Hosts {
		if name != dgName {
			continue
		}

		if info.HostType() != manifest.HostTypeDgraphCloud {
			return nil, fmt.Errorf("host %s is not a dgraph cloud host", dgName)
		}

		host := info.(manifest.DgraphCloudHostInfo)
		if host.Endpoint == "" {
			return nil, fmt.Errorf("dgraph host %s has empty address", dgName)
		}

		confKey, err := secrets.ApplyHostSecretsToString(ctx, info, host.Key)

		conn, err := dgo.DialCloud(host.Endpoint, confKey)
		if err != nil {
			return nil, err
		}

		ds := &dgraphConnector{
			conn:     conn,
			dgClient: dgo.NewDgraphClient(api.NewDgraphClient(conn)),
		}
		dr.dgraphConnectorCache[dgName] = ds
		return ds, nil
	}

	return nil, fmt.Errorf("dgraph host %s not found", dgName)
}
