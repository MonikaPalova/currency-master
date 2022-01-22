package wallet

type Purchase struct {
	Asset    Asset
	Quantity float64
	PriceUSD float64
}

func (p Purchase) Total() float64 {
	return p.Quantity * p.PriceUSD
}
