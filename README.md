# solaredge-exporter
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/clambin/solaredge-exporter?color=green&label=Release&style=plastic)
![Codecov](https://img.shields.io/codecov/c/gh/clambin/solaredge-exporter?style=plastic)
![Build](https://github.com/clambin/solaredge-exporter/workflows/Build/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/clambin/solaredge-exporter)
![GitHub](https://img.shields.io/github/license/clambin/solaredge-exporter?style=plastic)

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

Currently, the exporter provides the following metrics:

```
# HELP solaredge_current_power Current Power in Watt
## TYPE solaredge_current_power gauge
solaredge_current_power{site="1"} 3400

# HELP solaredge_day_energy Today's produced energy in WattHours
# TYPE solaredge_day_energy gauge
solaredge_day_energy{site="1"} 10

# HELP solaredge_inverter_ac_current AC current reported by the inverter(s)
# TYPE solaredge_inverter_ac_current gauge
solaredge_inverter_ac_current{site="1"} 15

# HELP solaredge_inverter_ac_voltage AC voltage reported by the inverter(s)
# TYPE solaredge_inverter_ac_voltage gauge
solaredge_inverter_ac_voltage{site="1"} 220

# HELP solaredge_inverter_dc_voltage DC voltage reported by the inverter(s)
# TYPE solaredge_inverter_dc_voltage gauge
solaredge_inverter_dc_voltage{site="1"} 300

# HELP solaredge_inverter_power_limit Power limit reported by the inverter(s)
# TYPE solaredge_inverter_power_limit gauge
solaredge_inverter_power_limit{site="1"} 100

# HELP solaredge_inverter_temperature Temperature reported by the inverter(s)
# TYPE solaredge_inverter_temperature gauge
solaredge_inverter_temperature{site="1"} 25

# HELP solaredge_month_energy This month's produced energy in WattHours
# TYPE solaredge_month_energy gauge
solaredge_month_energy{site="1"} 100

# HELP solaredge_year_energy This year's produced energy in WattHours
# TYPE solaredge_year_energy gauge
solaredge_year_energy{site="1"} 1000
```

## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.