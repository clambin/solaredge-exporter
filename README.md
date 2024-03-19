# solaredge-exporter
[![release](https://img.shields.io/github/v/tag/clambin/solaredge-exporter?color=green&label=release&style=plastic)](https://github.com/clambin/solaredge-exporter/releases)
[![codecov](https://img.shields.io/codecov/c/gh/clambin/solaredge-exporter?style=plastic)](https://app.codecov.io/gh/clambin/solaredge-exporter)
[![build](https://github.com/clambin/solaredge-exporter/actions/workflows/build.yaml/badge.svg)](https://github.com/clambin/solaredge-exporter/actions)
[![go report card](https://goreportcard.com/badge/github.com/clambin/solaredge-exporter)](https://goreportcard.com/report/github.com/clambin/solaredge-exporter)
[![license](https://img.shields.io/github/license/clambin/solaredge-exporter?style=plastic)](LICENSE.md)

Basic Prometheus exporter for SolarEdge power inverters.

## Installation

Binaries are available on the [release](https://github.com/clambin/solaredge-exporter/releases) page. Docker images are available on [ghcr.io](https://github.com/clambin/solaredge-exporter/pkgs/container/solaredge-exporter).

## Running
### Command-line options

The following command-line arguments can be passed:

```
Usage:
  solaredge-exporter [flags]

Flags:
      --addr string     Listener address for Prometheus metrics (default ":9090")
      --apikey string   SolarEdge API key
      --debug           Log debug messages
  -h, --help            help for solaredge-exporter
  -v, --version         version for solaredge-exporter
```

where apikey is your SolarEdge API key.

## Prometheus metrics

| metric | type |  labels | help |
| --- | --- |  --- | --- |
| solaredge_current_power | GAUGE | site|Current Power in Watt |
| solaredge_day_energy | GAUGE | site|Today's produced energy in WattHours |
| solaredge_exporter_http_request_duration_seconds | SUMMARY | code, method, path|http request duration in seconds |
| solaredge_exporter_http_requests_total | COUNTER | code, method, path|total number of http requests |
| solaredge_month_energy | GAUGE | site|This month's produced energy in WattHours |
| solaredge_year_energy | GAUGE | site|This year's produced energy in WattHours |

## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
