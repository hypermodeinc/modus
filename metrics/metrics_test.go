/*
 * Copyright 2024 Hypermode, Inc.
 */

package metrics_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"hypruntime/httpserver"
	"hypruntime/metrics"

	"github.com/prometheus/common/expfmt"
)

const (
	graphqlEndpoint = "/graphql"
	adminEndpoint   = "/admin"
	healthEndpoint  = "/health"
	metricsEndpoint = "/metrics"
)

func httpGet(t *testing.T, s *httptest.Server, endpoint string) []byte {
	resp, err := http.Get(s.URL + endpoint)
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ensureValidMetrics(t *testing.T, s *httptest.Server, totalRequests int) {
	metricsOutput := httpGet(t, s, metricsEndpoint)

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

	mux := httpserver.GetHandlerMux()
	s := httptest.NewServer(mux)
	defer s.Close()

	_ = httpGet(t, s, adminEndpoint)
	ensureValidMetrics(t, s, 1)

	_ = httpGet(t, s, graphqlEndpoint)
	ensureValidMetrics(t, s, 2)

	_ = httpGet(t, s, healthEndpoint)
	ensureValidMetrics(t, s, 2)
}

func BenchmarkSummary(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.FunctionExecutionDurationMillisecondsSummary.WithLabelValues("test").Observe(float64(i))
	}
}

func BenchmarkHistogram(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.FunctionExecutionDurationMilliseconds.WithLabelValues("test").Observe(float64(i))
	}
}
