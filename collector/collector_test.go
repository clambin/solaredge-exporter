package collector_test

import (
	"bytes"
	"errors"
	"github.com/clambin/solaredge"
	"github.com/clambin/solaredge-exporter/collector"
	"github.com/clambin/solaredge-exporter/collector/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	s := mocks.NewSite(t)
	s.EXPECT().GetPowerOverview(mock.Anything).
		Return(solaredge.PowerOverview{
			LastUpdateTime: solaredge.Time{},
			LifeTimeData: struct {
				Energy  float64 `json:"energy"`
				Revenue float64 `json:"revenue"`
			}{
				Energy: 10000,
			},
			LastYearData: struct {
				Energy  float64 `json:"energy"`
				Revenue float64 `json:"revenue"`
			}{
				Energy: 1000,
			},
			LastMonthData: struct {
				Energy  float64 `json:"energy"`
				Revenue float64 `json:"revenue"`
			}{
				Energy: 100,
			},
			LastDayData: struct {
				Energy  float64 `json:"energy"`
				Revenue float64 `json:"revenue"`
			}{
				Energy: 10,
			},
			CurrentPower: struct {
				Power float64 `json:"power"`
			}{
				Power: 3400,
			},
		}, nil)
	s.EXPECT().GetID().Return(1)

	i := mocks.NewInverter(t)
	i.EXPECT().GetTelemetry(mock.Anything, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
		Return([]solaredge.InverterTelemetry{
			{
				L1Data: struct {
					AcCurrent     float64 `json:"acCurrent"`
					AcFrequency   float64 `json:"acFrequency"`
					AcVoltage     float64 `json:"acVoltage"`
					ActivePower   float64 `json:"activePower"`
					ApparentPower float64 `json:"apparentPower"`
					CosPhi        float64 `json:"cosPhi"`
					ReactivePower float64 `json:"reactivePower"`
				}{AcCurrent: 15, AcFrequency: 50, AcVoltage: 220},
				Temperature: 25,
				PowerLimit:  100,
				DcVoltage:   300,
			},
		}, nil)

	c := collector.Collector{
		Sites: []collector.Site{s},
		Inverters: map[int][]collector.Inverter{
			1: {i},
		},
	}

	r := prometheus.NewPedanticRegistry()
	r.MustRegister(&c)

	const expected = `# HELP solaredge_current_power Current Power in Watt
# TYPE solaredge_current_power gauge
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
`
	assert.NoError(t, testutil.GatherAndCompare(r, bytes.NewBufferString(expected)))
}

func TestCollector_Collect_NoTelemetry(t *testing.T) {
	s := mocks.NewSite(t)
	s.EXPECT().GetPowerOverview(mock.Anything).
		Return(solaredge.PowerOverview{
			LastUpdateTime: solaredge.Time{},
			LifeTimeData: struct {
				Energy  float64 `json:"energy"`
				Revenue float64 `json:"revenue"`
			}{
				Energy: 10000,
			},
			LastYearData: struct {
				Energy  float64 `json:"energy"`
				Revenue float64 `json:"revenue"`
			}{
				Energy: 1000,
			},
			LastMonthData: struct {
				Energy  float64 `json:"energy"`
				Revenue float64 `json:"revenue"`
			}{
				Energy: 100,
			},
			LastDayData: struct {
				Energy  float64 `json:"energy"`
				Revenue float64 `json:"revenue"`
			}{
				Energy: 10,
			},
			CurrentPower: struct {
				Power float64 `json:"power"`
			}{
				Power: 3400,
			},
		}, nil)
	s.EXPECT().GetID().Return(1)

	i := mocks.NewInverter(t)
	i.EXPECT().GetTelemetry(mock.Anything, mock.Anything, mock.AnythingOfType("time.Time")).
		Return([]solaredge.InverterTelemetry{}, nil)

	c := collector.Collector{
		Sites: []collector.Site{s},
		Inverters: map[int][]collector.Inverter{
			1: {i},
		},
	}

	r := prometheus.NewPedanticRegistry()
	r.MustRegister(&c)

	const expected = `# HELP solaredge_current_power Current Power in Watt
# TYPE solaredge_current_power gauge
solaredge_current_power{site="1"} 3400
# HELP solaredge_day_energy Today's produced energy in WattHours
# TYPE solaredge_day_energy gauge
solaredge_day_energy{site="1"} 10
# HELP solaredge_month_energy This month's produced energy in WattHours
# TYPE solaredge_month_energy gauge
solaredge_month_energy{site="1"} 100
# HELP solaredge_year_energy This year's produced energy in WattHours
# TYPE solaredge_year_energy gauge
solaredge_year_energy{site="1"} 1000
`
	assert.NoError(t, testutil.GatherAndCompare(r, bytes.NewBufferString(expected)))
}

func TestCollector_Collect_Failure(t *testing.T) {
	s := mocks.NewSite(t)
	i := mocks.NewInverter(t)
	c := collector.Collector{
		Sites:     []collector.Site{s},
		Inverters: map[int][]collector.Inverter{1: {i}},
	}
	r := prometheus.NewPedanticRegistry()
	r.MustRegister(&c)

	s.EXPECT().GetPowerOverview(mock.Anything).Return(solaredge.PowerOverview{}, errors.New("error"))
	s.EXPECT().GetID().Return(1)
	i.EXPECT().GetTelemetry(mock.Anything, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
		Return(nil, errors.New("error"))

	_, err := r.Gather()
	assert.Error(t, err)
}
