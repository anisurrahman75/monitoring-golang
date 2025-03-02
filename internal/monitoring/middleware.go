package monitoring

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func summary(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()
		httpRequestSummary.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}

func httpRequestCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		method := r.Method
		path := r.URL.Path
		status := strconv.Itoa(recorder.statusCode)
		next.ServeHTTP(recorder, r)
		httpRequestCounter.WithLabelValues(status, path, method).Inc()
	})
}

func activeRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activeRequestsGauge.Inc()
		defer activeRequestsGauge.Dec()
		delay := time.Duration(rand.Intn(900)) * time.Millisecond
		time.Sleep(delay)
		next.ServeHTTP(w, r)
	})
}
