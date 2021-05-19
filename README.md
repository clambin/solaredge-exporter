# solaredge-exporter
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/clambin/solaredge-exporter?color=green&label=Release&style=plastic)
![Codecov](https://img.shields.io/codecov/c/gh/clambin/solaredge-exporter?style=plastic)
![Build](https://github.com/clambin/solaredge-exporter/workflows/Build/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/clambin/solaredge-exporter)
![GitHub](https://img.shields.io/github/license/clambin/solaredge-exporter?style=plastic)

Basic Prometheus exporter for SolarEdge power inverters.

## Installation

Binaries are available on the [release](https://github.com/clambin/solaredge-exporter/releases) page. Docker images are available on [docker hub](https://hub.docker.com/r/clambin/solaredge-monitor).

Alternatively, you can clone the repository and build from source:

```
git clone https://github.com/clambin/solaredge-exporter.git
cd solaredge-exporter
go build
```

You will need to have Go 1.16 installed on your system.

## Running
### Command-line options

The following command-line arguments can be passed:

```
usage: solaredge-exporter --apikey=APIKEY [<flags>]

solaredge-exporter

Flags:
  -h, --help           Show context-sensitive help (also try --help-long and --help-man).
  -v, --version        Show application version.
  -d, --debug          Log debug messages
  -p, --port=8080      Prometheus listener port
  -i, --interval=15m   Measurement interval
  -a, --apikey=APIKEY  SolarEdge API key
```

where APIKEY is the API key you created on the Solaredge portal.

## Prometheus metrics

Currently, the exporter provides a single metric:

```
* solaredge_current_power: current output, in Watt
```

## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.