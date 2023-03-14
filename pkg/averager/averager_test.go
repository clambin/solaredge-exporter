package averager_test

import (
	"github.com/clambin/solaredge-exporter/pkg/averager"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAverager_Integer(t *testing.T) {
	a := averager.Averager[int]{}

	for i := 1; i < 6; i++ {
		a.Add(i)
	}
	assert.Equal(t, 3, a.Average())
	assert.Zero(t, a.Average())
}

func TestAverager_Float(t *testing.T) {
	a := averager.Averager[float64]{}

	for i := 1; i < 6; i++ {
		a.Add(float64(i))
	}
	assert.Equal(t, 3.0, a.Average())
	assert.Zero(t, a.Average())
}
