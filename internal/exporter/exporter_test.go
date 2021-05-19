package exporter_test

import (
	"github.com/clambin/gotools/metrics"
	"github.com/clambin/solaredge-exporter/internal/exporter"
	"github.com/clambin/solaredge-exporter/pkg/solaredge"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	go func() {
		err := exporter.Run(&APIMock{}, 15*time.Minute)
		assert.NoError(t, err)
	}()

	assert.Eventually(t, func() bool {
		value, err := metrics.LoadValue("solaredge_current_power", "1")
		return err == nil && value == 1.0
	}, 5000*time.Millisecond, 10*time.Millisecond)
}

type APIMock struct{}

func (api *APIMock) GetSiteIDs() (siteIDs []int, err error) {
	return []int{1}, nil
}

func (api *APIMock) GetPower(_ int, startTime, endTime time.Time) (entries []solaredge.PowerMeasurement, err error) {
	var value float64

	for startTime.Before(endTime) {
		entries = append(entries, solaredge.PowerMeasurement{
			Time:  startTime,
			Value: value,
		})
		startTime = startTime.Add(15 * time.Minute)
		value += 100.0
	}
	return
}

func (api *APIMock) GetPowerOverview(_ int) (lifeTime, lastYear, lastMonth, lastDay, current float64, err error) {
	return 10000.0, 1000.0, 100.0, 10.0, 1.0, nil
}
