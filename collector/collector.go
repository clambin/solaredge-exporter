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
	dayEnergy    *prometheus.Desc
	monthEnergy  *prometheus.Desc
	yearEnergy   *prometheus.Desc
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
		),
		dayEnergy: prometheus.NewDesc(
			prometheus.BuildFQName("solaredge", "", "day_energy"),
			"Today's produced energy in WattHours",
			[]string{"site"},
			nil,
		),
		monthEnergy: prometheus.NewDesc(
			prometheus.BuildFQName("solaredge", "", "month_energy"),
			"This month's produced energy in WattHours",
			[]string{"site"},
			nil,
		),
		yearEnergy: prometheus.NewDesc(
			prometheus.BuildFQName("solaredge", "", "year_energy"),
			"This year's produced energy in WattHours",
			[]string{"site"},
			nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.currentPower
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()
	sites, err := c.GetSiteIDs(ctx)

	if err == nil {
		var year, month, day, current float64
		for _, site := range sites {
			_, year, month, day, current, err = c.GetPowerOverview(ctx, site)
			if err != nil {
				break
			}
			ch <- prometheus.MustNewConstMetric(c.currentPower, prometheus.GaugeValue, current, strconv.Itoa(site))
			ch <- prometheus.MustNewConstMetric(c.dayEnergy, prometheus.GaugeValue, day, strconv.Itoa(site))
			ch <- prometheus.MustNewConstMetric(c.monthEnergy, prometheus.GaugeValue, month, strconv.Itoa(site))
			ch <- prometheus.MustNewConstMetric(c.yearEnergy, prometheus.GaugeValue, year, strconv.Itoa(site))
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
