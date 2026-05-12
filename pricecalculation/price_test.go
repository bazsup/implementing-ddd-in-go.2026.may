package pricecalculation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"implementing-ddd-in-go/pricecalculation"
)

func TestNewPrice_ValidPrice(t *testing.T) {
	price, err := pricecalculation.NewPriceFromUSD(7.35)
	require.NoError(t, err)
	assert.Equal(t, 7.35, price.Amount())
}

func TestNewPrice_BelowZero(t *testing.T) {
	_, err := pricecalculation.NewPriceFromUSD(-1)
	assert.EqualError(t, err, pricecalculation.InvalidPriceAmount.Error())
}

func TestPrice_Equals_SameAmount(t *testing.T) {
	a, _ := pricecalculation.NewPriceFromUSD(7.35)
	b, _ := pricecalculation.NewPriceFromUSD(7.35)
	assert.True(t, a.Equals(b))
}

func TestPrice_Equals_DifferentAmount(t *testing.T) {
	a, _ := pricecalculation.NewPriceFromUSD(7.35)
	b, _ := pricecalculation.NewPriceFromUSD(1.00)
	assert.False(t, a.Equals(b))
}

func TestPrice_Add(t *testing.T) {
	a, _ := pricecalculation.NewPriceFromUSD(1.50)
	b, _ := pricecalculation.NewPriceFromUSD(5.85)
	sum, err := a.Add(b)
	require.NoError(t, err)
	expected, _ := pricecalculation.NewPriceFromUSD(7.35)
	assert.True(t, sum.Equals(expected))
}
