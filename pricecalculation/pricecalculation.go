package pricecalculation

type CalculatedPrice struct {
	PersonID      string
	VisitID       string
	PriceAmount   float64
	PriceCurrency string
}

func CalculatePrice(visit Visit) CalculatedPrice {
	return CalculatedPrice{
		PersonID:      visit.personID,
		VisitID:       visit.visitID,
		PriceAmount:   1,
		PriceCurrency: "USD",
	}
}
