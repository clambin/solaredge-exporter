package averager

import "golang.org/x/exp/constraints"

type Value interface {
	constraints.Integer | constraints.Float
}

type Averager[T Value] struct {
	total T
	count int
}

func (a *Averager[T]) Add(value T) {
	a.total += value
	a.count++
}

func (a *Averager[T]) Average() T {
	var avg T
	if a.count > 0 {
		avg = a.total / T(a.count)
	}
	a.total = 0
	a.count = 0
	return avg
}
