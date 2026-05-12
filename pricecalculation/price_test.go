package pricecalculation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"implementing-ddd-in-go/pricecalculation"
)

func TestNewPrice_ValidPrice(t *testing.T) {
	price := pricecalculation.NewPriceFromUSD(735)
	
	assert.Equal(t, 7.35, price.Amount())
	assert.Equal(t, "USD", price.Currency())
}

func TestPrice_Equals_SameAmount(t *testing.T) {
	a := pricecalculation.NewPriceFromUSD(735)
	b := pricecalculation.NewPriceFromUSD(735)

	assert.True(t, a.Equals(b))
}

func TestPrice_Equals_DifferentAmount(t *testing.T) {
	a := pricecalculation.NewPriceFromUSD(735)
	b := pricecalculation.NewPriceFromUSD(100)

	assert.False(t, a.Equals(b))
}

func TestPrice_Add(t *testing.T) {
	a := pricecalculation.NewPriceFromUSD(150)
	b := pricecalculation.NewPriceFromUSD(585)
	sum := a.Add(b)
	
	assert.Equal(t, 7.35, sum.Amount())
}
