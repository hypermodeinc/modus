/*
 * Copyright 2024 Hypermode, Inc.
 */

package metrics_test

import (
	"bytes"
	"context"
	"fmt"
	"hmruntime/config"
	"hmruntime/server"
	"io"
	"net/http"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/prometheus/common/expfmt"
)

const (
	configPort = 8080
)

var (
	graphqlEndpoint = fmt.Sprintf("http://localhost:%d/graphql", configPort)
	adminEndpoint   = fmt.Sprintf("http://localhost:%d/admin", configPort)
	healthEndpoint  = fmt.Sprintf("http://localhost:%d/health", configPort)
	metricsEndpoint = fmt.Sprintf("http://localhost:%d/metrics", configPort)
)

func setupRuntime(t *testing.T) {
	// configure the port to listen on
	config.Port = configPort

	// configure a wasm plugin
	_, thisFilePath, _, _ := runtime.Caller(0)
	config.StoragePath = path.Join(thisFilePath, "..", "testutil", "data", "test-as")

	go func() {
		if err := server.Start(context.Background()); err != nil {
			t.Logf("error while shutting down the server: %v", err)
		}
	}()

	for i := 0; i < 10; i++ {
		resp, err := http.Get(metricsEndpoint)
		if err == nil {
			resp.Body.Close()
			break
		}
		time.Sleep(1 * time.Second)
	}

	// It is possible that server did not come up but we have tried enough times.
}

func httpGet(t *testing.T, url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ensureValidMetrics(t *testing.T, totalRequests int) {
	metricsOutput := httpGet(t, metricsEndpoint)

	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(bytes.NewReader(metricsOutput))
	if err != nil {
		t.Fatal(err)
	}

	expValue := int(*mf["runtime_http_requests_total_num"].Metric[0].Counter.Value)
	if expValue != totalRequests {
		t.Fatalf("expected [%v] for runtime_http_requests_total_num, got: %v", totalRequests, expValue)
	}
}

func TestRuntimeMetrics(t *testing.T) {
	setupRuntime(t)

	_ = httpGet(t, adminEndpoint)
	ensureValidMetrics(t, 1)

	_ = httpGet(t, graphqlEndpoint)
	ensureValidMetrics(t, 2)

	_ = httpGet(t, healthEndpoint)
	ensureValidMetrics(t, 2)
}
