package collector

import (
	"context"
	"fmt"
	"github.com/clambin/solaredge"
	solaredge2 "github.com/clambin/solaredge-exporter/solaredge"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"strconv"
)

type Collector struct {
	solaredge2.API
	currentPower *prometheus.Desc
	dayEnergy    *prometheus.Desc
	monthEnergy  *prometheus.Desc
	yearEnergy   *prometheus.Desc
}

func New(token string) *Collector {
	return &Collector{
		API: &solaredge.Client{
			Token: token,
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
	ch <- c.dayEnergy
	ch <- c.monthEnergy
	ch <- c.yearEnergy
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()
	sites, err := c.GetSites(ctx)

	if err == nil {
		for _, site := range sites {
			c.SetActiveSiteID(site.ID)
			var overview solaredge.PowerOverview
			if overview, err = c.GetPowerOverview(ctx); err != nil {
				break
			}
			ch <- prometheus.MustNewConstMetric(c.currentPower, prometheus.GaugeValue, overview.CurrentPower.Power, strconv.Itoa(site.ID))
			ch <- prometheus.MustNewConstMetric(c.dayEnergy, prometheus.GaugeValue, overview.LastDayData.Energy, strconv.Itoa(site.ID))
			ch <- prometheus.MustNewConstMetric(c.monthEnergy, prometheus.GaugeValue, overview.LastMonthData.Energy, strconv.Itoa(site.ID))
			ch <- prometheus.MustNewConstMetric(c.yearEnergy, prometheus.GaugeValue, overview.LastYearData.Energy, strconv.Itoa(site.ID))
		}
	}

	if err != nil {
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewDesc("solaredge_error",
				"Error while retrieving SolarEdge metrics", nil, nil),
			fmt.Errorf("error retrieving SolarEdge metrics: %w", err))
		slog.Error("failed to retrieve SolarEdge metrics", err)
	}
}
