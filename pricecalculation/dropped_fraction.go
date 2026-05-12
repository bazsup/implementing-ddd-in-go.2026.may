package pricecalculation

var pricePerKg = map[string]Price{
	"Green waste":        NewPriceFromUSD(10),
	"Construction waste": NewPriceFromUSD(15),
}

type DroppedFraction struct {
	AmountKg     uint
	FractionType string
}

func (df DroppedFraction) CalculatePrice() Price {
	return pricePerKg[df.FractionType].Times(df.AmountKg)
}