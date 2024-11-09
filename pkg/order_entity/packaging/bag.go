package packaging

type Bag struct{}

func (b *Bag) Name() string {
	return "bag"
}

func (b *Bag) Cost() float32 {
	return 5
}

func (b *Bag) MaxWeight() float64 {
	return 10
}
