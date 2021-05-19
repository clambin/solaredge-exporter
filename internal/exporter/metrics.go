package exporter

import (
	"github.com/clambin/gotools/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	currentPower = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "solaredge_current_power",
		Help: "Current Power",
	}, []string{"site"})
)
