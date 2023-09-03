package server

import (
	"fmt"
	"github.com/prometheus/common/model"
	"log"
	"net/http"
	"prometheus-transformer/config"
	"prometheus-transformer/prometheus"
)

func StartServer() {
	http.HandleFunc("/metrics", MetricsHandler)
	log.Printf("HTTP server started on port %s", config.ServerPort)
	err := http.ListenAndServe(":"+config.ServerPort, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics, err := prometheus.FetchPrometheusData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, metric := range metrics {
		metricName := metric.Metric[model.LabelName("__name__")]
		delete(metric.Metric, model.LabelName("__name__")) // Remove the __name__ label for proper formatting
		fmt.Fprintf(w, "%s%s %f %d\n", metricName, metric.Metric.String(), metric.Value, metric.Timestamp.Unix())
	}
}
