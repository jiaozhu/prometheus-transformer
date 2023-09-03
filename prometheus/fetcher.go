package prometheus

import (
	"fmt"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
	"io/ioutil"
	"log"
	"net/http"
	"prometheus-transformer/config"
	"strings"

	"prometheus-transformer/utils"
)

func FetchPrometheusData() (model.Vector, error) {
	federateURL := fmt.Sprintf("%s/federate?match[]={job=\"%s\"}", config.PrometheusFederateURL, config.JobName)
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
					transformedName := utils.CamelCaseToSnakeCase(label.GetValue())
					labels[model.LabelName("__name__")] = model.LabelValue(strings.ToLower(transformedName))
				} else {
					labels[model.LabelName(label.GetName())] = model.LabelValue(label.GetValue())
				}

			}

			if _, exists := labels[model.LabelName("__name__")]; !exists {
				transformedName := utils.CamelCaseToSnakeCase(m.GetName())
				labels[model.LabelName("__name__")] = model.LabelValue(strings.ToLower(transformedName))
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
