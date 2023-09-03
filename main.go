package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
)

var (
	prometheusFederateURL string
	serverPort            string
	jobName               string
	startTime             = time.Now()
	requestCount          uint64
)

func init() {
	flag.StringVar(&prometheusFederateURL, "prometheus-url", "", "URL of the Prometheus federate endpoint")
	flag.StringVar(&serverPort, "port", "8191", "Port on which the server will run")
	flag.StringVar(&jobName, "job-name", "", "Name of the job to query metrics for")
	flag.Parse()
}

func fetchPrometheusData() (model.Vector, error) {
	federateURL := fmt.Sprintf("%s/federate?match[]={job=\"%s\"}", prometheusFederateURL, jobName)
	log.Printf("Fetching data from: %s", federateURL)
	resp, err := http.Get(federateURL)
	if err != nil {
		log.Printf("Error fetching data from Prometheus: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	var parser expfmt.TextParser
	metricFamilies, err := parser.TextToMetricFamilies(strings.NewReader(string(body)))
	if err != nil {
		log.Printf("Error parsing Prometheus data: %v", err)
		return nil, err
	}

	var metrics model.Vector
	for _, m := range metricFamilies {
		for _, metric := range m.Metric {
			labels := make(model.Metric)
			for _, label := range metric.Label {
				if label.GetName() == "__name__" {
					labels[model.LabelName("__name__")] = model.LabelValue(strings.ToLower(label.GetValue()))
				} else {
					labels[model.LabelName(label.GetName())] = model.LabelValue(label.GetValue())
				}
			}
			if _, exists := labels[model.LabelName("__name__")]; !exists {
				labels[model.LabelName("__name__")] = model.LabelValue(strings.ToLower(m.GetName()))
			}

			var value float64
			switch *m.Type {
			case dto.MetricType_COUNTER:
				if counter := metric.GetCounter(); counter != nil {
					value = counter.GetValue()
				}
			case dto.MetricType_GAUGE:
				if gauge := metric.GetGauge(); gauge != nil {
					value = gauge.GetValue()
				}
			case dto.MetricType_SUMMARY:
				if summary := metric.GetSummary(); summary != nil {
					value = summary.GetSampleSum()
				}
			case dto.MetricType_HISTOGRAM:
				if histogram := metric.GetHistogram(); histogram != nil {
					value = histogram.GetSampleSum()
				}
			case dto.MetricType_UNTYPED:
				if untyped := metric.GetUntyped(); untyped != nil {
					value = untyped.GetValue()
				}
			}

			metrics = append(metrics, &model.Sample{
				Metric:    labels,
				Value:     model.SampleValue(value),
				Timestamp: model.TimeFromUnixNano(metric.GetTimestampMs() * 1e6),
			})
		}
	}

	return metrics, nil
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&requestCount, 1)

	uptime := time.Since(startTime).Seconds()
	fmt.Fprintf(w, "# HELP transformer_uptime_seconds Total time the transformer has been up in seconds.\n")
	fmt.Fprintf(w, "# TYPE transformer_uptime_seconds gauge\n")
	fmt.Fprintf(w, "transformer_uptime_seconds %f\n", uptime)

	fmt.Fprintf(w, "# HELP transformer_requests_total Total number of requests processed by the transformer.\n")
	fmt.Fprintf(w, "# TYPE transformer_requests_total counter\n")
	fmt.Fprintf(w, "transformer_requests_total %d\n", requestCount)

	metrics, err := fetchPrometheusData()
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

func main() {
	log.Println("Starting Prometheus Transformer...")
	log.Printf("Target Prometheus URL: %s", prometheusFederateURL)
	log.Printf("Listening on port: %s", serverPort)
	log.Printf("Filtering by job name: %s", jobName)
	http.HandleFunc("/metrics", metricsHandler)
	err := http.ListenAndServe(":"+serverPort, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
