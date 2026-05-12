package pricecalculation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"implementing-ddd-in-go/pricecalculation"
)

func usd(t *testing.T) pricecalculation.Currency {
	t.Helper()
	c, err := pricecalculation.NewCurrency("USD")
	require.NoError(t, err)
	return c
}

func TestNewPrice_ValidPrice(t *testing.T) {
	price, err := pricecalculation.NewPrice(7.35, usd(t))
	require.NoError(t, err)
	assert.Equal(t, 7.35, price.Amount())
	assert.Equal(t, usd(t), price.Currency())
}

func TestNewPrice_BelowZero(t *testing.T) {
	_, err := pricecalculation.NewPrice(-1, usd(t))
	assert.EqualError(t, err, pricecalculation.InvalidPriceAmount.Error())
}

func TestPrice_Equals_SameAmountAndCurrency(t *testing.T) {
	a, _ := pricecalculation.NewPrice(7.35, usd(t))
	b, _ := pricecalculation.NewPrice(7.35, usd(t))
	assert.True(t, a.Equals(b))
}

func TestPrice_Equals_DifferentAmount(t *testing.T) {
	a, _ := pricecalculation.NewPrice(7.35, usd(t))
	b, _ := pricecalculation.NewPrice(1.00, usd(t))
	assert.False(t, a.Equals(b))
}

func TestPrice_Add(t *testing.T) {
	a, _ := pricecalculation.NewPrice(1.50, usd(t))
	b, _ := pricecalculation.NewPrice(5.85, usd(t))
	sum, err := a.Add(b)
	require.NoError(t, err)
	expected, _ := pricecalculation.NewPrice(7.35, usd(t))
	assert.True(t, sum.Equals(expected))
}

func TestPrice_Add_MismatchCurrency(t *testing.T) {
	a, _ := pricecalculation.NewPrice(1.50, usd(t))
	b, _ := pricecalculation.NewPrice(5.85, pricecalculation.Currency{})
	_, err := a.Add(b)
	assert.EqualError(t, err, pricecalculation.CurrenciesDontMatch.Error())
}
