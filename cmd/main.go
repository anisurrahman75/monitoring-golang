package main

import (
	"fmt"
	"github.com/anisurrahman75/monitoring-golang/internal/monitoring"
	"log"
)

func main() {
	fmt.Println("Starting Monitoring Service...")
	server := monitoring.NewMonitoringServer()

	port := "8080"
	log.Printf("Metrics server running on port %s\n", port)
	log.Fatal(server.ListenAndServe(fmt.Sprintf(":%s", port)))
}
