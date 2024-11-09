package packaging

type Box struct{}

func (b *Box) Name() string {
	return "box"
}

func (b *Box) Cost() float32 {
	return 20
}

func (b *Box) MaxWeight() float64 {
	return 30
}
