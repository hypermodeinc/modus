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

var (
	// runtimePromRegistry is the registry for all the runtime metrics.
	runtimePromRegistry = prometheus.NewRegistry()

	// MetricsHandler is the handler for the /metrics endpoint.
	MetricsHandler = promhttp.HandlerFor(runtimePromRegistry, promhttp.HandlerOpts{})
)

var (
	// # of series = 1
	httpRequestsInFlightNum = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "runtime_http_requests_in_flight_num",
			Help: "A gauge of requests currently being served excluding /health & /metrics endpoints",
		},
	)
	// # of series = # of HTTP response codes
	httpRequestsTotalNum = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "runtime_http_requests_total_num",
			Help: "A counter for HTTP requests excluding /health & /metrics endpoints",
		},
		[]string{"code"},
	)
	// # of series = # of handlers x 5
	httpRequestsDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "runtime_http_requests_duration_seconds",
			Help:    "A histogram of latencies for requests excluding /health & /metrics endpoints",
			Buckets: []float64{.25, .5, 1, 5, 10},
		},
		[]string{"handler"},
	)
	// # of series = # of handlers x 4
	httpResponseSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "runtime_http_response_size_bytes",
			Help:    "A histogram of response sizes for requests excluding /health & /metrics endpoints",
			Buckets: []float64{10, 100, 1000, 10000},
		},
		[]string{"handler"},
	)

	// FunctionExecutionsNum is a counter for number of function executions done in this runtime instance.
	// # of series = 1
	FunctionExecutionsNum = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "runtime_function_executions_num",
			Help: "Number of function executions",
		},
	)
	// FunctionExecutionDurationMilliseconds is a histogram of latencies for wasm function executions of user plugins.
	// # of series = # of functions x 49
	FunctionExecutionDurationMilliseconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "runtime_function_execution_duration_milliseconds",
			Help: "A histogram of latencies for wasm function executions of user plugins",
			Buckets: []float64{
				10, 15, 20, 30, 40, 60, 80, 100, 125, 150, 175, 200, 225, 250, 275, 300,
				350, 400, 450, 500, 550, 600, 650, 700, 750, 800, 850, 900, 950, 1000,
				1100, 1200, 1300, 1400, 1500, 1600, 1700, 1800, 1900, 2000,
				2500, 3000, 3500, 4000, 5000, 10000, 20000, 40000, 60000,
			},
		},
		[]string{"function_name"},
	)

	// FunctionExecutionDurationMillisecondsSummary is a summary of latencies for wasm function executions of user plugins.
	FunctionExecutionDurationMillisecondsSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "runtime_function_execution_duration_milliseconds_summary",
			Help:       "A summary of latencies for wasm function executions of user plugins",
			Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.025, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"function_name"},
	)

	DroppedInferencesNum = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "runtime_dropped_inferences_num",
			Help: "Number of dropped inference requests",
		},
	)
)

func init() {
	runtimePromRegistry.MustRegister(
		httpRequestsInFlightNum,
		httpRequestsTotalNum,
		httpRequestsDurationSeconds,
		httpResponseSizeBytes,
		FunctionExecutionsNum,
		FunctionExecutionDurationMilliseconds,
		FunctionExecutionDurationMillisecondsSummary,
		DroppedInferencesNum,
	)
}

// InstrumentHandler wraps the provided http.Handler with metrics instrumentation.
func InstrumentHandler(handler http.Handler, handlerName string) http.Handler {
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
