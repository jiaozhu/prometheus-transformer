# Prometheus Transformer

## Overview

The Prometheus Transformer is a Go-based service designed to fetch metrics from a Prometheus instance, transform metric names to lowercase, and expose them via an HTTP endpoint. This is especially useful for services like `vmagent` to scrape.

### Features

- **Fetch Metrics from Prometheus Federate Endpoint**: Easily retrieve metrics from a remote Prometheus federate endpoint.

- **Metric Name Transformation**: Converts the metric names in two stages:
  - Transforms camelCase portions after the last underscore into snake_case.
  - Converts the entire metric name into lowercase.

- **Filter by Job Name**: Provides an option to only fetch metrics from a specific job name.

- **Automated Builds with GitHub Actions**: Supports automated cross-compilation for x64 Linux and ARM64 architectures and attaches binaries to GitHub releases.


## Prerequisites

- Go (version 1.16 or newer recommended)
- Access to a Prometheus instance

## Usage

```bash
./prometheus-transformer -prometheus-url=http://your-prometheus-server -port=8191 -job-name=desired_job_name
```
- prometheus-url: The URL of the Prometheus federate endpoint.
- port: The port on which the transformer service will run. Default is 8191.
- job-name: The specific job name to filter metrics by.

Once the service is running, you can access the transformed metrics at: `http://localhost:8191/metrics`

## Contributing

Contributions are welcome! Feel free to open issues or pull requests if you have suggestions, improvements, or bug fixes. If you have any questions or need further clarification on any matter, you can reach out to the maintainer directly at [weijie@linux.com](mailto:weijie@linux.com).
