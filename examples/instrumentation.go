package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define Prometheus metrics
var (
	requestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hello_worlds_total",
			Help: "Hello Worlds requested.",
		},
	)
	exceptionCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hello_world_exceptions_total",
			Help: "Exceptions serving Hello World.",
		},
	)
	inProgress = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "hello_worlds_inprogress",
			Help: "Number of Hello Worlds in progress.",
		},
	)
	lastServedTime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "hello_world_last_time_seconds",
			Help: "The last time a Hello World was served.",
		},
	)
	latencySummary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "hello_world_latency_seconds",
			Help: "Time for a request Hello World.",
		},
	)
	latencyHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "hello_world_latency_histogram_seconds",
			Help:    "Histogram of request duration for Hello World.",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(requestCounter, exceptionCounter, inProgress, lastServedTime, latencySummary, latencyHistogram)
	rand.Seed(time.Now().UnixNano())
}

// MyHandler handles HTTP requests and responds with "Hello World".
func MyHandler(w http.ResponseWriter, r *http.Request) {
	inProgress.Inc()                               // Track in-progress requests
	timer := prometheus.NewTimer(latencyHistogram) // Start histogram timer
	defer timer.ObserveDuration()
	startTime := time.Now()
	requestCounter.Inc()

	// Simulate a 20% chance of failure
	if rand.Float64() < 0.2 {
		exceptionCounter.Inc()
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		inProgress.Dec()
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
	lastServedTime.Set(float64(time.Now().Unix()))          // Set last served time
	latencySummary.Observe(time.Since(startTime).Seconds()) // Observe request duration
	inProgress.Dec()                                        // Decrement in-progress counter
}

func main() {
	// Create a ServeMux for the Prometheus metrics server
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())

	// Start Prometheus metrics server on port 8000
	go func() {
		log.Println("Starting Prometheus metrics server on http://localhost:8000")
		log.Fatal(http.ListenAndServe("localhost:8000", metricsMux))
	}()

	// Create a ServeMux for the main HTTP server
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/", MyHandler)

	// Start HTTP server on port 8001
	log.Println("Starting main HTTP server on http://localhost:8001")
	log.Fatal(http.ListenAndServe("localhost:8001", mainMux))
}
