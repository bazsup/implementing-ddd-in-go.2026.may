package pricecalculation

import "implementing-ddd-in-go/pricecalculation/domain"

type CalculatedPrice struct {
	PersonID      string
	VisitID       string
	PriceAmount   float64
	PriceCurrency string
}

func CalculatePrice(visit Visit) CalculatedPrice {
	var total domain.Price
	for _, df := range visit.droppedFractions {
		total = total.Add(df.CalculatePrice())
	}
	return CalculatedPrice{
		PersonID:      visit.personID,
		VisitID:       visit.visitID,
		PriceAmount:   total.Amount(),
		PriceCurrency: total.Currency(),
	}
}
