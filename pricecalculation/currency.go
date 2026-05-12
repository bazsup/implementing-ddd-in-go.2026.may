package pricecalculation

import (
	"fmt"
	"slices"
)

var supportedCurrencies = []string{"USD"}

type Currency struct {
	value string
}

func NewCurrency(value string) (Currency, error) {
	if !slices.Contains(supportedCurrencies, value) {
		return Currency{}, fmt.Errorf("unsupported currency %q, supported: %v", value, supportedCurrencies)
	}
	return Currency{value: value}, nil
}

func (c Currency) Equals(other Currency) bool {
	return c.value == other.value
}
