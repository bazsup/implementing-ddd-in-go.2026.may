package pricecalculation

import "errors"

var InvalidPriceAmount = errors.New("invalid amount for price")

type Price struct {
	amountInSmallestUnit uint
}

func NewPriceFromUSDcents(amount uint) Price {
	return Price{amountInSmallestUnit: amount}
}

func (p Price) Amount() float64 {
	return float64(p.amountInSmallestUnit) / 100
}

func (p Price) Currency() string {
	return "USD"
}

func (p Price) Add(other Price) Price {
	return Price{p.amountInSmallestUnit + other.amountInSmallestUnit}
}

func (p Price) Times(factor uint) Price {
	return Price{p.amountInSmallestUnit * factor}
}

func (p Price) Equals(other Price) bool {
	return p.amountInSmallestUnit == other.amountInSmallestUnit
}
