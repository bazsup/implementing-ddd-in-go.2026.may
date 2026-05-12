package pricecalculation

import (
	"errors"
	"slices"
)

var InvalidCurrency = errors.New("invalid currency")

var supportedCurrencies = []string{"USD"}

type Currency struct {
	value string
}

func NewCurrency(value string) (Currency, error) {
	if !slices.Contains(supportedCurrencies, value) {
		return Currency{}, InvalidCurrency
	}
	return Currency{value: value}, nil
}

func (c Currency) Equals(other Currency) bool {
	return c.value == other.value
}
