package pricecalculation

var pricePerKg = map[string]float64{
	"Green waste":        0.10,
	"Construction waste": 0.15,
}

type CalculatedPrice struct {
	PersonID      string
	VisitID       string
	PriceAmount   float64
	PriceCurrency string
}

func CalculatePrice(visit Visit) CalculatedPrice {
	var total float64
	for _, f := range visit.droppedFractions {
		total += float64(f.AmountKg) * pricePerKg[f.FractionType]
	}
	return CalculatedPrice{
		PersonID:      visit.personID,
		VisitID:       visit.visitID,
		PriceAmount:   total,
		PriceCurrency: "USD",
	}
}
