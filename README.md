# Prometheus Transformer

## Overview

The Prometheus Transformer is a simple service written in Go that fetches metrics from a Prometheus federate endpoint, transforms the metric names to lowercase, and then provides an HTTP endpoint for services like `vmagent` to scrape.

## Features

- Fetch metrics from a remote Prometheus federate endpoint.
- Transform metric names (`__name__` label) to lowercase.
- Provide an HTTP endpoint to expose the transformed metrics.

## Prerequisites

- Go (version 1.16 or newer recommended)
- Access to a Prometheus federate endpoint

## Usage

### Building the Project

```bash
go build -o prometheus-transformer
```

### Running the Service

```bash
./prometheus-transformer -prometheus-url=http://your-prometheus-server/federate -port=8080
```
- prometheus-url: The URL of the Prometheus federate endpoint.
- port: The port on which the transformer service will run. Default is 8080.

Once the service is running, you can access the transformed metrics at: `http://localhost:8080/metrics`

 ## Contributing

Feel free to open issues or pull requests if you have suggestions or improvements.