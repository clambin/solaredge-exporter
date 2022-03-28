package collector_test

import (
	"fmt"
	"github.com/clambin/solaredge"
	"github.com/clambin/solaredge-exporter/collector"
	"github.com/clambin/solaredge/mocks"
	"github.com/prometheus/client_golang/prometheus"
	pcg "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCollector_Describe(t *testing.T) {
	c := collector.New("sometoken")
	c.API = &mocks.API{}

	ch := make(chan *prometheus.Desc)
	go c.Describe(ch)

	m := <-ch

	assert.Contains(t, m.String(), "Desc{fqName: \"solaredge_current_power\"")
}

func TestCollector_Collect(t *testing.T) {
	c := collector.New("sometoken")
	mockAPI := &mocks.API{}
	c.API = mockAPI

	mockAPI.On("GetSiteIDs", mock.Anything).Return([]int{123}, nil)
	mockAPI.On("GetPowerOverview", mock.Anything, 123).Return(0.0, 1000.0, 100.0, 10.0, 3400.0, nil)

	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	for _, expected := range []float64{3400, 10, 100, 1000} {
		m := <-ch
		readMetric := pcg.Metric{}
		err := m.Write(&readMetric)
		require.NoError(t, err)
		assert.Equal(t, expected, readMetric.GetGauge().GetValue())
	}

	mockAPI.AssertExpectations(t)
}

func TestCollector_Collect_Failure(t *testing.T) {
	c := collector.New("sometoken")
	mockAPI := &mocks.API{}
	c.API = mockAPI

	mockAPI.On("GetSiteIDs", mock.Anything).Return(nil, fmt.Errorf("get failed: %w", &solaredge.HTTPError{
		StatusCode: http.StatusForbidden,
		Status:     "403 Forbidden",
	}))

	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	m := <-ch
	readMetric := pcg.Metric{}
	err := m.Write(&readMetric)
	require.Error(t, err)
	assert.Equal(t, "error retrieving SolarEdge metrics: get failed: 403 Forbidden", err.Error())
	mockAPI.AssertExpectations(t)
}
