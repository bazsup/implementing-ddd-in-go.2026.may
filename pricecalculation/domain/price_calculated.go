package domain

type PriceCalculated struct {
	PersonID      string
	VisitID       string
	CustomerType  string
	PriceAmount   float64
	PriceCurrency string
	Email         string
}

func (e PriceCalculated) Type() string { return "PriceCalculated" }