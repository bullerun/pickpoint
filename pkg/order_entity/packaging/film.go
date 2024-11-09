package packaging

import "math"

type Film struct{}

func (f *Film) Name() string {
	return "film"
}

func (f *Film) Cost() float32 {
	return 1
}

func (f *Film) MaxWeight() float64 {
	return math.MaxFloat64
}
