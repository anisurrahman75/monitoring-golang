package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

var httpRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests received",
}, []string{"status", "path", "method"})

var activeRequestsGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "http_active_requests",
		Help: "Number of active connections to the service",
	},
)

var dbQueryHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "db_query_duration_seconds",
		Help:    "Histogram of database query latencies in seconds.",
		Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10},
	},
	[]string{"query_type"},
)

var httpRequestSummary = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "http_request_duration_summary_seconds",
		Help:       "Summary of HTTP request durations in seconds.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method", "path"},
)

var newReg = prometheus.NewRegistry()

func init() {
	newReg.MustRegister(httpRequestCounter, activeRequestsGauge, dbQueryHistogram, httpRequestSummary)
}
