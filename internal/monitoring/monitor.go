package monitoring

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	router *chi.Mux
}

func NewMonitoringServer() *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(httpRequestCount)
	r.Use(activeRequests)
	r.Use(summary)

	r.HandleFunc("/", helloWorld)
	r.HandleFunc("/db-query", dbQuery)
	r.Handle("/metrics", promhttp.HandlerFor(newReg, promhttp.HandlerOpts{}))

	return &Server{
		router: r,
	}
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.router)
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello world!"))
}

func dbQuery(w http.ResponseWriter, r *http.Request) {
	executeDBQuery := func(queryType string) {
		start := time.Now()
		delay := time.Duration(rand.Intn(900)) * time.Millisecond
		time.Sleep(delay)
		dbQueryHistogram.WithLabelValues(queryType).Observe(time.Since(start).Seconds())
	}

	queryTypes := []string{"SELECT"}
	randomQuery := queryTypes[rand.Intn(len(queryTypes))]
	executeDBQuery(randomQuery)
	w.Write([]byte(fmt.Sprintf("Executed %s query\n", randomQuery)))
}
