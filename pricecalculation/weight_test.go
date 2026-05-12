package pricecalculation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"implementing-ddd-in-go/pricecalculation"
)

func TestNewWeight_Zero(t *testing.T) {
	weight := pricecalculation.NewWeightFromKG(0)

	assert.Equal(t, uint(0), weight.Amount())
}

func TestNewWeight_Positive(t *testing.T) {
	weight := pricecalculation.NewWeightFromKG(50)

	assert.Equal(t, uint(50), weight.Amount())
}

func TestWeight_Equals_SameWeight(t *testing.T) {
	a := pricecalculation.NewWeightFromKG(50)
	b := pricecalculation.NewWeightFromKG(50)

	assert.True(t, a.Equals(b))
}

func TestWeight_Equals_DifferentWeight(t *testing.T) {
	a := pricecalculation.NewWeightFromKG(50)
	b := pricecalculation.NewWeightFromKG(100)
	
	assert.False(t, a.Equals(b))
}
