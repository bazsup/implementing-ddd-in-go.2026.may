package pricecalculation

import "errors"

var InvalidPriceAmount = errors.New("invalid amount for price")
var CurrenciesDontMatch = errors.New("price currencies don't match")

type Price struct {
	amount   float64
	currency Currency
}

func NewPrice(amount float64, currency Currency) (Price, error) {
	if amount < 0 {
		return Price{}, InvalidPriceAmount
	}
	return Price{amount: amount, currency: currency}, nil
}

func (p Price) Amount() float64 {
	return p.amount
}

func (p Price) Currency() Currency {
	return p.currency
}

func (p Price) Add(other Price) (Price, error) {
	if !p.currency.Equals(other.currency) {
		return Price{}, CurrenciesDontMatch
	}
	return NewPrice(p.amount+other.amount, p.currency)
}

func (p Price) Equals(other Price) bool {
	return p.amount == other.amount && p.currency.Equals(other.currency)
}
