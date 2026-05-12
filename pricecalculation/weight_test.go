package pricecalculation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"implementing-ddd-in-go/pricecalculation"
)

func TestNewWeight_BelowZero(t *testing.T) {
	_, err := pricecalculation.NewWeightFromKG(-1)
	assert.EqualError(t, err, pricecalculation.InvalidWeight.Error())
}

func TestNewWeight_Zero(t *testing.T) {
	_, err := pricecalculation.NewWeightFromKG(0)
	require.NoError(t, err)
}

func TestNewWeight_Positive(t *testing.T) {
	weight, err := pricecalculation.NewWeightFromKG(50)
	require.NoError(t, err)
	assert.Equal(t, 50, weight.Amount())
}

func TestWeight_Equals_SameWeight(t *testing.T) {
	a, _ := pricecalculation.NewWeightFromKG(50)
	b, _ := pricecalculation.NewWeightFromKG(50)
	assert.True(t, a.Equals(b))
}

func TestWeight_Equals_DifferentWeight(t *testing.T) {
	a, _ := pricecalculation.NewWeightFromKG(50)
	b, _ := pricecalculation.NewWeightFromKG(100)
	assert.False(t, a.Equals(b))
}
