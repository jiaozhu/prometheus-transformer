package config

import (
	"flag"
)

var (
	PrometheusFederateURL string
	ServerPort            string
	JobName               string
)

func InitConfig() {
	flag.StringVar(&PrometheusFederateURL, "prometheus-url", "", "URL of the Prometheus federate endpoint")
	flag.StringVar(&ServerPort, "port", "8191", "Port on which the server will run")
	flag.StringVar(&JobName, "job-name", "", "Specific job name to filter metrics by")
	flag.Parse()
}
