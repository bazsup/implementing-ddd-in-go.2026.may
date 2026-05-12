package pricecalculation

var pricePerKg = map[string]Price{
	"Green waste":        NewPriceFromUSD(10),
	"Construction waste": NewPriceFromUSD(15),
}

type CalculatedPrice struct {
	PersonID      string
	VisitID       string
	PriceAmount   float64
	PriceCurrency string
}

func CalculatePrice(visit Visit) CalculatedPrice {
	var total Price
	for _, f := range visit.droppedFractions {
		total = total.Add(pricePerKg[f.FractionType].Times(f.AmountKg))
	}
	return CalculatedPrice{
		PersonID:      visit.personID,
		VisitID:       visit.visitID,
		PriceAmount:   total.Amount(),
		PriceCurrency: total.Currency(),
	}
}
