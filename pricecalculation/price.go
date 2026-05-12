package pricecalculation

import "errors"

var InvalidPriceAmount = errors.New("invalid amount for price")

type Price struct {
	amount float64
}

func NewPriceFromUSD(amount float64) (Price, error) {
	if amount < 0 {
		return Price{}, InvalidPriceAmount
	}
	return Price{amount: amount}, nil
}

func (p Price) Amount() float64 {
	return p.amount
}

func (p Price) Add(other Price) (Price, error) {
	return NewPriceFromUSD(p.amount + other.amount)
}

func (p Price) Equals(other Price) bool {
	return p.amount == other.amount
}
