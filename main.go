package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"org.donghyuns.com/exporter/basic/metrics"
)

func main() {
	if loadErr := godotenv.Load(".env"); loadErr != nil {
		log.Fatalf("Error loading .env file: %v", loadErr)
	}

	convInt, convErr := strconv.Atoi(os.Getenv("METRICS_INTERVAL"))
	if convErr != nil {
		log.Fatalf("Error converting METRICS_INTERVAL to integer: %v", convErr)
	}

	interval := time.Duration(convInt)

	metrics.MetricsScheduler(interval)

	// /metrics endpoint 노출
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Exporter listening on :9468/metrics")
	log.Printf("Metrics update interval: %d seconds", interval)
	log.Fatal(http.ListenAndServe(":9468", nil))
}
