package collector

import (
	"bytes"
	"errors"
	"github.com/clambin/solaredge"
	"github.com/clambin/solaredge-exporter/collector/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	s := mocks.NewSite(t)
	s.
		On("GetPowerOverview", mock.AnythingOfType("*context.emptyCtx")).
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
	s.On("GetID").Return(1)

	c := Collector{Sites: []Site{s}}

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
	s.
		On("GetPowerOverview", mock.AnythingOfType("*context.emptyCtx")).
		Return(solaredge.PowerOverview{}, errors.New("error"))
	c := Collector{Sites: []Site{s}}

	r := prometheus.NewPedanticRegistry()
	r.MustRegister(&c)
	_, err := r.Gather()
	assert.Error(t, err)
}
