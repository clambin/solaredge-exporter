package collector

import (
	"context"
	"fmt"
	"github.com/clambin/solaredge"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Collector struct {
	solaredge.API
	currentPower *prometheus.Desc
}

func New(token string) *Collector {
	return &Collector{
		API: &solaredge.Client{
			Token:      token,
			HTTPClient: &http.Client{},
			APIURL:     "",
		},
		currentPower: prometheus.NewDesc(
			prometheus.BuildFQName("solaredge", "", "current_power"),
			"Current Power in Watt",
			[]string{"site"},
			nil,
		)}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.currentPower
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()
	sites, err := collector.GetSiteIDs(ctx)

	if err == nil {
		for _, site := range sites {
			var current float64
			_, _, _, _, current, err = collector.GetPowerOverview(ctx, site)
			if err == nil {
				ch <- prometheus.MustNewConstMetric(collector.currentPower, prometheus.GaugeValue, current, strconv.Itoa(site))
			}
		}
	}

	if err != nil {
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewDesc("solaredge_error",
				"Error while retrieving SolarEdge metrics", nil, nil),
			fmt.Errorf("error retrieving SolarEdge metrics: %w", err))
		log.WithError(err).Warning("failed to retrieve SolarEdge metrics")
	}
}
