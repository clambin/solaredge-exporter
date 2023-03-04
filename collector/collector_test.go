package collector_test

import (
	"bytes"
	"fmt"
	"github.com/clambin/solaredge"
	"github.com/clambin/solaredge-exporter/collector"
	"github.com/clambin/solaredge-exporter/solaredge/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	c := collector.New("sometoken")
	mockAPI := mocks.NewAPI(t)
	c.API = mockAPI

	mockAPI.On("GetSites", mock.Anything).Return([]solaredge.Site{{ID: 123}}, nil)
	mockAPI.On("SetActiveSiteID", 123)
	mockAPI.On("GetPowerOverview", mock.Anything).Return(
		solaredge.PowerOverview{
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
			}(struct{ Power float64 }{Power: 3400}),
		}, nil)

	r := prometheus.NewPedanticRegistry()
	r.MustRegister(c)

	const expected = `# HELP solaredge_current_power Current Power in Watt
# TYPE solaredge_current_power gauge
solaredge_current_power{site="123"} 3400
# HELP solaredge_day_energy Today's produced energy in WattHours
# TYPE solaredge_day_energy gauge
solaredge_day_energy{site="123"} 10
# HELP solaredge_month_energy This month's produced energy in WattHours
# TYPE solaredge_month_energy gauge
solaredge_month_energy{site="123"} 100
# HELP solaredge_year_energy This year's produced energy in WattHours
# TYPE solaredge_year_energy gauge
solaredge_year_energy{site="123"} 1000
`
	assert.NoError(t, testutil.GatherAndCompare(r, bytes.NewBufferString(expected)))
}

func TestCollector_Collect_Failure(t *testing.T) {
	c := collector.New("sometoken")
	mockAPI := mocks.NewAPI(t)
	c.API = mockAPI

	mockAPI.On("GetSites", mock.Anything).Return(nil, fmt.Errorf("get failed: %w", &solaredge.HTTPError{
		StatusCode: http.StatusForbidden,
		Status:     "403 Forbidden",
	}))

	r := prometheus.NewPedanticRegistry()
	r.MustRegister(c)
	_, err := r.Gather()
	assert.Error(t, err)
}
