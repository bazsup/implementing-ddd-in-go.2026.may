package pricecalculation

var pricePerKg = map[string]Price{
	"Green waste":        NewPriceFromUSD(10),
	"Construction waste": NewPriceFromUSD(15),
}

type DroppedFraction struct {
	weight       Weight
	FractionType string
}

func NewDroppedFraction(fractionType string, weight Weight) DroppedFraction {
	return DroppedFraction{FractionType: fractionType, weight: weight}
}

func (df DroppedFraction) CalculatePrice() Price {
	return pricePerKg[df.FractionType].Times(df.weight.Amount())
}
