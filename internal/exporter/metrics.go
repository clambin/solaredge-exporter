package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	currentPower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "solaredge_current_power",
		Help: "Current Power",
	}, []string{"site"})
)
