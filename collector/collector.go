package collector

import (
	"context"
	"fmt"
	"github.com/clambin/solaredge"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"strconv"
	"time"
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
	temperature = prometheus.NewDesc(
		prometheus.BuildFQName("solaredge", "inverter", "temperature"),
		"Temperature reported by the inverter(s)",
		[]string{"site"},
		nil,
	)
	acVoltage = prometheus.NewDesc(
		prometheus.BuildFQName("solaredge", "inverter", "ac_voltage"),
		"AC voltage reported by the inverter(s)",
		[]string{"site"},
		nil,
	)
	acCurrent = prometheus.NewDesc(
		prometheus.BuildFQName("solaredge", "inverter", "ac_current"),
		"AC current reported by the inverter(s)",
		[]string{"site"},
		nil,
	)
	dcVoltage = prometheus.NewDesc(
		prometheus.BuildFQName("solaredge", "inverter", "dc_voltage"),
		"DC voltage reported by the inverter(s)",
		[]string{"site"},
		nil,
	)
	powerLimit = prometheus.NewDesc(
		prometheus.BuildFQName("solaredge", "inverter", "power_limit"),
		"Power limit reported by the inverter(s)",
		[]string{"site"},
		nil,
	)
)

type Collector struct {
	Sites     []Site
	Inverters map[int][]Inverter
}

//go:generate mockery --name Site
type Site interface {
	GetID() int
	GetPowerOverview(ctx context.Context) (solaredge.PowerOverview, error)
	GetInverters(ctx context.Context) ([]solaredge.Inverter, error)
}

//go:generate mockery --name Inverter
type Inverter interface {
	GetTelemetry(ctx context.Context, start time.Time, end time.Time) ([]solaredge.InverterTelemetry, error)
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- currentPower
	ch <- dayEnergy
	ch <- monthEnergy
	ch <- yearEnergy
	ch <- temperature
	ch <- acVoltage
	ch <- acCurrent
	ch <- dcVoltage
	ch <- powerLimit
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()
	for _, site := range c.Sites {
		c.collectPowerOverview(ctx, site, ch)

		if c.Inverters != nil {
			if inverters, ok := c.Inverters[site.GetID()]; ok {
				for _, inverter := range inverters {
					c.collectInverterTelemetry(ctx, site.GetID(), inverter, ch)
				}
			}
		}
	}
}

func (c *Collector) collectPowerOverview(ctx context.Context, site Site, ch chan<- prometheus.Metric) {
	overview, err := site.GetPowerOverview(ctx)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewDesc(
				"solaredge_error",
				"Error while retrieving SolarEdge metrics",
				nil,
				nil),
			fmt.Errorf("solaredge: %w", err))
		slog.Error("failed to retrieve SolarEdge metrics", "err", err)
		return
	}

	id := strconv.Itoa(site.GetID())
	ch <- prometheus.MustNewConstMetric(currentPower, prometheus.GaugeValue, overview.CurrentPower.Power, id)
	ch <- prometheus.MustNewConstMetric(dayEnergy, prometheus.GaugeValue, overview.LastDayData.Energy, id)
	ch <- prometheus.MustNewConstMetric(monthEnergy, prometheus.GaugeValue, overview.LastMonthData.Energy, id)
	ch <- prometheus.MustNewConstMetric(yearEnergy, prometheus.GaugeValue, overview.LastYearData.Energy, id)

}

func (c *Collector) collectInverterTelemetry(ctx context.Context, siteID int, inverter Inverter, ch chan<- prometheus.Metric) {
	end := time.Now()
	start := end.Add(-10 * time.Minute)
	telemetry, err := inverter.GetTelemetry(ctx, start, end)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewDesc(
				"solaredge_error",
				"Error while retrieving SolarEdge metrics",
				nil,
				nil),
			fmt.Errorf("solaredge: %w", err),
		)
		slog.Error("failed to retrieve SolarEdge metrics", "err", err)
		return
	}

	if len(telemetry) == 0 {
		return
	}

	siteIDString := strconv.Itoa(siteID)
	measurement := telemetry[len(telemetry)-1]

	slog.Debug("telemetry received", "count", len(telemetry), "newest", time.Since(time.Time(measurement.Time)))

	ch <- prometheus.MustNewConstMetric(temperature, prometheus.GaugeValue, measurement.Temperature, siteIDString)
	ch <- prometheus.MustNewConstMetric(acVoltage, prometheus.GaugeValue, measurement.L1Data.AcVoltage, siteIDString)
	ch <- prometheus.MustNewConstMetric(acCurrent, prometheus.GaugeValue, measurement.L1Data.AcCurrent, siteIDString)
	ch <- prometheus.MustNewConstMetric(dcVoltage, prometheus.GaugeValue, measurement.DcVoltage, siteIDString)
	ch <- prometheus.MustNewConstMetric(powerLimit, prometheus.GaugeValue, measurement.PowerLimit, siteIDString)
}
