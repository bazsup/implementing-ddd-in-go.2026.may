package pricecalculation

import "errors"

type Price struct {
	amount   float64
	currency Currency
}

func NewPrice(amount float64, currency Currency) Price {
	return Price{amount: amount, currency: currency}
}

func (p Price) Amount() float64 {
	return p.amount
}

func (p Price) Currency() Currency {
	return p.currency
}

func (p Price) Add(other Price) (Price, error) {
	if !p.currency.Equals(other.currency) {
		return Price{}, errors.New("cannot add prices with different currencies")
	}
	return NewPrice(p.amount+other.amount, p.currency), nil
}

func (p Price) Equals(other Price) bool {
	return p.amount == other.amount && p.currency.Equals(other.currency)
}
