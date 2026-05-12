package pricecalculation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"implementing-ddd-in-go/pricecalculation"
)

func TestNewCurrency_ValidCurrency(t *testing.T) {
	_, err := pricecalculation.NewCurrency("USD")
	require.NoError(t, err)
}

func TestNewCurrency_InvalidCurrency(t *testing.T) {
	_, err := pricecalculation.NewCurrency("EUR")
	assert.EqualError(t, err, pricecalculation.InvalidCurrency.Error())
}

func TestCurrency_Equals_SameCurrency(t *testing.T) {
	a, _ := pricecalculation.NewCurrency("USD")
	b, _ := pricecalculation.NewCurrency("USD")
	assert.True(t, a.Equals(b))
}

func TestCurrency_Equals_DifferentCurrency(t *testing.T) {
	usd, _ := pricecalculation.NewCurrency("USD")
	var empty pricecalculation.Currency
	assert.False(t, usd.Equals(empty))
}
