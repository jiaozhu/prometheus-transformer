# Prometheus Transformer

## Overview

The Prometheus Transformer is a Go-based service designed to fetch metrics from a Prometheus instance, transform metric names to lowercase, and expose them via an HTTP endpoint. This is especially useful for services like `vmagent` to scrape.

## Features

- Fetch metrics from a specified Prometheus instance.
- Transform metric names (`__name__` label) to lowercase.
- Provide an HTTP endpoint to expose the transformed metrics.
- Filter metrics based on a specified job name.
- Detailed startup and operational logs for better monitoring and debugging.

## Prerequisites

- Go (version 1.16 or newer recommended)
- Access to a Prometheus instance

## Usage

### Building the Project

```bash
go build -o prometheus-transformer
```

### Running the Service

```bash
./prometheus-transformer -prometheus-url=http://your-prometheus-server -port=8080 -job-name=desired_job_name
```
- prometheus-url: The URL of the Prometheus federate endpoint.
- port: The port on which the transformer service will run. Default is 8080.
- job-name: The specific job name to filter metrics by.

Once the service is running, you can access the transformed metrics at: `http://localhost:{port}/metrics`

 ## Contributing

Feel free to open issues or pull requests if you have suggestions or improvements.