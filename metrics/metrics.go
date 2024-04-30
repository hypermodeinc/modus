/*
 * Copyright 2024 Hypermode, Inc.
 */

package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// This package defines all the runtime metrics that we expose on /metrics endpoint.
// Each metric name should be named as below:
//      runtime_{name}_{unit}
//  where unit is one of the following: num, seconds or bytes.
//  and name is the name of the metric.

const (
	// LabelFunction is the label for function name in a wasm module.
	LabelFunction = "function"
	// labelPlugin is the label for plugin name in a wasm module.
	LabelPlugin = "plugin"
)

var (
	// runtimePromRegistry is the registry for all the runtime metrics.
	runtimePromRegistry = prometheus.NewRegistry()

	// MetricsHandler is the handler for the /metrics endpoint.
	MetricsHandler = promhttp.HandlerFor(runtimePromRegistry, promhttp.HandlerOpts{})
)

var (
	httpRequestsInFlightNum = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "runtime_http_requests_in_flight_num",
			Help: "A gauge of requests currently being served excluding /health & /metrics endpoints",
		},
	)
	httpRequestsTotalNum = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "runtime_http_requests_total_num",
			Help: "A counter for HTTP requests excluding /health & /metrics endpoints",
		},
		[]string{"code"},
	)
	httpRequestsDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "runtime_http_requests_duration_seconds",
			Help:    "A histogram of latencies for requests excluding /health & /metrics endpoints",
			Buckets: []float64{.25, .5, 1, 5, 10},
		},
		[]string{"handler"},
	)
	httpResponseSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "runtime_http_response_size_bytes",
			Help:    "A histogram of response sizes for requests excluding /health & /metrics endpoints",
			Buckets: []float64{10, 100, 1000, 10000},
		},
		[]string{"handler"},
	)

	// FunctionExecutionsNum is a counter for number of function executions done in this runtime instance.
	FunctionExecutionsNum = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "runtime_function_executions_num",
			Help: "Number of function executions",
		},
		[]string{LabelPlugin, LabelFunction},
	)
)

func init() {
	runtimePromRegistry.MustRegister(
		httpRequestsInFlightNum,
		httpRequestsTotalNum,
		httpRequestsDurationSeconds,
		httpResponseSizeBytes,
		FunctionExecutionsNum,
	)
}

// InstrumentHandler wraps the provided http.Handler with metrics instrumentation.
func InstrumentHandler(handler http.HandlerFunc, handlerName string) http.Handler {
	return promhttp.InstrumentHandlerInFlight(
		httpRequestsInFlightNum,
		promhttp.InstrumentHandlerCounter(
			httpRequestsTotalNum,
			promhttp.InstrumentHandlerDuration(
				httpRequestsDurationSeconds.MustCurryWith(prometheus.Labels{"handler": handlerName}),
				promhttp.InstrumentHandlerResponseSize(
					httpResponseSizeBytes.MustCurryWith(prometheus.Labels{"handler": handlerName}),
					handler,
				),
			),
		),
	)
}
