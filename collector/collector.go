package collector

import (
	"context"
	"fmt"
	"github.com/clambin/solaredge"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"strconv"
)

var (
	currentPower = prometheus.NewDesc(
		prometheus.BuildFQName("solaredge", "", "current_power"),
		"Current Power in Watt",
		[]string{"site"},
		nil,
	)
	dayEnergy = prometheus.NewDesc(
		prometheus.BuildFQName("solaredge", "", "day_energy"),
		"Today's produced energy in WattHours",
		[]string{"site"},
		nil,
	)
	monthEnergy = prometheus.NewDesc(
		prometheus.BuildFQName("solaredge", "", "month_energy"),
		"This month's produced energy in WattHours",
		[]string{"site"},
		nil,
	)
	yearEnergy = prometheus.NewDesc(
		prometheus.BuildFQName("solaredge", "", "year_energy"),
		"This year's produced energy in WattHours",
		[]string{"site"},
		nil,
	)
)

type Collector struct {
	Sites []Site
}

//go:generate mockery --name Site
type Site interface {
	GetID() int
	GetPowerOverview(ctx context.Context) (solaredge.PowerOverview, error)
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- currentPower
	ch <- dayEnergy
	ch <- monthEnergy
	ch <- yearEnergy
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()
	for _, site := range c.Sites {
		overview, err := site.GetPowerOverview(ctx)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(
				prometheus.NewDesc("solaredge_error",
					"Error while retrieving SolarEdge metrics", nil, nil),
				fmt.Errorf("error retrieving SolarEdge metrics: %w", err))
			slog.Error("failed to retrieve SolarEdge metrics", err)
			return
		}

		id := strconv.Itoa(site.GetID())
		ch <- prometheus.MustNewConstMetric(currentPower, prometheus.GaugeValue, overview.CurrentPower.Power, id)
		ch <- prometheus.MustNewConstMetric(dayEnergy, prometheus.GaugeValue, overview.LastDayData.Energy, id)
		ch <- prometheus.MustNewConstMetric(monthEnergy, prometheus.GaugeValue, overview.LastMonthData.Energy, id)
		ch <- prometheus.MustNewConstMetric(yearEnergy, prometheus.GaugeValue, overview.LastYearData.Energy, id)
	}
}
