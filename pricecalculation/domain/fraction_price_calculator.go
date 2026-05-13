package domain

type FractionPriceCalculator interface {
	CalculatePrice(droppedFraction DroppedFraction) Price
}

type FlatRatePriceCalculator struct {
	pricePerKG Price
}

func NewFlatRatePriceCalculator(pricePerKG Price) FlatRatePriceCalculator {
	return FlatRatePriceCalculator{pricePerKG: pricePerKG}
}

func (c FlatRatePriceCalculator) CalculatePrice(df DroppedFraction) Price {
	return c.pricePerKG.Times(df.weight.Amount())
}
