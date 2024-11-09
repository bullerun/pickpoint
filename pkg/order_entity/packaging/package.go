package packaging

type Type interface {
	Name() string
	Cost() float32
	MaxWeight() float64
}
